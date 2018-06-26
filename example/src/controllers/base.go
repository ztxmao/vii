package controllers

import (
	"{@project}/src/common"
	"encoding/json"
	"fmt"
	"html"
	"net/http"
	"strconv"
	"time"
)

//所有Controller类的父类
type BaseController struct {
	rw             http.ResponseWriter
	r              *http.Request
	startTime      time.Time
	ControllerName string `当前controller`
	ActionName     string `当前action`
	OutputDirect   bool   //是否直接输出到http
}

func (this *BaseController) Init(rw http.ResponseWriter, r *http.Request, controllerName, actionName string) {
	this.startTime = time.Now()
	this.rw = rw
	this.r = r
	this.ControllerName = controllerName
	this.ActionName = actionName
	this.OutputDirect = true
}

func (this *BaseController) BeforeAction() bool {
	return true
}

//如果有异常，server会回调这个方法
func (this *BaseController) OutputError(err interface{}) {
	errno := common.ERR_SYSTEM.ErrorNo
	errmsg := common.ERR_SYSTEM.ErrorMsg

	switch errinfo := err.(type) {
	case string:
		errmsg = errinfo
	case common.Err:
		errno = errinfo.ErrorNo
		errmsg = errinfo.Error()
	case common.ErrorInfo:
		errno = errinfo.ErrorNo
		errmsg = errinfo.Error()
	case error:
		errmsg = errinfo.Error()
	default:
		errmsg = fmt.Sprint(errinfo)
	}
	this.toJson(map[string]interface{}{
		"errno":  errno,
		"error":  errmsg,
		"result": false,
	})
	common.Logger.Error(errmsg)
}

func (this *BaseController) GetString(key string, defaultValue string) string {
	ret := html.EscapeString(this.r.FormValue(key))
	if ret == "" {
		ret = defaultValue
	}
	return ret
}

func (this *BaseController) GetStrings(key string) []string {
	if this.r.Form == nil {
		return []string{}
	}
	vs := this.r.Form[key]
	return vs
}

func (this *BaseController) GetInt(key string, defaultValue int) int {
	ret, err := strconv.ParseInt(this.r.FormValue(key), 10, 64)
	if err != nil {
		return defaultValue
	}
	return int(ret)
}

func (this *BaseController) GetInt64(key string, defaultValue int64) int64 {
	ret, err := strconv.ParseInt(this.r.FormValue(key), 10, 64)
	if err != nil {
		return defaultValue
	}
	return ret
}

func (this *BaseController) GetBool(key string, defaultValue bool) bool {
	ret, err := strconv.ParseBool(this.r.FormValue(key))
	if err != nil {
		ret = defaultValue
	}
	return ret
}

//api统一json输出
func (this *BaseController) Output(data interface{}) {
	this.toJson(map[string]interface{}{
		"errno":  0,
		"error":  "ok",
		"result": data,
	})
}

//字符串输出
func (this *BaseController) OutputString(data string) {
	this.writeToWriter([]byte(data))
}

//原始json格式输出
func (this *BaseController) OutputNoFmt(data interface{}) {
	this.toJson(data)
}

func (this *BaseController) toJson(data interface{}) {
	content, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		common.Logger.Error("output jsonencode err", err)
		content = ([]byte)("{\"errno\":1,\"error\":\"系统开了小差!\",\"result\":false}")
	}
	this.rw.Header().Set("Content-Type", "application/json;charset=UTF-8")
	this.writeToWriter(content)
}

//获取需要打印到日志的信息
func (this *BaseController) genLog() map[string]interface{} {
	ret := make(map[string]interface{})
	for key, value := range this.r.Form {
		if len(value) > 1 {
			ret[key] = value
		} else {
			ret[key] = value[0]
		}
	}
	//访问ip
	ret["user_ip"] = this.r.RemoteAddr
	//请求路径
	ret["req_uri"] = this.r.URL.Path
	//时间消耗
	ret["time_cost"] = this.cost()
	return ret
}

func (this *BaseController) cost() int64 {
	return time.Now().Sub(this.startTime).Nanoseconds() / 1000 / 1000
}

func (this *BaseController) writeToWriter(rb []byte) {
	//this.rw.Header().Set("Content-Length", strconv.Itoa(len(rb)))
	if this.OutputDirect {
		this.rw.Write(rb)
	}
	common.Logger.Access(this.genLog(), len(rb))
}

func (this *BaseController) httpErr(code int, info string) {
	http.Error(this.rw, info, code)
}
