package main

import (
	"regexp"
	"strconv"
	"strings"
	"time"
)

const dropAmazonData = `drop table amazonData`
const insertAmazonData = `INSERT INTO amazonData (searchQuery, searchResults, timestamp, caseterm,themeterm, modelterm) VALUES ($1,$2,$3,$4,$5,$6) RETURNING id`

const amazonDataTable = `CREATE TABLE IF NOT EXISTS amazonData (
	id  SERIAL PRIMARY KEY,
	timestamp   timestamp DEFAULT current_timestamp,
	searchQuery   VARCHAR(255) NOT NULL,
	caseterm   VARCHAR(255) NOT NULL,
	themeterm   VARCHAR(255) NOT NULL,
	modelterm   VARCHAR(255) NOT NULL,
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

	var parseModel = func(str string) string {
		models := getModels()
		for _, model := range models {
			if strings.Contains(str, model[0]) {
				return model[0]
			}
		}
		return ""
	}

	var parseTheme = func(str string) string {
		cases := []string{"cases", "hülle", "schutz", "case"}

		for _, token := range cases {
			if strings.Contains(str, token) {
				return token
			}
		}
		return ""
	}

	var parseCase = func(str string) string {
		trimmedModel := strings.TrimPrefix(str, parseModel(str)+" ")
		caseqry := strings.TrimPrefix(trimmedModel, parseTheme(trimmedModel)+" ")

		return caseqry
	}

	var parseQryTerm = func(body string) string {
		var re = regexp.MustCompile(`Ergebnissen oder Vorschlägen für <span><span class="a-color-state a-text-bold">&#034;([\w\süäö])+&#034;<\/span>`)

		result := re.FindString(body)
		result = strings.TrimPrefix(result, "Ergebnissen oder Vorschlägen für <span><span class=\"a-color-state a-text-bold\">&#034;")
		result = strings.TrimSuffix(result, "&#034;</span>")
		Configuration.Logger.Info.Println(result)

		return result
	}

	var storeInDatabase = func(qry string, results int) {
		var appId int

		err := DbConn().QueryRow(insertAmazonData, qry, results, Configuration.StartTime.Format(time.RFC3339), parseTheme(qry), parseCase(qry), parseModel(qry)).Scan(&appId)
		if err != nil {
			Configuration.Logger.Warning.Println("Failed to save searchresults ", err)
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
