package main

import (
	"fmt"
	"net/url"
	"time"
)

func savePage(qry, uri string) string {
	body := Download(uri)
	StorePage(qry, uri, body)
	return string(body)
}

func scheduler() {
	if !Configuration.DryRun {
		caseToken := []string{"cases", "hülle", "schutz"}
		models := getModels()

		for _, model := range models {
			for _, token := range caseToken {
				qry := fmt.Sprintf("%s %s", model[0], token)

				ParseAmazon(savePage(qry, qryAmazon(qry)))

				suggestions := savePage(qry, qryAmazonSuggestions(qry))
				decompSuggestions := decomposeSuggestions(qry, suggestions)
				for _, sug := range decompSuggestions {
					ParseAmazon(savePage(sug, qryAmazon(sug)))
				}
			}

			// No Google Trends in this version
			// for i := j + 1; i < len(models); i += 1 {
			// 	savePage(fmt.Sprintf("%s-%s", model[0], models[i][0]), qryGoogleTrend(model[0], models[i][0]))
			// }
		}
	} else {
		qry := "iphone 6s case"

		// One Iteration
		ParseAmazon(savePage(qry, qryAmazon(qry)))
		suggestions := savePage(qry, qryAmazonSuggestions(qry))
		decompSuggestions := decomposeSuggestions(qry, suggestions)
		body := savePage(decompSuggestions[0], qryAmazon(decompSuggestions[0]))
		ParseAmazon(body)
	}
}

func qryAmazon(str string) string {
	urlTemplate := `https://www.amazon.de/s/?field-keywords=%s`
	return fmt.Sprintf(urlTemplate, url.QueryEscape(str))
}

func qryAmazonSuggestions(str string) string {
	urlTemplate := `https://completion.amazon.co.uk/search/complete?method=completion&mkt=4&p=Search&l=de_DE&sv=desktop&client=amazon-search-ui&x=String&search-alias=aps&q=%s&qs=&cf=1&fb=1&sc=1`
	return fmt.Sprintf(urlTemplate, url.QueryEscape(str))
}

func qryGoogleTrend(a, b string) string {
	queryFormat := `{"time":"%d-%02d-%02dT00\\:01\\:00 %02d-%02d-%02dT00\\:01\\:00","resolution":"HOUR","locale":"de","comparisonItem":[{"geo":{"country":"DE"},"complexKeywordsRestriction":{"keyword":[{"type":"BROAD","value":"%s"}]}},{"geo":{"country":"DE"},"complexKeywordsRestriction":{"keyword":[{"type":"BROAD","value":"%s"}]}}],"requestOptions":{"property":"","backend":"CM","category":0}}`
	urlTemplate := `https://trends.google.com/trends/api/widgetdata/multiline/csv?req=%s&token=%s&tz=-60`

	fromDate := time.Now()
	fromDate = fromDate.AddDate(0, 0, -7)
	ny, nm, nd := time.Now().Date()
	by, bm, bd := fromDate.Date()

	qryJson := fmt.Sprintf(queryFormat, by, bm, bd, ny, nm, nd, a, b)

	return fmt.Sprintf(urlTemplate, url.QueryEscape(qryJson), "APP6_UEAAAAAWK-Jk8lNDN4yfRMYEaXhLkPD0J3U0ugC")
}
