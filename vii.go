package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"
)

var (
	h    bool
	v    bool
	info bool
	p    string
)

func init() {
	flag.BoolVar(&h, "h", false, "帮助")
	flag.BoolVar(&v, "v", false, "版本")
	flag.BoolVar(&info, "i", false, "详细信息")
	flag.StringVar(&p, "p", "example_project", "项目名称 example_project")

	// 改变默认的 Usage
	flag.Usage = usage
}
func main() {
	flag.Parse()
	if len(os.Args) <= 1 || h {
		flag.Usage()
		os.Exit(0)
	}
	if v {
		fmt.Fprintf(os.Stdout, "1.10.0(beta)\n")
		os.Exit(0)
	}
	if p != "" {
		var stdout bytes.Buffer
		var stderr bytes.Buffer
		prj_path := os.Getenv("GOPATH")
		dst_file := prj_path + "/src/github.com/ztxmao/vii/example"
		if _, err := os.Open(dst_file); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		pArr := strings.Split(p, "/")
		plen := len(pArr)
		if plen > 1 {
			os.MkdirAll(p, 0755)
		}
		cmd1 := exec.Command("cp", "-rf", dst_file, p)
		cmd1.Stdout = &stdout
		cmd1.Stderr = &stderr
		if err := cmd1.Run(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		showinfo(stdout.String())

		stdout.Reset()
		stderr.Reset()
		cmd2 := exec.Command("grep", "{@project}", "-rl", p)
		cmd2.Stdout = &stdout
		cmd2.Stderr = &stderr
		if err := cmd2.Run(); err != nil {
			fmt.Println(stderr.String())
			os.Exit(2)
		}
		showinfo(stdout.String())
		lineArr := strings.Split(stdout.String(), "\n")
		var replaceContent string
		for _, line := range lineArr {
			if buf, err := ioutil.ReadFile(line); err == nil {
				content := string(buf)
				if path.Ext(line) == ".go" {
					replaceContent = p
				} else {
					replaceContent = strings.Join(pArr, "_")
				}
				newContent := strings.Replace(content, "{@project}", replaceContent, -1)
				ioutil.WriteFile(line, []byte(newContent), 0)
			}
		}
		stdout.Reset()
		stderr.Reset()
	}
	os.Exit(0)
}
func usage() {
	fmt.Fprintf(os.Stdout, `
vii version: 1.10.0(beta)
Usage: vii [-p projectname] 
Options
`)
	flag.PrintDefaults()
}

func showinfo(msg string) {
	if info {
		fmt.Println(msg)
	}
}
