package exsets

import (
	"github.com/emirpasic/gods/maps/hashmap"
)

import (
	"fmt"
	"sync"
	"time"
)

/**
* 带有过期时间的集合
* 删除策略 定时删除+惰性删除
**/
type ESet struct {
	m          *(hashmap.Map) `hash map 方便查找元素`
	em         *(hashmap.Map) `过期时间维护`
	item_type  string         `数据单元类型标示，在集合运算要求同种类型的数据单元`
	createline int64
	mutex_lock *sync.Mutex `锁住整个集合`
}

var (
	p = fmt.Println
)

//过期时间字典

// New instantiates a hash map.
func New(dataType string) *ESet {
	var set ESet
	set.m = hashmap.New()
	set.em = hashmap.New()
	set.item_type = dataType
	set.createline = time.Now().Unix()
	set.mutex_lock = new(sync.Mutex)
	return &set
}

//添加元素到集合
func (sset *ESet) Add(key interface{}, value interface{}) {
	sset.mutex_lock.Lock()
	sset.m.Put(key, value)
	sset.mutex_lock.Unlock()
}

//添加元素到集合
func (sset *ESet) AddEx(key interface{}, value interface{}, expire int64) {
	nowTemp := time.Now().Unix()
	sset.mutex_lock.Lock()
	sset.m.Put(key, value)
	if expire > 0 {
		sset.em.Put(key, nowTemp+expire)
	}
	sset.mutex_lock.Unlock()
}
func (sset *ESet) getExpire(key interface{}) int64 {
	if expireValue, found := sset.em.Get(key); found {
		if int64Var, ok := expireValue.(int64); ok {
			return int64Var
		}
	}
	return 0

}

//获取集合所有元素
func (sset *ESet) Values() (values map[interface{}]interface{}) {
	values = make(map[interface{}]interface{})
	keys := sset.m.Keys()
	nowTemp := time.Now().Unix()
	for _, key := range keys {
		expireTemp := sset.getExpire(key)
		if expireTemp > 0 && expireTemp < nowTemp {
			sset.Remove(key)
		} else {
			values[key], _ = sset.m.Get(key)
		}
	}
	return
}

//是否是空集合
func (sset *ESet) Empty() bool {
	return sset.m.Empty()
}

//集合元素个数包含过期key
func (sset *ESet) Size() int {
	return sset.m.Size()
}

//设置过期时间的key个数
func (sset *ESet) ExSize() int {
	return sset.em.Size()
}

//查看value过期时间
func (sset *ESet) Ttl(key interface{}) int64 {
	expireTemp := sset.getExpire(key)
	nowTemp := time.Now().Unix()
	if expireTemp > 0 && expireTemp > nowTemp {
		return expireTemp - nowTemp
	}

	if expireTemp > 0 && expireTemp <= nowTemp {
		return 0
	}
	return -1

}

//清空集合
func (sset *ESet) Clear() {
	sset.m.Clear()
	sset.em.Clear()
}

// Remove removes one or more elements from the list with the supplied indices.
func (sset *ESet) Remove(key interface{}) {
	sset.mutex_lock.Lock()
	sset.m.Remove(key)
	sset.em.Remove(key)
	sset.mutex_lock.Unlock()
}

// String returns a string representation of container
func (sset *ESet) String() string {
	str := "ESet\n" + sset.item_type + "\n"
	mStr := sset.m.String() + "\n"
	emStr := sset.em.String() + "\n"
	str = str + mStr + emStr
	return str
}

//按照切片索引获取item
func (sset *ESet) Get(key interface{}) (interface{}, bool) {
	ttl := sset.Ttl(key)
	if ttl > 0 || ttl == -1 {
		return sset.m.Get(key)
	}
	sset.Remove(key)
	return nil, false
}

func (sset *ESet) GetCreateline() int64 {
	return sset.createline
}

func (sset *ESet) Keys() (resKeys []interface{}) {
	keys := sset.m.Keys()
	nowTemp := time.Now().Unix()
	for _, key := range keys {
		expireTemp := sset.getExpire(key)
		if expireTemp > 0 && expireTemp < nowTemp {
			sset.Remove(key)
		} else {
			resKeys = append(resKeys)
		}
	}
	return
}

func (sset *ESet) Expire(key interface{}, expire int64) {
	if _, ok := sset.Get(key); ok {
		sset.mutex_lock.Lock()
		if expire > 0 {
			sset.em.Put(key, time.Now().Unix()+expire)
		} else {
			sset.em.Remove(key)
		}
		sset.mutex_lock.Unlock()
	}
}
func (sset *ESet) CronClear() {
	keys := sset.m.Keys()
	nowTemp := time.Now().Unix()
	for _, key := range keys {
		expireTemp := sset.getExpire(key)
		if expireTemp > 0 && expireTemp < nowTemp {
			sset.Remove(key)
		}
	}
	return
}

//TODO
func (sset *ESet) IncrBy(key interface{}, value int64) (int64, bool) {
	sset.mutex_lock.Lock()
	var int64Var int64
	if item, ok := sset.Get(key); ok {
		if int64Var, ok = item.(int64); !ok {
			sset.mutex_lock.Unlock()
			return 0, false
		}
	}
	sset.m.Put(key, int64Var+value)
	sset.mutex_lock.Unlock()
	return int64Var + value, true
}
