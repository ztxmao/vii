package logic

import (
	"fmt"
)

//所有logic需要初始换的东西都写到这里
//静态类
var (
	msg string
)

// logic 所有初始化函数在这里执行，初始化失败 可以直接panic
func Init() {
	msg = "logic static obj init succ!"
	fmt.Println(msg)
}
