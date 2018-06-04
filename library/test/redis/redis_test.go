package redis_test

import (
	"fmt"
	l "github.com/ztxmao/vii/library"
	"testing"
	//"time"

	"github.com/garyburd/redigo/redis"
)

var (
	p = fmt.Println
)

func makeRedisOption() *l.RedisOption {
	/* {{{*/
	serverPort := "10.16.57.78:6392"
	passwd := "719c7b2735318152"
	dbIndex := "0"
	connectTimeout := "50ms"
	readTimeout := "30ms"
	writeTimeout := "30ms"
	poolIdleTimeout := "500ms"
	poolMaxActiveLink := "50"
	poolMaxIdleLink := "20"
	waitConn := true

	ro := l.NewRedisOption()
	ro.Server = serverPort
	ro.Passwd = passwd
	ro.DbIndex = l.Str2Int(dbIndex)
	ro.ConnectTimeout = l.MsStr2Duration(connectTimeout)
	ro.ReadTimeout = l.MsStr2Duration(readTimeout)
	ro.WriteTimeout = l.MsStr2Duration(writeTimeout)
	ro.PoolMaxActive = l.Str2Int(poolMaxActiveLink)
	ro.PoolMaxIdle = l.Str2Int(poolMaxIdleLink)
	ro.PoolIdleTimeout = l.MsStr2Duration(poolIdleTimeout)
	ro.WaitConn = waitConn

	return ro
} /*}}}*/

func TestRedis(t *testing.T) {
	/*{{{*/
	rro := makeRedisOption()
	wro := makeRedisOption()
	r := l.NewRedisExecutor(rro, wro)

	fmt.Printf("redis executor: %#v\n", r)

	rconn, err := r.GetReadConn()
	if err != nil {
		panic(err.Error())
	} else {
		defer rconn.Close()
	}
	k := "fqTestKey1"
	v := "valtest"
	expire := 60
	res, err := redis.String(rconn.Do("GET", k))
	connCount := r.GetReadActiveConnCount()
	rconn.Close()
	fmt.Printf("raw get: res=%#v, err=%s, count=%d\n", res, err, connCount)

	wconn, err := r.GetReadConn()
	res2, err2 := wconn.Do("SET", k, v)
	connCount2 := r.GetWriteActiveConnCount()
	wconn.Close()
	fmt.Printf("raw set: res=%#v, err=%s, count=%d\n", res2, err2, connCount2)

	res3, err3 := r.Get(k)
	fmt.Printf("method get: res=%#v, err=%s\n", res3, err3)
	res6, err6 := r.GetString(k)
	fmt.Printf("method getstring: res=%#v, err=%s\n", res6, err6)
	var k2 string = "k2"
	var k3 string = "k3"
	res7, err7 := r.MgetString(k, k2, k3)
	fmt.Printf("method mget: res=%#v, err=%s\n", res7, err7)

	res4, err4 := r.Set(k2, v)
	fmt.Printf("method set: res=%#v, err=%s\n", res4, err4)

	res5, err5 := r.Setex(k, expire, v)
	fmt.Printf("method setex: res=%#v, err=%s\n", res5, err5)

} /*}}}*/
