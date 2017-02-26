package main

import (
	"encoding/csv"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const dropGoogleTrends = `drop table googleTrends`

const avoidsDupletsGoogleTrends = `select id from googleTrends where terma=$1 AND termb=$2 AND timestamp=$3`
const insertGoogleTrends = `INSERT INTO googleTrends (terma, termb, timestamp, vala,valb) VALUES ($1,$2,$3,$4,$5) RETURNING id`

const googleTrendsTable = `CREATE TABLE IF NOT EXISTS googleTrends (
    id  SERIAL PRIMARY KEY,
    timestamp   timestamp DEFAULT current_timestamp,
    terma   VARCHAR(255) NOT NULL,
    termb   VARCHAR(255) NOT NULL,
    vala   VARCHAR(255) NOT NULL,
    valb   VARCHAR(255) NOT NULL)`

func ParseGoogleTrends(body string) {

	var duplicates int

	var parseCSVContent = func(body string) string {
		var re = regexp.MustCompile(`(?m)^([-\w\d])+,([\w\d:\(\)\s])+,([\w\d:\(\) ])+`)
		matches := re.FindAllString(body, -1)
		Configuration.Logger.Info.Println("Loaded ", len(matches), " Datasets")
		if len(matches) <= 0 {
			return ""
		}
		return strings.Join(matches, "\n")
	}

	var parseData = func(body string) [][]string {
		if body == "" {
			return nil
		}
		reader := csv.NewReader(strings.NewReader(body))
		reader.FieldsPerRecord = 3
		records, _ := reader.ReadAll()
		return records
	}

	var storeInDatabase = func(a, b string, vala, valb int, tstp time.Time) bool {
		var dupleRowId, appId int

		dupleRow := DbConn().QueryRow(avoidsDupletsGoogleTrends, a, b, tstp.Format(time.RFC3339))
		errAD := dupleRow.Scan(&dupleRowId)
		if errAD != nil {

			err := DbConn().QueryRow(insertGoogleTrends, a, b, tstp.Format(time.RFC3339), vala, valb).Scan(&appId)
			if err != nil {
				Configuration.Logger.Warning.Println(err)
				return false
			}
			return true
		} else {
			duplicates += 1
			return true
		}
	}

	var saveData = func(records [][]string) bool {
		if records != nil {
			duplicates = 0
			terma := strings.TrimSuffix(records[0][1], ": (Deutschland)")
			termb := strings.TrimSuffix(records[0][2], ": (Deutschland)")
			woheaders := records[1:]
			Configuration.Logger.Info.Println(terma, termb)
			for _, rec := range woheaders {
				tstp := fmt.Sprintf("%s-%s-%sT%s:00:00+00:00", rec[0][:4], rec[0][5:7], rec[0][8:10], rec[0][11:13])
				saveTime, _ := time.Parse(time.RFC3339, tstp)
				vala, _ := strconv.Atoi(rec[1])
				valb, _ := strconv.Atoi(rec[2])
				if !storeInDatabase(terma, termb, vala, valb, saveTime) {
					Configuration.Logger.Warning.Println("Failed to save row", rec)
					Configuration.Logger.Warning.Println("Stopped saving data")
					break
					return false
				}

			}
		}
		return true
	}

	func(body string) {
		csvdata := parseCSVContent(body)
		data := parseData(csvdata)
		if saveData(data) {
			Configuration.Logger.Info.Println("saved ", (len(data) - duplicates), " rows")
		}

	}(body)

}
