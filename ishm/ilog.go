package ishm

import (
	"fmt"
)

// DEBUG decide print some debug info or not
var DEBUG bool

func init() {
	DEBUG = false
}

func iPrintln(v ...interface{}) {
	if DEBUG {
		fmt.Println(v...)
	}
}

func iPrintf(format string, v ...interface{}) {
	if DEBUG {
		fmt.Printf(format, v...)
	}
}
