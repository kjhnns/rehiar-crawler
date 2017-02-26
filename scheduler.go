package main

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

func savePage(qry, uri string) string {
	body := Download(uri)
	StorePage(qry, uri, body)
	return string(body)
}

func scheduler() {
	caseToken := []string{"cases", "case", "hülle", "schutz"}
	models := getModels()

	for j, model := range models {
		for _, token := range caseToken {
			qry := fmt.Sprintf("%s %s", model[0], token)

			modelData := savePage(qry, qryAmazon(qry))
			ParseAmazon(modelData)

			suggestions := savePage(qry, qryAmazonSuggestions(qry))
			decompSuggestions := decomposeSuggestions(qry, suggestions)

			for _, sug := range decompSuggestions {
				themeData := savePage(sug, qryAmazon(sug))
				ParseAmazon(themeData)
			}
		}

		for i := j + 1; i < len(models); i += 1 {
			qry := fmt.Sprintf("%s-%s", model[0], models[i][0])

			trendData := savePage(qry, qryGoogleTrendCSV(model[0], models[i][0]))
			ParseGoogleTrends(trendData)
		}
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

func qryGoogleTrendToken(a, b string) (token, tsp string) {

	queryFormat := `{"comparisonItem":[{"keyword":"%s","geo":"DE","time":"now 7-d"},{"keyword":"%s","geo":"DE","time":"now 7-d"}],"category":0,"property":""}`

	urlTemplate := `https://trends.google.com/trends/api/explore?hl=en-US&tz=-60&req=%s&tz=-60`

	qryJson := fmt.Sprintf(queryFormat, a, b)

	tokenResponse := string(Download(fmt.Sprintf(urlTemplate, url.QueryEscape(qryJson))))

	// Acquire Timestamp
	var reTsp = regexp.MustCompile(`"time":"([\d-T\\:]+) ([\d-T\\:]+)"`)
	resultTsp := reTsp.FindString(tokenResponse)
	resultTsp = strings.TrimPrefix(resultTsp, `"time":"`)
	resultTsp = strings.TrimSuffix(resultTsp, `"`)
	Configuration.Logger.Info.Println("Timestamp:", resultTsp)

	// Acquire Token
	var reToken = regexp.MustCompile(`"token":"(.*)","id":"TIMESERIES","type":"fe_line_chart","title":"Interest over time"`)
	resultToken := reToken.FindString(tokenResponse)
	resultToken = strings.TrimPrefix(resultToken, `"token":"`)
	resultToken = strings.TrimSuffix(resultToken, `","id":"TIMESERIES","type":"fe_line_chart","title":"Interest over time"`)
	Configuration.Logger.Info.Println("Token: ", resultToken)

	return resultToken, resultTsp
}

func qryGoogleTrendCSV(a, b string) string {
	queryFormat := `{"time":"%s","resolution":"HOUR","locale":"en-US","comparisonItem":[{"geo":{"country":"DE"},"complexKeywordsRestriction":{"keyword":[{"type":"BROAD","value":"%s"}]}},{"geo":{"country":"DE"},"complexKeywordsRestriction":{"keyword":[{"type":"BROAD","value":"%s"}]}}],"requestOptions":{"property":"","backend":"CM","category":0}}`
	urlTemplate := `https://trends.google.com/trends/api/widgetdata/multiline/csv?req=%s&token=%s&tz=-60`

	token, tsp := qryGoogleTrendToken(a, b)
	qryJson := fmt.Sprintf(queryFormat, tsp, a, b)

	requrl := fmt.Sprintf(urlTemplate, url.QueryEscape(qryJson), token)
	requrl = strings.Replace(requrl, "+", "%20", -1)

	return requrl

}
