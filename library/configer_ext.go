package library

import (
	"bufio"
	//	"fmt"
	"errors"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

type ConfigExt struct {
	filepath string //your ini file path directory+file
	env      string
	conflist map[string]map[string]map[string]string //configuration information slice
}

var ConfigerExt = &ConfigExt{}

func (c *ConfigExt) Init(filepath, env string) {
	if c.filepath != "" {
		return
	}
	c.filepath = filepath
	c.env = env
	c.ReadList()
	if _, ok := c.conflist[env]; !ok {
		panic(errors.New("env node [" + env + "] not exist"))
	}
}

//Get string value of the key values
func (c *ConfigExt) GetStr(section, name string) string {
	sec, ok := c.conflist[c.env][section][name]
	if ok {
		return sec
	}
	return ""
}

//Get int value of the key values
func (c *ConfigExt) GetInt(section, name string) int {
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
func (c *ConfigExt) GetInt64(section, name string) int64 {
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
func (c *ConfigExt) ReadList() map[string]map[string]map[string]string {
	file, err := os.Open(c.filepath)
	if err != nil {
		c.checkErr(err)
	}
	defer file.Close()
	var block, blockKey, section string
	c.conflist = make(map[string]map[string]map[string]string)
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
		case line[0] == '{' && line[len(line)-1] == '}':
			blockKey = strings.TrimSpace(line[1 : len(line)-1])
			blockFields := strings.Split(blockKey, ":")
			block = blockFields[0]
			c.conflist[block] = make(map[string]map[string]string)
			//child node  copy parents val
			if len(blockFields) == 2 {
				parents := blockFields[1]
				if _, ok := c.conflist[parents]; ok {
					for psection, vals := range c.conflist[parents] {
						c.conflist[block][psection] = make(map[string]string)
						for field, val := range vals {
							c.conflist[block][psection][field] = val
						}
					}
				} else {
					panic(errors.New("parents [" + parents + "] node conf not exits"))
				}
			}
		case line[0] == '[' && line[len(line)-1] == ']':
			section = strings.TrimSpace(line[1 : len(line)-1])
			if _, ok := c.conflist[block][section]; !ok {
				c.conflist[block][section] = make(map[string]string)
			}
		case line[0] == '#':
		case line[0] == ';':
		default:
			i := strings.IndexAny(line, "=")
			value := strings.TrimSpace(line[i+1 : len(line)])
			c.conflist[block][section][strings.TrimSpace(line[0:i])] = value
		}
	}
	return c.conflist
}

//Get int value of the key values
func (c *ConfigExt) GetBool(section, name string) bool {
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
func (c *ConfigExt) GetSection(section string) (map[string]string, bool) {
	sec, ok := c.conflist[c.env][section]
	return sec, ok
}
func (c *ConfigExt) GetDuration(section, name string) time.Duration {
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
func (c *ConfigExt) checkErr(err error) {
	if err != nil {
		panic("Error is :" + err.Error())
	}
}
