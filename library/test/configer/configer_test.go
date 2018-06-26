package ini_configer_test

import (
	"fmt"
	"github.com/ztxmao/vii/library"
	"testing"
)

var (
	configFile    = "app.conf"
	configIniFile = "app.ini"
	c             = library.Configer
	ce            = library.ConfigerExt
	p             = fmt.Println
)

func init() {
	c.Init(configFile)
	ce.Init(configIniFile, "dev")
}

func TestConfiger(t *testing.T) {
	t1 := c.GetStr("database", "username")
	t2 := c.GetInt("database", "port")
	t3, _ := c.GetSection("database")
	t4 := c.GetBool("database", "debug")
	p(t1, t2, t3, t4)
}

func TestConfigerExt(t *testing.T) {
	t1 := ce.GetStr("user", "identityCookie.domain")
	t2 := ce.GetInt("user", "authTimeout")
	t3, _ := ce.GetSection("picasso")
	t4 := ce.GetBool("log", "debug")
	p(t1, t2, t3, t4)
}
