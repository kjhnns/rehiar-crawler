package main

import (
	"crypto/tls"
	"errors"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"
)

/*
	TODOS
 * - Let the browser agents vary randomly
*/

func Download(url string) []byte {
	var resp *http.Response
	var body []byte

	sleeper()

	resp, body = httpRequest(url)
	if resp != nil {
		if resp.StatusCode != 200 {
			Configuration.Logger.Warning.Printf("[%d] StatusCode - %s\n", resp.StatusCode, url)
		}
	} else {
		Configuration.Logger.Warning.Printf("BodyNil - %s\n", url)
	}
	return body
}

func sleeper() {
	rand.Seed(time.Now().UnixNano())
	randomSleepTime := calcRandWithVariance(Configuration.SleepTime, 10)
	Configuration.Logger.Info.Println("RandomSleeping: ", randomSleepTime)
	time.Sleep(time.Duration(randomSleepTime) * time.Second)
}

func calcRandWithVariance(base, variance int) int {
	return base + rand.Intn(variance)*(rand.Intn(3)-1)
}

var cookieJar http.CookieJar

func httpRequest(uri string) (*http.Response, []byte) {
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
