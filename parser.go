package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"time"
)

func resetDatabase() {
	DbConn().Exec(dropAmazonData)
	DbConn().Exec(amazonDataTable)

	DbConn().Exec(dropGoogleTrends)
	DbConn().Exec(googleTrendsTable)
}

func initParseMode() {
	var iterateStorageFS = func(iterator func(path, folder, filename string)) {
		folders, err := ioutil.ReadDir("./data/")
		if err != nil {
			Configuration.Logger.Error.Println("Couldn't open data folder", err)
			return
		}

		for _, folder := range folders {
			if folder.Name() != ".DS_Store" {
				files, err_sf := ioutil.ReadDir("./data/" + folder.Name())
				if err_sf != nil {
					Configuration.Logger.Error.Println("Couldn't open data folder", "./data/"+folder.Name(), err_sf)
					return
				}

				for _, file := range files {
					if file.Name() != ".DS_Store" {
						path := "./data/" + folder.Name() + "/" + file.Name()
						iterator(path, folder.Name(), file.Name())
					}
				}
			}
		}
	}

	var parseMetadata = func(folder, filename string) (time.Time, string) {
		tstp := fmt.Sprintf("%s-%s-%sT%s:%s:00+00:00", folder[:4], folder[4:6], folder[6:8], folder[8:10], folder[10:12])
		fname := strings.Split(filename, "-")
		resTstp, err := time.Parse(time.RFC3339, tstp)
		if err != nil {
			Configuration.Logger.Warning.Println("Couldn't parse time", err)
		}
		return resTstp, fname[1]
	}

	var loadBody = func(path string) string {
		body, err := ioutil.ReadFile(path)
		if err != nil {
			Configuration.Logger.Warning.Println("Couldn't open ", path, err)
		}
		return string(body)
	}

	resetDatabase()

	iterateStorageFS(func(path, folder, filename string) {
		var domain string
		Configuration.StartTime, domain = parseMetadata(folder, filename)
		body := loadBody(path)

		Configuration.Logger.Info.Println("parsing ", filename)
		switch domain {
		case "www.amazon.de":
			ParseAmazon(body)
		case "trends.google.com":
			ParseGoogleTrends(string(body))
		}
	})

}

func decomposeSuggestions(qry, suggestions string) []string {
	var re = regexp.MustCompile(fmt.Sprintf(`"%s(\w|\s)+"`, qry))

	var matches []string
	for _, match := range re.FindAllString(suggestions, -1) {
		matches = append(matches, strings.Trim(match, "\""))
	}
	return matches
}

func getModels() [][]string {
	fileHandler, _ := os.Open("models.csv")
	defer fileHandler.Close()

	// Create a new reader.
	csvReader := csv.NewReader(bufio.NewReader(fileHandler))
	csvReader.Comma = ';'
	csvReader.Comment = '#'

	records, err := csvReader.ReadAll()
	if err != nil {
		Configuration.Logger.Error.Println("models err", err)
	}

	return records
}
