package main

import (
	"fmt"
	"golang.org/x/net/html"
	"io/ioutil"
	// "strconv"
	"os"
	"strings"
)
import "encoding/csv"

type Company struct {
	Id        string
	Name      string
	Url       string
	Zip       string
	Place     string
	Employees string
	Revenue   string
}

func (comp Company) Slice() []string {
	var res []string
	res = append(res, ""+comp.Id)
	res = append(res, ""+comp.Name)
	res = append(res, ""+comp.Url)
	res = append(res, ""+comp.Zip)
	res = append(res, ""+comp.Place)
	res = append(res, ""+comp.Employees)
	res = append(res, ""+comp.Revenue)

	return res
}

func ParseHtml() {
	var resultSlice [][]string

	resultSlice = append(resultSlice, Company{
		Id:        "Id",
		Name:      "Name",
		Url:       "Url",
		Zip:       "Zip",
		Place:     "Place",
		Employees: "Employees",
		Revenue:   "Revenue",
	}.Slice())

	fmt.Printf("Parsing Input ... ")

	body, _ := ioutil.ReadFile("./input.html")

	// Processing html
	doc, err := html.Parse(strings.NewReader(string(body)))
	if err != nil {
		HandleErr(err)
	}

	companiesTable := doc.FirstChild.FirstChild.NextSibling.FirstChild.NextSibling
	if companiesTable == nil {
		HandleErr("Failed to find CompanyTable")
	}
	companiesTable = companiesTable.FirstChild
	companiesTable = companiesTable.NextSibling

	count := 0
	for company := companiesTable.FirstChild; company != nil; company = company.NextSibling {
		if company.DataAtom.String() == "tr" {
			count += 1

			fmt.Println("Company", count)
			tmpCompany := Company{}

			row := 0
			for companyData := company.FirstChild; companyData != nil; companyData = companyData.NextSibling {
				if companyData.DataAtom.String() == "td" {
					row += 1
					if len(companyData.Attr) > 0 && (companyData.Attr[0].Val == "hlm" || companyData.Attr[0].Val == "hlc") {
						// fmt.Println(companyData, row)

						// Heres where we fetch the ID
						if row == 1 {
							tmpCompany.Id = strings.TrimLeft(companyData.FirstChild.Attr[0].Val, "hdb_trefferliste.cgi?select=")
						}
						// Heres where we fetch the Name and the Url
						if row == 3 {
							tmpCompany.Url = strings.TrimLeft(companyData.FirstChild.Attr[0].Val, "hdb_firma_xml.cgi?hoco=")
							tmpCompany.Name = companyData.FirstChild.Attr[1].Val
						}
						// Heres where we fetch the Name and the Url
						if row == 5 {
							tmpCompany.Zip = companyData.FirstChild.Data
						}
						if row == 6 {
							tmpCompany.Place = companyData.FirstChild.Data
						}
						if row == 7 {
							tmpCompany.Employees = companyData.FirstChild.Data
						}
						if row == 8 {
							tmpCompany.Revenue = companyData.FirstChild.Data
						}
					}
				}
			}
			// result.Companies = append(result.Companies, tmpCompany)
			// tmpCompany.FetchInformation()
			// fmt.Println(tmpCompany)
			resultSlice = append(resultSlice, tmpCompany.Slice())
		}
	}

	file, _ := os.Create("./out.csv")
	defer file.Close()

	writer := csv.NewWriter(file)
	writer.Comma = ';'
	writer.WriteAll(resultSlice)

	if err := writer.Error(); err != nil {
		HandleErr("error writing csv:", err)
	}

	fmt.Println("done.")

}
