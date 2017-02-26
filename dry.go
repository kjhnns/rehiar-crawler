package main

import (
	"fmt"
)

func initDryMode() {
	qry := "iphone 6s"
	qry2 := "iphone 7"

	bodygt := savePage(fmt.Sprintf("%s-%s", qry, qry2), qryGoogleTrendCSV(qry, qry2))
	ParseGoogleTrends(bodygt)

	ParseAmazon(savePage(qry+" cases", qryAmazon(qry+" cases")))
	suggestions := savePage(qry, qryAmazonSuggestions(qry))
	decompSuggestions := decomposeSuggestions(qry, suggestions)
	bodya := savePage(decompSuggestions[0], qryAmazon(decompSuggestions[0]))
	ParseAmazon(bodya)
}
