package main

import (
	"flag"
	"fmt"
	"os"
	"time"
)

func main() {
	fmt.Println("Load Configuration ...")
	loadConfig()
	fmt.Println("done.")
	SendMail("Hallo", "Test")

	// fmt.Println("Starting Crawler ... ")
	// scheduler()
	// fmt.Println("done.")

	// fmt.Println("Starting Crawler ... ")
}

func loadConfig() {
	// Default Configuration
	Configuration = Settings{
		StartTime: time.Now(),
	}

	dry := flag.Bool("dry", false, "dry run")
	flag.Parse()
	Configuration.DryRun = *dry

	if os.Getenv("DATABASE_URL") != "" {
		Configuration.DatabaseUrl = os.Getenv("DATABASE_URL")
	}

	if *dry {
		fmt.Println("DRY RUN")
		InitLogger("\t", os.Stdout, os.Stdout, os.Stderr)
	} else {
		y, m, d := Configuration.StartTime.Date()
		hr := Configuration.StartTime.Hour()
		mn := Configuration.StartTime.Minute()
		fileName := fmt.Sprintf("%d%02d%02d%02d%02d", y, m, d, hr, mn)
		file := "./logs/" + fileName

		logfile, err := os.Create(file)
		if err != nil {
			fmt.Println("Couldn't create a log file")
			InitLogger("\t", os.Stdout, os.Stdout, os.Stderr)
		} else {
			fmt.Sprintf("created log file - %s\n", file)
			InitLogger("", logfile, logfile, logfile)
		}
	}

}
