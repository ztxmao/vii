package controllers

import (
	c "{@project}/src/common"
	"strings"
)

//示例Controller
type ReloadController struct {
	BaseController
}

func (this *ReloadController) BeforeAction() bool {
	ipArr := make([]string, 4)
	ipArr = strings.Split(this.r.RemoteAddr, ":")
	ip := ipArr[0]
	if ip == "127.0.0.1" || ip == "[" {
		return true
	}
	this.httpErr(403, "ip limit")
	return false
}
func (this *ReloadController) ConfAction() {
	c.Logger.Debug("Reload config...")
	go c.Configer.ReadList()
	this.rw.Header().Set("Content-Type", "text/html;charset=UTF-8")
	this.writeToWriter(([]byte)("ok\r\n"))
	return

}
