package main

import (
	"flag"
	"fmt"
	"time"
)

func main() {
	fmt.Print("Load Configuration ...")
	loadConfig()
	fmt.Println(" done.")

	ParseHtml()
}

func loadConfig() {
	// Default Configuration
	Configuration = Settings{
		StartTime: time.Now(),
		Debugging: DEBUG_NONE,
	}

	v := flag.Bool("v", false, "Verbose Output for Debugging purpose")
	vv := flag.Bool("vv", false, "Verbose Output for Debugging purpose")
	vvv := flag.Bool("vvv", false, "Verbose Output for Debugging purpose")
	flag.Parse()

	if *v {
		Configuration.Debugging = Configuration.Debugging | DEBUG_VERBOSE
	}

	if *vv {
		Configuration.Debugging = Configuration.Debugging | DEBUG_VERBOSE
		Configuration.Debugging = Configuration.Debugging | DEBUG_DEBUGGING
	}

	if *vvv {
		Configuration.Debugging = Configuration.Debugging | DEBUG_VERBOSE
		Configuration.Debugging = Configuration.Debugging | DEBUG_DEBUGGING
		Configuration.Debugging = Configuration.Debugging | DEBUG_REPAIR
	}

}
