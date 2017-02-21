package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"golang.org/x/net/html"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

var wzkeys []int
var wzTitles map[int]string

var wzItCodes = []int{61100, 61200, 61300, 61901, 61909, 62011, 62019, 62020, 62030, 62090, 63110, 63120, 63910, 63990}

type Company struct {
	Id        int
	Name      string
	Mail      string
	Zip       string
	Place     string
	Employees string
	Revenue   string
	Form      string
	Founded   string
	Job       string
	Locations string
	WZs       []int
}

func (comp Company) Slice() []string {
	var res []string
	res = append(res, ""+strconv.Itoa(comp.Id))
	res = append(res, ""+comp.Name)
	res = append(res, ""+comp.Mail)
	res = append(res, ""+comp.Zip)
	res = append(res, ""+comp.Place)
	res = append(res, ""+comp.Employees)
	res = append(res, ""+comp.Revenue)
	res = append(res, ""+comp.Form)
	res = append(res, ""+comp.Founded)
	res = append(res, ""+comp.Job)
	res = append(res, ""+comp.Locations)

	for k, _ := range wzkeys {
		found := false
		for vv := range comp.WZs {
			if vv == k {
				found = true
			}
		}

		if found {
			res = append(res, "1")
		} else {
			res = append(res, "0")
		}
	}

	return res
}

func ParseHtml() {
	var resultSlice [][]string

	var cmpIds map[int]int
	var companyStash []Company
	wzTitles = make(map[int]string)
	cmpIds = make(map[int]int)

	fmt.Printf("Parsing Input ... ")

	counting := 0
	files, _ := ioutil.ReadDir("./input/")
	for _, file := range files {
		if file.Name() != ".DS_Store" {
			counting += 1
			fmt.Println("./input/"+file.Name(), counting)
			body, _ := ioutil.ReadFile("./input/" + file.Name())

			// Processing html
			doc, err := html.Parse(strings.NewReader(string(body)))
			if err != nil {
				fmt.Println(err)
			}
			iId, _ := strconv.Atoi(file.Name())
			tmpCompany := Company{Id: iId}

			for table := doc.FirstChild.FirstChild.NextSibling.NextSibling.FirstChild; table != nil; table = table.NextSibling {
				if table.DataAtom.String() == "table" {
					companiesTable := table.FirstChild.NextSibling

					for detail := companiesTable.FirstChild; detail != nil; detail = detail.NextSibling {
						if detail.DataAtom.String() == "tr" {
							if detail.FirstChild != nil && detail.FirstChild.FirstChild != nil && detail.FirstChild.FirstChild.Data == "Branche WZ 2008 :" {

								for wzCode := detail.FirstChild.NextSibling.FirstChild; wzCode != nil; wzCode = wzCode.NextSibling {
									if wzCode.DataAtom.String() == "span" {
										iWz, _ := strconv.Atoi(wzCode.FirstChild.Data)
										wzTitles[iWz] = wzCode.Attr[0].Val
										tmpCompany.WZs = append(tmpCompany.WZs, iWz)
									}
								}
							}

							if detail.FirstChild != nil && detail.FirstChild.FirstChild != nil && detail.FirstChild.FirstChild.Data == "Gründung :" {
								tmpCompany.Founded = strings.TrimSpace(strings.TrimSuffix(detail.FirstChild.NextSibling.FirstChild.Data, "&#160;"))
							}

							if detail.FirstChild != nil && detail.FirstChild.FirstChild != nil && detail.FirstChild.FirstChild.Data == "Rechtsform :" {
								tmpCompany.Form = strings.TrimSpace(detail.FirstChild.NextSibling.FirstChild.Data)
							}

							if detail.FirstChild != nil && detail.FirstChild.FirstChild != nil && detail.FirstChild.FirstChild.Data == "Geschäftstätigkeit :" {
								tmpCompany.Job = detail.FirstChild.NextSibling.FirstChild.Data
							}
							if detail.FirstChild != nil && detail.FirstChild.FirstChild != nil && detail.FirstChild.FirstChild.Data == "Niederlassung(en)" {
								tmpCompany.Locations = strings.TrimSpace(detail.NextSibling.NextSibling.NextSibling.FirstChild.NextSibling.FirstChild.Data)
							}
						}
					}

				}
			}
			cmpIds[tmpCompany.Id] = len(companyStash)
			companyStash = append(companyStash, tmpCompany)

		}
	}

	hl := Company{
		Id:        000,
		Name:      "Name",
		Mail:      "Mail",
		Zip:       "Zip",
		Place:     "Place",
		Employees: "Employees",
		Revenue:   "Revenue",
		Form:      "Form",
		Founded:   "Founded",
		Job:       "Job",
		Locations: "Locations",
	}.Slice()
	hl[0] = "Id"
	for _, k := range wzItCodes {
		// for k, val := range wzTitles {
		wzkeys = append(wzkeys, k)
		hl = append(hl, wzTitles[k]+" ("+strconv.Itoa(k)+")")
	}
	resultSlice = append(resultSlice, hl)

	companyStash = parseCompanyDetails(companyStash, cmpIds)
	companyStash = parseCompanyMails(companyStash, cmpIds)

	for _, comp := range companyStash {
		if comp.Mail != "invalid" && comp.Mail != "" {
			resultSlice = append(resultSlice, comp.Slice())
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

func parseCompanyMails(companyStash []Company, cmpIds map[int]int) []Company {
	fileHandler, _ := os.Open("./mails.csv")
	defer fileHandler.Close()

	// Create a new reader.
	csvReader := csv.NewReader(bufio.NewReader(fileHandler))
	csvReader.Comma = ';'
	csvReader.Comment = '#'

	notfound := 0
	for {
		records, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			HandleErr("-> ", err)
		}

		cmpId, _ := strconv.Atoi(records[0])

		if cmpIds[cmpId] > 0 {
			companyStash[cmpIds[cmpId]].Mail = records[2]

		} else {
			notfound += 1
		}
	}
	fmt.Println("mails: ", notfound)
	return companyStash

}

func parseCompanyDetails(companyStash []Company, cmpIds map[int]int) []Company {
	fileHandler, _ := os.Open("./details.csv")
	defer fileHandler.Close()

	// Create a new reader.
	csvReader := csv.NewReader(bufio.NewReader(fileHandler))
	csvReader.Comma = ';'
	csvReader.Comment = '#'

	notfound := 0
	for {
		records, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			HandleErr("-> ", err)
		}

		cmpId, _ := strconv.Atoi(records[0])

		if cmpIds[cmpId] > 0 {
			companyStash[cmpIds[cmpId]].Name = records[1]
			companyStash[cmpIds[cmpId]].Zip = records[3]
			companyStash[cmpIds[cmpId]].Place = records[4]
			companyStash[cmpIds[cmpId]].Employees = records[5]
			companyStash[cmpIds[cmpId]].Revenue = records[6]
		} else {
			notfound += 1
		}
	}
	fmt.Println("Details : ", notfound)
	return companyStash
}
