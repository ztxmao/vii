package main

import (
	"{@project}/src/common"
	"{@project}/src/controllers"
	"{@project}/src/logic"
	"{@project}/src/models"
	"flag"
	"fmt"
	"github.com/ztxmao/vii/frame"
	"github.com/ztxmao/vii/frame/router"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"time"
)

func main() {
	defer func() {
		managePid(false) //删除pid文件
		println("Server Exit")
	}()
	envInit()
	managePid(true) //生成pid文件
	common.Logger.Debug("Init Models start...")
	modelsInit()
	common.Logger.Debug("Init Models done...")
	common.Logger.Debug("Init Data start...")
	dataInit()
	common.Logger.Debug("Init Data done...")
	common.Logger.Debug("Server start...")
	startServer()
	common.Logger.Debug("Server done...")
}

//所有初始化操作
func envInit() {
	os.Chdir(path.Dir(os.Args[0]))
	defaultConfPath := os.Getenv("GOPATH") + "/src/{@project}/conf/server.ini"
	confiPath := flag.String("f", defaultConfPath, "config file")
	env := flag.String("env", "", "环境参数")
	if *env == "" {
		*env = os.Getenv("GO_PROJECT_ENV")
	}
	flag.Parse()
	if *confiPath == "" {
		panic("config file missing")
	}
	common.Configer.Init(*confiPath, *env)
	logPath := common.Configer.GetStr("log", "path")
	logName := common.Configer.GetStr("log", "name")
	logLevel := common.Configer.GetInt("log", "level")
	common.Logger.Init(logPath, logName, logLevel)
}

func dataInit() {
	logic.Init()
}
func modelsInit() {
	models.Init()
}

//启动httpserver
func startServer() {
	//路由初始化
	defaultRouter := router.NewStaticRouter("", "monitor", "status")
	defaultRouter.AddController(&controllers.MonitorController{})
	defaultRouter.AddController(&controllers.ReloadController{})
	defaultRouter.AddController(&controllers.ConsoleController{})
	//配置文件获取
	port := common.Configer.GetInt("server", "port")
	rTimeOut := common.Configer.GetInt("server", "r_time")
	wTimeOut := common.Configer.GetInt("server", "w_time")
	pprofAddr := common.Configer.GetStr("server", "pprof_addr")
	//server 初始化
	server := frame.NewHttpServer("", port, rTimeOut, wTimeOut, pprofAddr)
	server.AddRouter("default", defaultRouter)

	fmt.Println("Server Start: ", time.Now())
	//server run
	server.Run()
}

//生成/删除当前进程id文件
func managePid(create bool) {
	pidFile := common.Configer.GetStr("server", "pidfile")
	fmt.Println(pidFile)
	if create {
		pid := os.Getpid()
		pidString := strconv.Itoa(pid)
		ioutil.WriteFile(pidFile, []byte(pidString), 0777)
		os.Chmod(pidFile, 0777)
	} else {
		os.Remove(pidFile)
	}
}
