package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"net/url"
	"os"
	"time"
)

func scheduler() {
	models := getModels()

	for j, model := range models {
		storePage(amazon(model[0]))

		for i := j + 1; i < len(models); i += 1 {
			storePage(googleTrend(model[0], models[i][0]))
		}
	}
}

func amazon(str string) string {
	urlTemplate := `https://www.amazon.de/s/?field-keywords=%s`

	return fmt.Sprintf(urlTemplate, url.QueryEscape(str))
}

func googleTrend(a, b string) string {
	queryFormat := `{"time":"%d-%02d-%02dT00\\:01\\:00 %02d-%02d-%02dT00\\:01\\:00","resolution":"HOUR","locale":"de","comparisonItem":[{"geo":{"country":"DE"},"complexKeywordsRestriction":{"keyword":[{"type":"BROAD","value":"%s"}]}},{"geo":{"country":"DE"},"complexKeywordsRestriction":{"keyword":[{"type":"BROAD","value":"%s"}]}}],"requestOptions":{"property":"","backend":"CM","category":0}}`
	urlTemplate := `https://trends.google.com/trends/api/widgetdata/multiline/csv?req=%s&token=%s&tz=-60`

	fromDate := time.Now()
	fromDate = fromDate.AddDate(0, 0, -7)
	ny, nm, nd := time.Now().Date()
	by, bm, bd := fromDate.Date()

	qryJson := fmt.Sprintf(queryFormat, by, bm, bd, ny, nm, nd, a, b)

	return fmt.Sprintf(urlTemplate, url.QueryEscape(qryJson), "APP6_UEAAAAAWKzj62TR8fNVzg7-YUfOPQCyk7o8Zgsd")
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
