package models

import (
	"fmt"
)

var (
	msg string
)

// model all init function
func Init() {
	msg = "models init succ!"
	fmt.Println(msg)
}
