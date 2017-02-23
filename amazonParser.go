package main

import (
	"regexp"
	"strconv"
	"strings"
)

const insertAmazonData = `INSERT INTO amazonData (searchQuery, searchResults) VALUES ($1,$2) RETURNING id`

const amazonDataTable = `CREATE TABLE IF NOT EXISTS amazonData (id  SERIAL PRIMARY KEY, created_at   timestamp DEFAULT current_timestamp, searchQuery   VARCHAR(255) NOT NULL,
searchResults   INT NOT NULL)`

func ParseAmazon(qry, body string) {

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

	var storeInDatabase = func(qry string, results int) {
		var appId int
		err := DbConn().QueryRow(insertAmazonData, qry, results).Scan(&appId)
		if err != nil {
			Configuration.Logger.Error.Println("Failed to save searchresults ", err)
		} else {
			Configuration.Logger.Info.Println("wrote to database ", appId)
		}
	}

	func(qry, body string) {
		searchResults := parseSearchResults(body)
		storeInDatabase(qry, searchResults)
	}(qry, body)

}
