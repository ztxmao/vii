package tasker_test

import (
	//	"errors"
	"fmt"
	"library"
	"testing"
	"time"
)

var (
	p = fmt.Println
	f = func(expire time.Duration) error {
		p(time.Now(), " task runing!")
		return nil
	}
)

func TestTasker(t *testing.T) {
	rate := time.Second * 1
	expire := time.Second * 2
	tasker, _ := library.NewTasker(rate, expire, f)
	go tasker.Run()
	time.Sleep(time.Second * 2)
	p("tasker pause")
	tasker.Pause()
	p("tasker start")
	tasker.Start()
	time.Sleep(time.Second * 2)
	p("tasker stop")
	tasker.Stop()
	time.Sleep(time.Second * 3)
}
