package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"time"
)

func storePage(uri string) bool {
	requrl, err := url.Parse(uri)
	if err != nil {
		Configuration.Logger.Warning.Printf("Parsing url failed - %s - %s\n", err, uri)
	}

	now := time.Now()
	y, m, d := now.Date()
	hr := now.Hour()
	mn := now.Minute()
	fileName := fmt.Sprintf("%d%02d%02d%02d%02d-%s-%s", y, m, d, hr, mn, requrl.Host, CalcHash(uri))
	path := fmt.Sprintf("./data/%d%02d%02d%02d%02d/", y, m, d, hr, mn)
	os.MkdirAll(path, 0777)
	file := path + fileName

	Configuration.Logger.Info.Printf("store page (%s)", file)

	if _, err := os.Stat(file); os.IsNotExist(err) {
		var resp *http.Response
		var body []byte

		resp, body = RetrieveUrl(uri)
		if resp != nil {
			if resp.StatusCode != 200 {
				Configuration.Logger.Warning.Printf("StatusCode [%d] Failed to download page - %s\n", resp.StatusCode, uri)
			}
		} else {
			Configuration.Logger.Warning.Printf("BodyNil - Failed - %s\n", uri)
		}
		err = ioutil.WriteFile(file, body, 0644)
		if err != nil {
			Configuration.Logger.Warning.Println("IOerror - Failed to store page - %s, %s", err, uri)
		}
		return false
	} else {
		return true
	}
}
