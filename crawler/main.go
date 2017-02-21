package main

import (
	"fmt"
	"os"
	"time"
)

func main() {
	fmt.Println("Load Configuration ...")
	loadConfig()
	fmt.Println("done.")
	scheduler()
}

func loadConfig() {
	// Default Configuration
	Configuration = Settings{
		StartTime: time.Now(),
	}

	y, m, d := Configuration.StartTime.Date()
	hr := Configuration.StartTime.Hour()
	mn := Configuration.StartTime.Minute()
	fileName := fmt.Sprintf("%d%02d%02d%02d%02d", y, m, d, hr, mn)
	file := "./logs/" + fileName

	logfile, err := os.Create(file)
	if err != nil {
		fmt.Println("Couldn't create a log file")
		InitLogger(os.Stdout, os.Stdout, os.Stderr)
	} else {
		InitLogger(logfile, logfile, logfile)
	}

}
