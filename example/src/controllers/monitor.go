package controllers

import (
	"errors"
)

import (
	c "{@project}/src/common"
)

//示例Controller
type MonitorController struct {
	BaseController
}

func (this *MonitorController) StatusAction() {
	this.rw.Header().Set("Content-Type", "text/html;charset=UTF-8")
	this.OutputString("ok\n")
	return

}

func (this *MonitorController) ErrorsAction() {

	err := errors.New("错误输出demo")
	this.OutputError(c.ErrorInfo{
		Err: c.ERR_SYSTEM,
		Raw: err,
	})
	return

}

func (this *MonitorController) PassAction() {
	ret := make(map[string]string)
	ret["REV"] = "OK"
	ret["TITLE"] = "monitor pass demo"
	ret["DATA"] = ""
	this.OutputNoFmt(ret)
	return //通过函数返回时为了兼容从MonitorController发起的调用
}

func (this *MonitorController) FailAction() {
	ret := make(map[string]string)
	ret["REV"] = "FAILED"
	ret["TITLE"] = "monitor fail demo"
	ret["DATA"] = "fail!"
	this.OutputNoFmt(ret)
	return //通过函数返回时为了兼容从MonitorController发起的调用
}
