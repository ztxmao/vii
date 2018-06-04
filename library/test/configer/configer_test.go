package ini_configer_test

import (
	"fmt"
	"library"
	"testing"
)

var (
	configFile = "app.conf"
	c          = library.Configer
	p          = fmt.Println
)

func init() {
	c.Init(configFile)
}

func TestConfiger(t *testing.T) {
	t1 := c.GetStr("database", "username")
	t2 := c.GetInt("database", "port")
	t3, _ := c.GetSection("database")
	t4 := c.GetBool("database", "debug")
	if t1 != "root" || t2 != 3306 || t3 == nil || t4 != true {
		t.Errorf("configer get data err! t1=%v t2=%v t3=%v t4=%v", t1, t2, t3, t4)
	}
}
