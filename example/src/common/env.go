package common

import (
	"github.com/ztxmao/vii/library"
	"os"
)

var (
	Configer    = library.Configer
	Logger      = library.Logger
	Hostname, _ = os.Hostname()
)
