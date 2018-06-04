package data_class_test

import (
	"fmt"
	"library"
	dc "library/data_class/exsets"
	"testing"
	"time"
)

var (
	ExSet = dc.New("VideoItem")
)

func InitExset() {
	ExSet.Clear()
	a := VideoItem{1, "test1", 3}
	b := VideoItem{2, "test2", 1}
	c := VideoItem{3, "test3", 2}
	ExSet.Add(a.Vid, a)
	ExSet.Add(b.Vid, b)
	ExSet.Add(c.Vid, c)
}

func TestExSetValues(t *testing.T) {
	InitExset()
	strSet := fmt.Sprint(ExSet.Values())
	fmt.Println(strSet)
	if len(strSet) < 10 {
		t.Errorf("ExSet Get Values Err!")
	}
}

func TestExSetString(t *testing.T) {
	InitExset()
	strSet := fmt.Sprintf("%s", ExSet)
	fmt.Println(strSet)
	if len(strSet) < 10 {
		t.Errorf("ExSet Get Values Err!")
	}
}

func TestExSetClear(t *testing.T) {
	InitExset()
	ExSet.Clear()
	values := ExSet.Values()
	count := ExSet.Size()
	item, ok := ExSet.Get("1")
	fmt.Println(values, count, item)
	if count != 0 || len(values) != 0 || ok {
		t.Errorf("ExSet size clear fail!")
	}
}

func TestExSetGet(t *testing.T) {
	InitExset()
	d := VideoItem{4, "test4", 11}
	ExSet.AddEx(d.Vid, d, 10)
	item, ok := ExSet.Get(d.Vid)
	fmt.Println(ExSet, item)
	if !ok {
		t.Errorf("ExSet size clear fail!")
	}
}

func TestExSetGetTtl(t *testing.T) {
	InitExset()
	d := VideoItem{4, "test4", 11}
	ExSet.AddEx(d.Vid, d, 1)
	time.Sleep(time.Second * 1)
	item, ok := ExSet.Get(d.Vid)
	fmt.Println(ExSet, item)
	if ok {
		t.Errorf("ExSet size clear fail!")
	}
}

func TestExSetEmpty(t *testing.T) {
	InitExset()
	if ExSet.Empty() == true {
		t.Errorf("ExSet verify empty (init) fail !")
	}
	ExSet.Clear()
	fmt.Println(ExSet, time.Now().Unix())
	if ExSet.Empty() == false {
		t.Errorf("ExSet verify empty (clear) fail !")
	}
}

func TestExSetRemove(t *testing.T) {
	InitExset()
	b := VideoItem{2, "test2", 1}
	ExSet.AddEx(b.Vid, b, 10)
	ExSet.Remove(b.Vid)
	fmt.Println(ExSet, time.Now().Unix())
	if _, ok := ExSet.Get(b.Vid); ok {
		t.Errorf("ExSet remove item fail!")
	}
}

func TestExSetTtl(t *testing.T) {
	InitExset()
	b := VideoItem{2, "test2", 10}
	ExSet.AddEx(b.Vid, b, 10)
	time.Sleep(time.Second)
	ttl := ExSet.Ttl(b.Vid)
	fmt.Println(ExSet, ttl)
	if ttl <= 0 {
		t.Errorf("ExSet get ttl fail!")
	}
}

func TestExSetExpire(t *testing.T) {
	InitExset()
	b := VideoItem{2, "test2", 10}
	ExSet.Expire(b.Vid, 2)
	time.Sleep(time.Second)
	ttl := ExSet.Ttl(b.Vid)
	if ttl <= 0 {
		t.Errorf("ExSet expire add fail!")
	}
	ExSet.Expire(b.Vid, 0)
	ttl = ExSet.Ttl(b.Vid)
	if ttl > 0 {
		t.Errorf("ExSet expire del fail!")
	}
}
func TestExSetIncr(t *testing.T) {
	InitExset()
	a, b := ExSet.IncrBy(22, 2)
	if !b {
		t.Errorf("ExSet Incr fail(1)!")
	}
	fmt.Println(ExSet, a, b)
	a, b = ExSet.IncrBy(int64(2), 2)
	fmt.Println(ExSet, a, b)
	if b {
		t.Errorf("ExSet Incr fail(2)!")
	}
}

func TestCron(t *testing.T) {
	InitExset()
	b := VideoItem{2, "test2", 4}
	ExSet.AddEx(b.Vid, b, 4)
	rate := time.Second * 1
	expire := time.Second * 2
	f := func(expire time.Duration) error {
		fmt.Println(time.Now(), " task runing!")
		ExSet.CronClear()
		return nil
	}
	tasker, _ := library.NewTasker(rate, expire, f)
	time.Sleep(time.Second * 5)
	fmt.Println(ExSet)
	go tasker.Run()
	time.Sleep(time.Second * 2)
	fmt.Println(ExSet)
	//	p("tasker pause")
	//	tasker.Pause()
	//	p("tasker start")
	//	tasker.Start()
	//	time.Sleep(time.Second * 2)
	//	p("tasker stop")
	//	tasker.Stop()
	//	time.Sleep(time.Second * 3)
}
