package library

/**
	工具包
	@author ztxmao
**/

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"reflect"
	"runtime"
	"runtime/pprof"
	"sort"
	"strconv"
	"strings"
	"time"
)

//获取本机ip
func GetLocalIp() string {
	addrs, _ := net.InterfaceAddrs()
	var ip string = ""
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok {
			ip = ipnet.IP.String()
			if ip != "127.0.0.1" {
			}
		}
	}
	return ip
}

func GetTimeFormate() string {
	return "2006-01-02 15:04:05"
}

func GetTimeNow() string {
	return time.Now().Format(GetTimeFormate())
}

//简单序列化成php的格式, 和已有的php系统交互的时候可能会用到
func SerializePhp(data map[string]interface{}) string {
	ret := fmt.Sprintf("a:%d:{", len(data))
	for key, value := range data {
		ret = ret + fmt.Sprintf("s:%d:\"%s\";", len(key), key)
		if valuemap, ok := value.(map[string]interface{}); ok {
			ret = ret + SerializePhp(valuemap)
		} else {
			valuestr := value.(string)
			ret = ret + fmt.Sprintf("s:%d:\"%s\";", len(valuestr), valuestr)
		}
	}
	ret = ret + "}"
	return ret
}

//根据键排序map
func sortMapByKey(m map[string]string) ([]string, map[string]string) {
	var mk []string
	for k, _ := range m {
		mk = append(mk, k)
	}

	sort.Strings(mk)
	ret := make(map[string]string)
	for _, value := range mk {
		rkey, _ := m[value]
		ret[value] = rkey
	}

	return mk, ret
}

//生成32位md5字串
func GetMd5String(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

/**
*	json 字符串转换 对象
*	return interface;
 */
func Json2Obj(j string) (m interface{}) {
	if j == "" {
		m = nil
		return
	}
	err := json.NewDecoder(strings.NewReader(j)).Decode(&m)
	if err != nil {
		m = nil
	}
	return
}

/**
*	字符串转bytes
*	return targets || 用 []byte( conv ) 代替
 */
func String2ByteArray(conv string) (targets []byte) {
	if conv == "" {
		targets = nil
	}
	targets = bytes.NewBufferString(conv).Bytes()
	return
}

/**
*	Unicode转utf8
*	return string
 */
func Unicode2Utf8(str string) (rs string) {
	r := []rune(str)
	for i := 0; i < len(r); i++ {
		rs += string(r[i])
	}
	return rs
}

/**
*	map,string 转json对象
*	return string
 */
func Obj2Json(v interface{}) (res string) {
	b, err := json.Marshal(v)
	if err != nil {
		res = ""
	} else {
		res = string(b)
	}
	return
}

/**
 * intstr to int
 *
 */
func Str2Int(intStr string) int {
	/*{{{*/
	i, err := strconv.ParseInt(intStr, 10, 0)
	if err != nil {
		return 0
	}

	return int(i)
} /*}}}*/

/**
 * 将数字字符串转换成时间类型的ms
 *
 */
func MsStr2Duration(ms string) time.Duration {
	/*{{{*/
	l := len(ms)
	if ms[l-2] != 'm' && ms[l-1] != 's' {
		ms += "ms"
	}

	dura, err := time.ParseDuration(ms)
	if err != nil {
		return 500 * time.Millisecond
	}

	return dura
} /*}}}*/

/**
 * 开始记录cpu使用情况
 *
 */
func StartCPUProfile(cpuProfile string) error {
	/*{{{*/
	//100us 获取一次记录
	runtime.SetCPUProfileRate(10000)
	if f, err := os.Create(cpuProfile); err == nil {
		if err := pprof.StartCPUProfile(f); err != nil {
			f.Close()
			return err
		} else {
			return nil
		}
	} else {
		return err
	}

} /*}}}*/

/**
 * 停止记录cpu使用情况
 *
 */
func StopCPUProfile() {
	pprof.StopCPUProfile()
}

/**
 * redis HMSET 使用
 *
 */
func Struct2Kv(obj interface{}) (data []interface{}) {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)
	data = make([]interface{}, 2)
	for i := 0; i < t.NumField(); i++ {
		data = append(data, []interface{}{t.Field(i).Name, v.Field(i).Interface()})
	}
	return
}
