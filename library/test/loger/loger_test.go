package loger_test

import (
	"fmt"
	"github.com/ztxmao/vii/library"
	"testing"
)

var (
	logFile = "./"
	l       = library.Logger
	p       = fmt.Println
)

func init() {
	l.Init(logFile, "TEST", 0)
}

func TestLogger(t *testing.T) {
	l.Access("test", "just loger access test")
	l.Debug("test", "just loger debug test")
	l.Warn("test", "just loger warn test")
	l.Error("test", "just loger error test")

}
