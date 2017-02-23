package main

import (
	"bufio"
	"fmt"
	"net/url"
	"os"
)

func StorePage(qry, uri string, body []byte) {
	errPrint := func(err error) {
		if err != nil {
			Configuration.Logger.Warning.Println("IOerror  %s - %s", err, uri)
		}
	}
	requrl, err := url.Parse(uri)
	errPrint(err)

	y, m, d := Configuration.StartTime.Date()
	hr := Configuration.StartTime.Hour()
	mn := Configuration.StartTime.Minute()
	fileName := fmt.Sprintf("%d%02d%02d%02d%02d-%s-%s", y, m, d, hr, mn, requrl.Host, CalcHash(uri))
	path := fmt.Sprintf("./data/%d%02d%02d%02d%02d/", y, m, d, hr, mn)
	os.MkdirAll(path, 0777)
	file := path + fileName

	if _, err := os.Stat(file); os.IsNotExist(err) {
		fileHandler, err := os.Create(file)
		errPrint(err)
		defer fileHandler.Close()
		bufioWriter := bufio.NewWriter(fileHandler)
		_, err = bufioWriter.WriteString(fmt.Sprintf("%s;%s;%d-%02d-%02d %02d:%02d;%s\n", qry, uri, y, m, d, hr, mn, requrl.Host))
		errPrint(err)
		_, err = bufioWriter.Write(body)
		errPrint(err)
		bufioWriter.Flush()
		Configuration.Logger.Info.Printf("stored page - %s", file)
	} else {
		Configuration.Logger.Info.Printf("file already exists - %s ", file)
	}
}
