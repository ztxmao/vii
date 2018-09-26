package frame

import (
	"fmt"
	"github.com/ztxmao/vii/frame/router"
	"net/http"
	_ "net/http/pprof"
	"reflect"
	"runtime"
	"time"

	"github.com/tabalt/gracehttp"
)

//http服务监听,路由
type HttpApiServer struct {
	HttpAddr     string
	HttpPort     int
	ReadeTimeout int
	WriteTimeout int
	IdleTimeout  int
	hanlder      *httpApiHandler
	pprofAddr    string
}

type httpApiHandler struct {
	routMap     map[string]map[string]reflect.Type //key:controller: {key:method value:reflect.type}
	routersPool map[string]router.Router
	timeout     time.Duration
}

//new server
func NewHttpServer(addr string, port, readTimout, witeTimeout, idleTimeout, execTimout int, pprofAddr string) *HttpApiServer {
	handleTimeout := time.Duration(execTimout) * time.Millisecond
	ret := &HttpApiServer{
		HttpAddr:     addr,
		HttpPort:     port,
		ReadeTimeout: readTimout,
		WriteTimeout: witeTimeout,
		IdleTimeout:  idleTimeout,
		pprofAddr:    pprofAddr,
		hanlder:      &httpApiHandler{routMap: make(map[string]map[string]reflect.Type), routersPool: make(map[string]router.Router), timeout: handleTimeout},
	}
	return ret
}

//添加路由
func (this *HttpApiServer) AddRouter(name string, route router.Router) {
	if route == nil {
		panic("route is nil!")
	}
	this.hanlder.routersPool[name] = route
}

// server run
func (this *HttpApiServer) Run() {
	//	runtime.GOMAXPROCS(runtime.NumCPU())
	addr := fmt.Sprintf("%s:%d", this.HttpAddr, this.HttpPort)
	readTimeout := time.Duration(this.ReadeTimeout) * time.Millisecond
	writeTimeout := time.Duration(this.WriteTimeout) * time.Millisecond
	idleTimeout := time.Duration(this.IdleTimeout) * time.Second
	if this.pprofAddr != "" {

		fmt.Println("HttpApiServer PProf Listen: ", this.pprofAddr)
		go http.ListenAndServe(this.pprofAddr, nil)
	}
	fmt.Println("HttpApiServer Listen: ", addr)
	if err := gracehttp.NewServer(addr, this.hanlder, readTimeout, writeTimeout, idleTimeout).ListenAndServe(); err != nil {
		panic(err)
	}
}
func stack() string {
	var buf [2 << 10]byte
	return string(buf[:runtime.Stack(buf[:], true)])
}
func (this *httpApiHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("ServeHTTP: ", err)
			//			fmt.Println(stack())
			http.Error(rw, fmt.Sprintln(err), http.StatusInternalServerError)
		}
	}()
	after := time.After(this.timeout)
	done := make(chan int, 1)
	rw.Header().Set("Server", "GoServer")
	r.ParseForm()
	var contollerType reflect.Type
	var ok bool
	var cname, mname string
	for _, route := range this.routersPool {
		if cname, mname, contollerType, ok = route.GetRouter(r); ok {
			break
		}
	}
	if ok == false {
		http.NotFound(rw, r)
		return
	}
	vc := reflect.New(contollerType)
	var in []reflect.Value
	var method reflect.Value

	defer func() {
		if err := recover(); err != nil {
			in = []reflect.Value{reflect.ValueOf(err)}
			method := vc.MethodByName("OutputError")
			method.Call(in)
		}
	}()
	in = make([]reflect.Value, 5)
	in[0] = reflect.ValueOf(rw)
	in[1] = reflect.ValueOf(r)
	in[2] = reflect.ValueOf(cname)
	in[3] = reflect.ValueOf(mname)
	in[4] = reflect.ValueOf(done)
	method = vc.MethodByName("Init")
	method.Call(in)
	in = make([]reflect.Value, 0)
	beforeMethod := vc.MethodByName("BeforeAction")
	beforeRes := beforeMethod.Call(in)
	if beforeRes[0].Bool() == false {
		return
	}
	method = vc.MethodByName(mname)
	go method.Call(in)
	select {
	case <-done:
		close(done)
		break
	case <-after:
		fmt.Println(time.Now(), "ServeHTTP: timeout")
		http.Error(rw, "timeout", http.StatusBadGateway)
		return
	}
}
