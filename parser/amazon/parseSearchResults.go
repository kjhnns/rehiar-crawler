package parser

import (
	"fmt"
	"regexp"
)

func ParseSearchResults(body string) {
	var re = regexp.MustCompile(`(von ((\d|\.)+) Ergebnissen oder Vorschlägen für|of ((\d|\.)+) results for)`)

	result := re.FindString(body)
	fmt.Println(result)
}
