package main

import (
	"crypto/tls"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"
)

/*
 * - Let the browser agents vary randomly
 */

var cookieJar http.CookieJar

func RetrieveUrl(uri string) (*http.Response, []byte) {
	Configuration.Logger.Info.Printf("http req: %s\n", uri)

	req := buildRequest(uri)
	client := buildClient()

	resp, err := client.Do(req)
	if err != nil {
		Configuration.Logger.Warning.Println(err)
		return nil, nil
	}

	defer resp.Body.Close()
	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		Configuration.Logger.Warning.Println(err)
		return nil, nil
	}

	return resp, contents
}

func buildClient() http.Client {
	if cookieJar == nil {
		cookieJar, _ = cookiejar.New(nil)
	}

	// 10 secs Request Timeout
	timeout := time.Duration(30 * time.Second)

	redirectPolicy := func(req *http.Request, via []*http.Request) error {
		if len(via) >= 10 {
			Configuration.Logger.Warning.Println("stopped after 10 redirects")
			return errors.New("stopped after 10 redirects")
		}
		return nil
	}

	// Deactivate SSL Encryption - unnecessary
	transport := &http.Transport{
		MaxIdleConnsPerHost: 250,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	return http.Client{
		Transport:     transport,
		Timeout:       timeout,
		CheckRedirect: redirectPolicy,
		Jar:           cookieJar,
	}
}

func buildRequest(uri string) *http.Request {
	requrl, err := url.Parse(uri)
	if !requrl.IsAbs() {
		Configuration.Logger.Warning.Printf("relative url: %s", requrl)
		return nil
	}

	if err != nil {
		Configuration.Logger.Warning.Println(err)
		return nil
	}

	req, err := http.NewRequest("GET", requrl.String(), nil)
	if err != nil {
		Configuration.Logger.Warning.Println(err)
		return nil
	}

	// Specify User-Agent
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/45.0.2454.85 Safari/537.36")
	req.Header.Set("Accept-Language", "de-DE,en;q=0.8,de;q=0.6")

	return req
}
