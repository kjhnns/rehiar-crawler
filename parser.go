package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"regexp"
	"strings"
)

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
