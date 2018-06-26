package common

import (
	"github.com/ztxmao/vii/library"
	"os"
)

var (
	Configer    = library.ConfigerExt
	Logger      = library.Logger
	Hostname, _ = os.Hostname()
)
