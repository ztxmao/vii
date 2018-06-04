package sortsets

import (
	"github.com/emirpasic/gods/lists/arraylist"
	"github.com/emirpasic/gods/maps/hashmap"
	"github.com/emirpasic/gods/utils"
)

import (
	"errors"
	"time"
	//	"fmt"
)

type SSet struct {
	list       *(arraylist.List)                   `简单切片列表`
	m          *(hashmap.Map)                      `hash map 方便查找元素`
	m_index    map[interface{}]int                 `index map 维护元素index`
	f          func(value interface{}) interface{} `hash key 生成函数`
	item_type  string                              `数据单元类型标示，在集合运算要求同种类型的数据单元`
	createline int64
}

// New instantiates a hash map.
func New(getkey func(value interface{}) interface{}, vtype string) *SSet {
	var set SSet
	set.list = arraylist.New()
	set.m = hashmap.New()
	set.m_index = make(map[interface{}]int)
	set.f = getkey
	set.item_type = vtype
	set.createline = time.Now().Unix()
	return &set
}

//添加元素到集合
func (sset *SSet) Add(value interface{}) {
	key := sset.f(value)
	if _, found := sset.m.Get(key); !found {
		sset.list.Add(value)
		sset.m.Put(key, value)
		index := sset.list.Size()
		sset.m_index[key] = index - 1
	}
}

//获取集合所有元素（有序）
func (sset *SSet) Values() []interface{} {
	return sset.list.Values()
}

//是否是空集合
func (sset *SSet) Empty() bool {
	return sset.list.Empty()
}

//集合元素个数
func (sset *SSet) Size() int {
	return sset.list.Size()
}

//清空集合
func (sset *SSet) Clear() {
	sset.list.Clear()
	sset.m.Clear()
	sset.m_index = make(map[interface{}]int)
}

//排序
func (sset *SSet) Sort(comparator utils.Comparator) {
	sset.list.Sort(comparator)
	sset.list.Each(func(index int, value interface{}) {
		key := sset.f(value)
		sset.m_index[key] = index
	})
}

//TODO 稳定插入 已经有的数据不改变顺序
func (sset *SSet) Insert(index int, value interface{}) {
	key := sset.f(value)
	if _, found := sset.m.Get(key); found {
		if oldIndex, found := sset.m_index[key]; found {
			sset.list.Remove(oldIndex)
		}
	}
	sset.list.Insert(index, value)
	sset.m.Put(key, value)
	sset.fixIndex()
}

// Remove removes one or more elements from the list with the supplied indices.
func (sset *SSet) Remove(value interface{}) {
	key := sset.f(value)
	if index, found := sset.m_index[key]; found {
		sset.list.Remove(index)
		sset.m.Remove(key)
		delete(sset.m_index, key)
		sset.fixIndex()
	}
}

//TODO  如何完美保证数据一致性
func (sset *SSet) fixIndex() {
	sset.m_index = make(map[interface{}]int)
	sset.list.Each(func(index int, value interface{}) {
		key := sset.f(value)
		sset.m_index[key] = index
	})
}

//TODO  如何完美保证数据一致性
func (sset *SSet) fixMap() {
	sset.m.Clear()
	sset.list.Each(func(index int, value interface{}) {
		key := sset.f(value)
		sset.m.Put(key, value)
	})
}

// String returns a string representation of container
func (sset *SSet) String() string {
	str := "SSet\n"
	list_values := sset.list.String() + "\n"
	map_values := sset.m.String() + "\n"
	str = str + list_values + map_values
	return str
}

//迭代器获取
func (sset *SSet) Iterator() arraylist.Iterator {
	return sset.list.Iterator()
}

//按照切片索引获取item
func (sset *SSet) Get(index int) (interface{}, bool) {
	return sset.list.Get(index)
}

//查找 指定item
func (sset *SSet) Find(value interface{}) (interface{}, int, bool) {
	key := sset.f(value)
	if index, found := sset.m_index[key]; found {
		if rdata, res := sset.Get(index); res {
			return rdata, index, true
		}
	}
	return nil, -1, false
}

func (sset *SSet) GetCreateline() int64 {
	return sset.createline
}

//和指定切片求差集
func (sset *SSet) ArrayDiff(values []interface{}) (res_values []interface{}) {
	for _, value := range values {
		key := sset.f(value)
		if _, found := sset.m.Get(key); !found {
			res_values = append(res_values, value)
		}
	}
	return
}

//和指定切片求交集
func (sset *SSet) ArrayIntersect(values []interface{}) (res_values []interface{}) {
	for _, value := range values {
		key := sset.f(value)
		if _, found := sset.m.Get(key); found {
			res_values = append(res_values, value)
		}
	}
	return
}

//两个集合求交集 时间复杂度 O[n]
func (sset *SSet) Intersect(new_sset *SSet) (res_set *SSet, err error) {
	if sset.item_type == new_sset.item_type && sset.Size() > 0 && new_sset.Size() > 0 {
		res_set = New(sset.f, sset.item_type)
		new_sset.list.Each(func(index int, value interface {
		}) {
			if a_item, _, found := sset.Find(value); found {
				res_set.Add(a_item)
			}
		})
		return res_set, nil
	} else {
		return nil, errors.New("type is not same!")
	}
}

//两个集合求差集 时间复杂度 O[n]
func (sset *SSet) Diff(new_sset *SSet) (res_set *SSet, err error) {

	if sset.item_type == new_sset.item_type && sset.Size() > 0 && new_sset.Size() > 0 {
		res_set = New(sset.f, sset.item_type)
		new_sset.list.Each(func(index int, value interface {
		}) {
			if _, _, found := sset.Find(value); !found {
				res_set.Add(value)
			}
		})
		return res_set, nil
	} else {
		return nil, errors.New("type is nsortsets.goot same!")
	}
}

//两个集合求并集 时间复杂度 O[n]
func (sset *SSet) Union(new_sset *SSet) (res_set *SSet, err error) {

	if sset.item_type == new_sset.item_type && sset.Size() > 0 && new_sset.Size() > 0 {
		res_set = New(sset.f, sset.item_type)
		new_sset.list.Each(func(index int, value interface {
		}) {
			res_set.Add(value)
		})
		sset.list.Each(func(index int, value interface {
		}) {
			res_set.Add(value)
		})
		return res_set, nil
	} else {
		return nil, errors.New("type is not same!")
	}
}
