package logic

import (
	"flag"
	"fmt"
	"{@project}/src/common"
	"testing"
	"time"
)

import (
	"{@project}/src/models"
)

var (
	p = fmt.Println
)

//其他的test文件就不能用init 函数 不自然不知道先加载的谁问导致未知错误
func init() {
	confiPath := flag.String("f", "../../conf/server.conf", "config file")
	flag.Parse()
	if *confiPath == "" {
		panic("config file missing")
	}
	common.Configer.Init(*confiPath)
	logPath := common.Configer.GetStr("log", "path")
	logName := common.Configer.GetStr("log", "name")
	logLevel := common.Configer.GetInt("log", "level")
	common.Logger.Init(logPath, logName, logLevel)
	models.Init()
}
