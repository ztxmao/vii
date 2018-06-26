package library

import (
	"bufio"
	//	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	filepath string                       //your ini file path directory+file
	conflist map[string]map[string]string //configuration information slice
}

var Configer = &Config{}

func (this *Config) Init(filepath string) {
	if this.filepath != "" {
		return
	}
	this.filepath = filepath
	this.ReadList()
}

//Get string value of the key values
func (c *Config) GetStr(section, name string) string {
	sec, ok := c.conflist[section][name]
	if ok {
		return sec
	}
	return ""
}

//Get int value of the key values
func (c *Config) GetInt(section, name string) int {
	str := c.GetStr(section, name)
	if str == "" {
		return 0
	}
	i, err := strconv.ParseInt(str, 10, 0)
	if err != nil {
		return 0
	}
	return int(i)
}

//Get int value of the key values
func (c *Config) GetInt64(section, name string) int64 {
	str := c.GetStr(section, name)

	if str == "" {
		return 0
	}
	result, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return 0
	}
	return result
}

//List all the configuration file
func (c *Config) ReadList() map[string]map[string]string {

	file, err := os.Open(c.filepath)
	if err != nil {
		c.checkErr(err)
	}
	defer file.Close()

	var section string
	c.conflist = make(map[string]map[string]string)

	buf := bufio.NewReader(file)
	for {
		l, err := buf.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				c.checkErr(err)
			}
			break
		}

		line := strings.TrimSpace(l)
		if len(line) == 0 {
			continue
		}

		switch {
		case len(line) == 0:
		case line[0] == '[' && line[len(line)-1] == ']':
			section = strings.TrimSpace(line[1 : len(line)-1])
			c.conflist[section] = make(map[string]string)
		case line[0] == '#':

		default:
			i := strings.IndexAny(line, "=")
			value := strings.TrimSpace(line[i+1 : len(line)])
			c.conflist[section][strings.TrimSpace(line[0:i])] = value
		}

	}

	return c.conflist
}

//Get int value of the key values
func (c *Config) GetBool(section, name string) bool {
	str := c.GetStr(section, name)
	if str == "" {
		return false
	}
	result, err := strconv.ParseBool(str)
	if err != nil {
		return false
	}
	return result
}

//Get int value of the key values
func (c *Config) GetSection(section string) (map[string]string, bool) {
	sec, ok := c.conflist[section]
	return sec, ok
}

func (c *Config) GetDuration(section, name string) time.Duration {
	str := c.GetStr(section, name)
	if str == "" {
		return 0
	}
	result, err := time.ParseDuration(str)
	if err != nil {
		return 0
	}
	return result
}
func (c *Config) checkErr(err error) {
	if err != nil {
		panic("Error is :" + err.Error())
	}
}
