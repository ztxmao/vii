package exsets

import (
	"sync"
	"testing"
)

var (
	ExSet = New("test")
)
var waitGroup sync.WaitGroup

func TestIncr(t *testing.T) {
	var f = func(key interface{}) {
		ExSet.IncrBy(key, 1)
		waitGroup.Add(-1)
	}
	for i := 0; i < 100; i++ {
		waitGroup.Add(1)
		go f(i)
	}
	waitGroup.Wait()
	//	p(ExSet)
	t.Errorf("ExSet Get Values Err!")
}

func TestIncrMap(t *testing.T) {
	tMap := make(map[interface{}]interface{})
	var f = func(key interface{}) {
		tMap[key] = key
	}
	for i := 0; i < 1000000; i++ {
		//waitGroup.Add(1)
		f(i)
	}
	//waitGroup.Wait()
	//	p(ExSet)
	t.Errorf("ExSet Get Values Err!")
}
