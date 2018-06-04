package mapsorter_test

import (
	//	"errors"
	"fmt"
	l "library"
	"testing"
)

func TestMapSorter(t *testing.T) {
	m := map[string]int{
		"One":   100,
		"Two":   2,
		"Three": 3,
		"Ten":   10,
		"Fifty": 50,
	}
	fmt.Println(m)
	fmt.Printf("map:%#v\n", m)
	vs := l.SortMapByValue(m)
	fmt.Printf("sort by val:%#v\n", vs)
	ks := l.SortMapByKey(m)
	fmt.Printf("sort by key:%#v\n", ks)
}
