package controllers

import (
	c "{@project}/src/common"
	"strings"
)

//示例Controller
type ConsoleController struct {
	BaseController
}

func (this *ConsoleController) BeforeAction() bool {
	ipArr := make([]string, 4)
	ipArr = strings.Split(this.r.RemoteAddr, ":")
	ip := ipArr[0]
	if ip == "127.0.0.1" || ip == "[" {
		return true
	}
	this.httpErr(403, "ip limit")
	return false
}
func (this *ConsoleController) HelpAction() {
	c.Logger.Debug("Console Qbus Source...")
	this.rw.Header().Set("Content-Type", "text/html;charset=UTF-8")
	this.OutputString("you help message")

}
