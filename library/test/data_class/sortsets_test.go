package data_class_test

import (
	"fmt"
	dc "github.com/ztxmao/vii/library/data_class/exsets"
	"testing"
)

type VideoItem struct {
	Vid      int64  `视频ID`
	Title    string `标题`
	Priority int    `排序权重`
}

var (
	getkey = func(item interface{}) (key interface{}) {
		key = fmt.Sprintf("%d", item.(VideoItem).Vid)
		return
	}
	tset = dc.New(getkey, "VideoItem")
)

//SSet
//ArrayList
//{1 test1 3}, {2 test2 1}, {3 test3 2}
//HashMap
//map[2:{2 test2 1} 3:{3 test3 2} 1:{1 test1 3}]
func InitSet() {
	tset.Clear()
	a := VideoItem{1, "test1", 3}
	b := VideoItem{2, "test2", 1}
	c := VideoItem{3, "test3", 2}
	tset.Add(a)
	tset.Add(b)
	tset.Add(c)
}
func TestSSetValues(t *testing.T) {
	InitSet()
	strSet := fmt.Sprint(tset.Values())
	if strSet != "[{1 test1 3} {2 test2 1} {3 test3 2}]" {
		t.Errorf("SSet Get Values Err!")
	}
}

func TestSSetString(t *testing.T) {
	InitSet()
	c := VideoItem{3, "test3", 2}
	values := tset.Values()
	count := tset.Size()
	item, index, rs := tset.Find(c)
	if item2, rs := tset.Get(2); !rs || item2 != c {
		t.Errorf("SSet get item  fail!")
	}
	if len(values) != 3 || count != 3 {
		t.Errorf("SSet size is error!")
	}
	if rs == false || index != 2 || item != c {
		t.Errorf("SSet fand item fail !")
	}
}
func TestSSetClear(t *testing.T) {
	InitSet()
	c := VideoItem{3, "test3", 2}
	tset.Clear()
	values := tset.Values()
	count := tset.Size()
	_, _, rs := tset.Find(c)
	if count != 0 || len(values) != 0 {
		t.Errorf("SSet size clear fail!")
	}
	if rs != false {
		t.Errorf("SSet item clear fail !")
	}
}

func TestSSetEmpty(t *testing.T) {
	InitSet()
	if tset.Empty() == true {
		t.Errorf("SSet verify empty (init) fail !")
	}
	tset.Clear()
	if tset.Empty() == false {
		t.Errorf("SSet verify empty (clear) fail !")
	}
}

func TestSSetInstert(t *testing.T) {
	InitSet()
	d := VideoItem{4, "test4", 10}
	e := VideoItem{5, "test5", 10}
	tset.Insert(1, d)
	tset.Insert(2, d)
	tset.Insert(2, e)
	strValues := fmt.Sprint(tset.Values())
	if strValues != "[{1 test1 3} {2 test2 1} {5 test5 10} {4 test4 10} {3 test3 2}]" {
		t.Errorf("SSet insert item fail !")
	}
	if item, index, rs := tset.Find(d); !rs || item != d || index != 3 {
		t.Errorf("SSet insert item location error!")
	}
}

func TestSSetSort(t *testing.T) {
	InitSet()
	d := VideoItem{4, "test4", 10}
	tset.Insert(1, d)
	fsort := func(a, b interface{}) int {
		a_priority := a.(VideoItem).Priority
		b_priority := b.(VideoItem).Priority
		switch {
		case a_priority > b_priority:
			return -1
		case a_priority < b_priority:
			return 1
		default:
			return 0
		}
	}
	tset.Sort(fsort)
	strValue := fmt.Sprint(tset.Values())
	if strValue != "[{4 test4 10} {1 test1 3} {3 test3 2} {2 test2 1}]" {
		t.Errorf("SSet values sort error!")
	}
}

func TestSSetIterator(t *testing.T) {
	InitSet()
	it := tset.Iterator()
	strValue := ""
	for it.Next() {
		strValue += fmt.Sprint(it.Value())
	}
	if strValue != "{1 test1 3}{2 test2 1}{3 test3 2}" {
		t.Errorf("SSet iterator get value error!")
	}
}

func TestSSetRemove(t *testing.T) {
	InitSet()
	b := VideoItem{2, "test2", 1}
	tset.Remove(b)
	strValue := fmt.Sprint(tset.Values())
	if strValue != "[{1 test1 3} {3 test3 2}]" {
		t.Errorf("SSet remove item fail!")
	}
}

func TestSSetDiffArray(t *testing.T) {
	InitSet()
	a := VideoItem{1, "test1", 3}
	b := VideoItem{2, "test2", 1}
	c := VideoItem{3, "test3", 2}
	d := VideoItem{4, "test4", 10}
	items := []interface{}{a, b, c, d}
	diff := tset.ArrayDiff(items)
	strValue := fmt.Sprint(diff)
	if strValue != "[{4 test4 10}]" {
		t.Errorf("SSet diff array result error!")
	}
}

func TestSSetIntersectArray(t *testing.T) {
	InitSet()
	a := VideoItem{1, "test1", 3}
	b := VideoItem{2, "test2", 1}
	c := VideoItem{3, "test3", 2}
	d := VideoItem{4, "test4", 10}
	items := []interface{}{a, b, c, d}
	diff := tset.ArrayIntersect(items)
	strValue := fmt.Sprint(diff)
	if strValue != "[{1 test1 3} {2 test2 1} {3 test3 2}]" {
		t.Errorf("SSet intersect array result error!")
	}
}

func TestSSetDiffSet(t *testing.T) {
	InitSet()
	InitSet()
	a := VideoItem{1, "test1", 3}
	b := VideoItem{2, "test2", 1}
	c := VideoItem{3, "test3", 2}
	d := VideoItem{4, "test4", 10}
	e := VideoItem{5, "test5", 10}
	set_b := dc.New(getkey, "VideoItem")
	set_b.Add(c)
	set_b.Add(b)
	set_b.Add(a)
	set_b.Add(e)
	set_b.Add(d)
	diff, err := tset.Diff(set_b)
	strValue := fmt.Sprint(diff.Values())
	if err != nil || strValue != "[{5 test5 10} {4 test4 10}]" {
		t.Errorf("SSet diff set result error!")
	}
}

func TestSSetIntersectSet(t *testing.T) {
	InitSet()
	a := VideoItem{1, "test1", 3}
	b := VideoItem{2, "test2", 1}
	c := VideoItem{3, "test3", 2}
	d := VideoItem{4, "test4", 10}

	set_b := dc.New(getkey, "VideoItem")
	set_b.Add(c)
	set_b.Add(b)
	set_b.Add(a)
	set_b.Add(d)
	diff, err := tset.Intersect(set_b)
	strValue := fmt.Sprint(diff.Values())
	if err != nil || strValue != "[{3 test3 2} {2 test2 1} {1 test1 3}]" {
		t.Errorf("SSet intersect set result error!")
	}
}
func TestSSetunionSet(t *testing.T) {
	InitSet()
	a := VideoItem{1, "test1", 3}
	b := VideoItem{2, "test2", 1}
	c := VideoItem{3, "test3", 2}
	d := VideoItem{4, "test4", 10}

	set_b := dc.New(getkey, "VideoItem")
	set_b.Add(c)
	set_b.Add(b)
	set_b.Add(a)
	set_b.Add(d)
	diff, err := tset.Union(set_b)
	strValue := fmt.Sprint(diff.Values())
	if err != nil || strValue != "[{3 test3 2} {2 test2 1} {1 test1 3} {4 test4 10}]" {
		t.Errorf("SSet intersect set result error!")
	}
}
