package library

import (
	//"fmt"
	//"reflect"
	"time"

	"github.com/garyburd/redigo/redis"
)

//Create an empty configuration file
func NewRedisExecutor(readRedisOption, writeRedisOption *RedisOption) *RedisExecutor {
	/*{{{*/
	r := new(RedisExecutor)
	r.Init(readRedisOption, writeRedisOption)

	return r
} /*}}}*/

func NewRedisOption() *RedisOption {
	/*{{{*/
	ro := &RedisOption{}
	ro.Network = "tcp"
	ro.Server = "127.0.0.1:6379"
	ro.Passwd = ""
	ro.DbIndex = 0
	ro.ConnectTimeout = time.Millisecond * 20
	ro.ReadTimeout = time.Millisecond * 30
	ro.WriteTimeout = time.Millisecond * 30
	ro.PoolMaxActive = 50
	ro.PoolMaxIdle = 30
	ro.PoolIdleTimeout = time.Millisecond * 1000
	ro.WaitConn = true

	return ro
} /*}}}*/

//创建redis链接所需的参数
type RedisOption struct {
	/*{{{*/
	Network         string        `协议`
	Server          string        `ip:port`
	Passwd          string        `redis auth password`
	DbIndex         int           `db index default 0`
	ConnectTimeout  time.Duration `connection timeout`
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	PoolMaxActive   int           `最大活跃链接数`
	PoolMaxIdle     int           `最大空闲链接数`
	PoolIdleTimeout time.Duration `链接的空闲时长`
	WaitConn        bool          `链接池耗尽时是否等着链接`
} /*}}}*/

type RedisExecutor struct {
	/*{{{*/
	readPool         *redis.Pool
	writePool        *redis.Pool
	readRedisOption  *RedisOption
	writeRedisOption *RedisOption
} /*}}}*/

//初始化
func (this *RedisExecutor) Init(readRedisOption, writeRedisOption *RedisOption) {
	/*{{{*/
	this.readRedisOption = readRedisOption
	this.writeRedisOption = writeRedisOption
	this.readPool = this.NewPool(this.readRedisOption)
	this.writePool = this.NewPool(this.writeRedisOption)
} /*}}}*/

//创建redis链接池
func (this *RedisExecutor) NewPool(ro *RedisOption) *redis.Pool {
	/*{{{*/
	return &redis.Pool{
		MaxActive:   ro.PoolMaxActive,
		MaxIdle:     ro.PoolMaxIdle,
		IdleTimeout: ro.PoolIdleTimeout,
		Wait:        ro.WaitConn,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial(ro.Network,
				ro.Server,
				redis.DialConnectTimeout(ro.ConnectTimeout),
				redis.DialReadTimeout(ro.ReadTimeout),
				redis.DialWriteTimeout(ro.WriteTimeout),
				redis.DialDatabase(ro.DbIndex),
				redis.DialPassword(ro.Passwd))
			if err != nil {
				return nil, err
			}

			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}
} /*}}}*/

//获取read链接
func (this *RedisExecutor) GetReadConn() (redis.Conn, error) {
	/*{{{*/
	conn := this.readPool.Get()
	if err := conn.Err(); err != nil {
		return nil, err
	}
	return conn, nil

} /*}}}*/

//获取write链接
func (this *RedisExecutor) GetWriteConn() (redis.Conn, error) {
	/*{{{*/
	conn := this.writePool.Get()
	if err := conn.Err(); err != nil {
		return nil, err
	}
	return conn, nil

} /*}}}*/

//read链接池关闭链接
func (this *RedisExecutor) ReadPoolClose() error {
	/*{{{*/
	return this.readPool.Close()
} /*}}}*/

//write链接池关闭链接
func (this *RedisExecutor) WritePoolClose() error {
	/*{{{*/
	return this.writePool.Close()
} /*}}}*/

//read链接池中可用链接数
func (this *RedisExecutor) GetReadActiveConnCount() int {
	/*{{{*/
	return this.readPool.ActiveCount()
} /*}}}*/

//write链接池中可用链接数
func (this *RedisExecutor) GetWriteActiveConnCount() int {
	/*{{{*/
	return this.writePool.ActiveCount()
} /*}}}*/

//------------- redis 操作 -----------------
//从redis中取数据
func (this *RedisExecutor) Get(key string) (interface{}, error) {
	/*{{{*/
	conn, err := this.GetReadConn()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	data, err := conn.Do("GET", key)
	return data, err
} /*}}}*/

//从redis中取数据,并转成string类型
func (this *RedisExecutor) GetString(key string) (string, error) {
	/*{{{*/
	data, err := redis.String(this.Get(key))
	if err != nil {
		return "", err
	}
	return data, err
} /*}}}*/

//从redis中取数据,并转成int类型
func (this *RedisExecutor) GetInt(key string) (int, error) {
	/*{{{*/
	data, err := redis.Int(this.Get(key))
	if err != nil {
		return 0, err
	}
	return data, err
} /*}}}*/

//从redis中取出多个key的数据
func (this *RedisExecutor) Mget(key string, otherKey ...interface{}) (interface{}, error) {
	/*{{{*/
	conn, err := this.GetReadConn()
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	keys := make([]interface{}, len(otherKey)+1)
	keys[0] = key
	for i := range otherKey {
		keys[i+1] = otherKey[i]
	}

	data, err := conn.Do("MGET", keys...)
	return data, err
} /*}}}*/

//从redis中取出多个key的数据,并转成string类型
func (this *RedisExecutor) MgetString(key string, otherKey ...interface{}) ([]string, error) {
	/*{{{*/
	data, err := redis.Strings(this.Mget(key, otherKey...))
	if err != nil {
		return nil, err
	}
	return data, err
} /*}}}*/

//存数据,有失效时间
func (this *RedisExecutor) Setex(key string, seconds int, value interface{}) (interface{}, error) {
	/*{{{*/
	conn, err := this.GetWriteConn()
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	return conn.Do("SETEX", key, seconds, value)
} /*}}}*/

//存数据,永久存储
func (this *RedisExecutor) Set(key string, value interface{}) (interface{}, error) {
	/*{{{*/
	conn, err := this.GetWriteConn()
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	return conn.Do("SET", key, value)
} /*}}}*/

func (this *RedisExecutor) Mrdo(cmd string, args ...interface{}) (interface{}, error) {
	// {{{
	conn, err := this.GetReadConn()
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	//fmt.Println(args)
	data, err := conn.Do(cmd, args...)

	return data, err
} // }}}

func (this *RedisExecutor) Wdo(cmd string, args ...interface{}) (interface{}, error) {
	// {{{
	conn, err := this.GetWriteConn()
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	return conn.Do(cmd, args...)
} // }}}
