package main

import (
	"fmt"
	"os"
)

func Debug(debug Debugging) bool {
	return ((Configuration.Debugging & debug) > 0)
}

func Shutdown() {
	fmt.Printf("\n\n-> success, stopped")
	os.Exit(0)
}

func HandleErr(err ...interface{}) {
	if Debug(DEBUG_VERBOSE) {
		fmt.Println(err)
	}
	os.Exit(3)
}

func HandleDebugging(msg ...interface{}) {
	if Debug(DEBUG_DEBUGGING) {
		fmt.Println(msg)
	}
}
func HandleVerbose(msg ...interface{}) {
	if Debug(DEBUG_VERBOSE) {
		fmt.Println(msg)
	}
}

func HandleRepair(msg ...interface{}) {
	if Debug(DEBUG_REPAIR) {
		fmt.Println(msg)
	}
}

func HandleWarn(err ...interface{}) {
	if Debug(DEBUG_VERBOSE) {
		fmt.Println(err)
	}
}
