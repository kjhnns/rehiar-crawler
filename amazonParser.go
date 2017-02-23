package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const insertAmazonData = `INSERT INTO amazonData (searchQuery, searchResults, timestamp) VALUES ($1,$2,$3) RETURNING id`

const amazonDataTable = `CREATE TABLE IF NOT EXISTS amazonData (id  SERIAL PRIMARY KEY, timestamp   timestamp DEFAULT current_timestamp, searchQuery   VARCHAR(255) NOT NULL,
searchResults   INT NOT NULL)`

func ParseAmazon(body string) {

	var parseSearchResults = func(body string) int {
		var re = regexp.MustCompile(`(von ((\d|\.)+) Ergebnissen oder Vorschlägen für|of ((\d|\.)+) results for)`)

		result := re.FindString(body)
		result = strings.TrimPrefix(result, "von ")
		result = strings.TrimPrefix(result, "of ")
		result = strings.TrimSuffix(result, " results for")
		result = strings.TrimSuffix(result, " Ergebnissen oder Vorschlägen für")

		result = strings.Replace(result, ".", "", -1)
		resInt, _ := strconv.Atoi(result)
		Configuration.Logger.Info.Println("Search Results: ", resInt)
		return resInt
	}

	var parseQryTerm = func(body string) string {
		var re = regexp.MustCompile(`&#034;(\w|\s)+&#034;`)

		result := re.FindString(body)
		result = strings.TrimPrefix(result, "&#034;")
		result = strings.TrimSuffix(result, "&#034;")
		Configuration.Logger.Info.Println(result)

		return result
	}

	var storeInDatabase = func(qry string, results int) {
		var appId int

		tsp := fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d-00:00",
			Configuration.StartTime.Year(), Configuration.StartTime.Month(), Configuration.StartTime.Day(),
			Configuration.StartTime.Hour(), Configuration.StartTime.Minute(), Configuration.StartTime.Second())
		Configuration.Logger.Info.Println(tsp)
		Configuration.Logger.Info.Println(Configuration.StartTime.Format(time.RFC3339))

		err := DbConn().QueryRow(insertAmazonData, qry, results, Configuration.StartTime.Format(time.RFC3339)).Scan(&appId)
		if err != nil {
			Configuration.Logger.Error.Println("Failed to save searchresults ", err)
		} else {
			Configuration.Logger.Info.Println("wrote to database ", appId)
		}
	}

	func(body string) {
		searchResults := parseSearchResults(body)
		qry := parseQryTerm(body)
		storeInDatabase(qry, searchResults)
	}(body)

}
