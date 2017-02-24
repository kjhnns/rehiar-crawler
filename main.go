package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

func main() {
	fmt.Println("Load Configuration ...")
	loadConfig()
	fmt.Println("done.")

}

func startparse() {
	files, _ := ioutil.ReadDir("./data/")
	for _, folder := range files {
		if folder.Name() != ".DS_Store" {

			subfolder, _ := ioutil.ReadDir("./data/" + folder.Name())
			timestamp := folder.Name()
			for _, file := range subfolder {

				if file.Name() != ".DS_Store" {
					fileName := file.Name()

					tstp := fmt.Sprintf("%s-%s-%sT%s:%s:00+00:00", timestamp[:4], timestamp[4:6], timestamp[6:8], timestamp[8:10], timestamp[10:12])

					domain := strings.Split(fileName, "-")

					if "www.amazon.de" == domain[1] {
						fmt.Println("parsing ", fileName, tstp, domain[1])
						body, _ := ioutil.ReadFile("./data/" + timestamp + "/" + fileName)
						Configuration.StartTime, _ = time.Parse(time.RFC3339, tstp)
						fmt.Println(Configuration.StartTime)
						ParseAmazon(string(body))

					}
				}
			}
		}
	}

}

func loadConfig() {
	// Default Configuration
	Configuration = Settings{
		StartTime: time.Now(),
		SleepTime: 20,
	}

	parse := flag.Bool("parse", false, "take all the historic crawl data, truncate the database and refill")
	dry := flag.Bool("dry", false, "dry run")
	flag.Parse()
	Configuration.DryRun = *dry

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

	if *dry {
		fmt.Println("DRY RUN")
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
			fmt.Println("Couldn't create a log file")
			InitLogger("\t", os.Stdout, os.Stdout, os.Stderr)
		} else {
			fmt.Sprintf("created log file - %s\n", file)
			InitLogger("", logfile, logfile, logfile)
		}
	}

	InitDatabase()

	if *parse {

		DbConn().Exec(dropAmazonData)
		DbConn().Exec(amazonDataTable)

		fmt.Println("Starting Parser ... ")
		startparse()
		fmt.Println("done.")
	} else {
		fmt.Println("Starting Crawler ... ")
		scheduler()
		fmt.Println("done.")
	}

}
