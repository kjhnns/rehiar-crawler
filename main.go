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

	fmt.Println("\nStarting")
	initialize()
}

func initialize() {
	switch Configuration.Mode {
	case DRYMODE:
		Configuration.Logger.Info.Println("dry mode")
		initDryMode()
	case PARSEMODE:
		Configuration.Logger.Info.Println("parse mode")
		initParseMode()
	default:
		Configuration.Logger.Info.Println("initialize crawler")
		scheduler()
	}
}

func loadConfig() {
	// Default Configuration
	Configuration = Settings{
		StartTime: time.Now(),
		SleepTime: 10,
	}

	mode := flag.String("mode", "", "modes: dry, parse")
	sleeptime := flag.Int("sleep", -1, "set the spleeptime between http requests")
	flag.Parse()

	switch *mode {
	case "dry":
		Configuration.Mode = DRYMODE
	case "parse":
		Configuration.Mode = PARSEMODE
	default:
		Configuration.Mode = NORMALMODE
	}

	if os.Getenv("DATABASE_URL") != "" {
		Configuration.DatabaseUrl = os.Getenv("DATABASE_URL")
	}

	if os.Getenv("MAILSENDER") != "" {
		Configuration.Mail.Sender = os.Getenv("MAILSENDER")
	}
	if os.Getenv("MAILSERVER") != "" {
		Configuration.Mail.Server = os.Getenv("MAILSERVER")
	}
	if os.Getenv("MAILUSER") != "" {
		Configuration.Mail.User = os.Getenv("MAILUSER")
	}
	if os.Getenv("MAILPASS") != "" {
		Configuration.Mail.Pass = os.Getenv("MAILPASS")
	}
	if os.Getenv("MAILPORT") != "" {
		Configuration.Mail.Pass = os.Getenv("MAILPORT")
	}

	if Configuration.Mode != NORMALMODE {
		fmt.Println("\t- Logging to console")
		Configuration.SleepTime = 1
		InitLogger("\t", os.Stdout, os.Stdout, os.Stderr)
	} else {
		y, m, d := Configuration.StartTime.Date()
		hr := Configuration.StartTime.Hour()
		mn := Configuration.StartTime.Minute()
		fileName := fmt.Sprintf("%d%02d%02d%02d%02d", y, m, d, hr, mn)
		file := "./logs/" + fileName

		logfile, err := os.Create(file)
		if err != nil {
			fmt.Println("\t- Couldn't create a log file")
			InitLogger("\t", os.Stdout, os.Stdout, os.Stderr)
		} else {
			fmt.Printf("\t- created log file - %s\n", file)
			InitLogger("", logfile, logfile, logfile)
		}
	}

	if *sleeptime >= 0 {
		Configuration.SleepTime = *sleeptime
	}
	fmt.Println("\t- sleeptime ", Configuration.SleepTime)

	InitDatabase()
}
