package main

import (
	"bytes"
	"mime/multipart"
	"net/http"
)

func SendMail(subject, content string) bool {
	var err error
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	writer.WriteField("from", "Mailgun Sandbox <postmaster@sandbox0e70a920c50c477e97fd0bc9e39da583.mailgun.org>")
	writer.WriteField("to", "REHIAR <teamrehiar@gmail.com>")
	writer.WriteField("subject", subject)
	writer.WriteField("text", content)

	err = writer.Close()
	if err != nil {
		Configuration.Logger.Error.Println("couldn't send mail! - ", subject)
		Configuration.Logger.Error.Println("MailContent: ", content)
		return false
	}

	req, err := http.NewRequest("POST", "https://api.mailgun.net/v3/sandbox0e70a920c50c477e97fd0bc9e39da583.mailgun.org/messages", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	if err != nil {
		Configuration.Logger.Error.Println("couldn't send mail! - ", subject)
		Configuration.Logger.Error.Println(content)
		return false
	}

	return true
}
