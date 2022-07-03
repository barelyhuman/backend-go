package debug

import (
	"log"
	"os"
)

var debug bool

func init() {
	envDebug := os.Getenv("DEBUG")
	if len(envDebug) > 0 {
		debug = true
	}
}

func DebugErr(err error) {
	if err == nil {
		return
	}

	if !debug {
		return
	}

	log.Println(err)
}

func DebugLog(v interface{}) {
	if !debug {
		return
	}
	log.Println(v)
}
