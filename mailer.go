package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"net/mail"
	"net/smtp"
	"strings"
)

const mailBoundaryMarker = "ACUSTOMANDUNIQUEBOUNDARY"

var sendMails MailSender

type MailSender func(string, smtp.Auth, string, []string, []byte) error

type EmailContext struct {
	From       string
	To         string
	Cc         string
	Subject    string
	Body       string
	Attachment File
	content    string
	Success    bool
}

func SendMail(subject, content string) bool {
	ctx := EmailContext()
	ctx.Subject = subject
	ctx.Body = content
	ctx.From = Configuration.Mail.Sender
	ctx.To = "teamrehiar@gmail.com"
	if !ctx.Send() {
		Configuration.Logger.Error.Println("Mail Configuration Error")
		Configuration.Logger.Error.Println(content)
		return false
	} else {
		Configuration.Logger.Info.Println("Send System online mail")
		return true
	}
}

func (context *EmailContext) encodeAttachment() {
	lineMaxLength := 500
	var attachmentBuf bytes.Buffer
	basedFile := base64.StdEncoding.EncodeToString(context.Attachment.File)

	nbrLines := len(basedFile) / lineMaxLength
	for i := 0; i < nbrLines; i++ {
		attachmentBuf.WriteString(basedFile[i*lineMaxLength:(i+1)*lineMaxLength] + "\n")
	}
	attachmentBuf.WriteString(basedFile[nbrLines*lineMaxLength:])

	context.content += fmt.Sprintf("\r\nContent-Type: %s; name=\"%s\"\r\nContent-Transfer-Encoding:base64\r\nContent-Disposition: attachment; filename=\"%s\"\r\n\r\n%s\r\n--%s--",
		context.Attachment.FileType, context.Attachment.Filename, context.Attachment.Filename, attachmentBuf.String(), mailBoundaryMarker)
}

func (context *EmailContext) encodeHeader() {
	context.content += fmt.Sprintf("From: %s <%s>\r\nTo: %s <%s>\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: multipart/mixed; boundary=%s\r\n--%s",
		context.From, context.From, context.To, context.To, context.Subject, mailBoundaryMarker, mailBoundaryMarker)
}

func (context *EmailContext) encodeBody() {
	encodedBody := base64.StdEncoding.EncodeToString([]byte(context.Body))
	context.content += fmt.Sprintf("\r\nContent-Type: text/plain; charset=\"utf-8\"\r\nContent-Transfer-Encoding: base64\r\n\r\n%s\r\n--%s", encodedBody, mailBoundaryMarker)
}

func (context *EmailContext) prepare() {
	context.encodeHeader()
	context.encodeBody()
	if context.Attachment.Filename != "" {
		context.encodeAttachment()
	}
}

func (context EmailContext) validate() bool {
	_, err := mail.ParseAddressList(context.To)
	if err != nil {
		return false
	}
	return true
}

func (context *EmailContext) Send() bool {
	Configuration.Logger.Info.Println("Sending Email > ", context.Subject)
	auth := smtp.PlainAuth(
		"",
		Configuration.Mail.User,
		Configuration.Mail.Pass,
		Configuration.Mail.Server,
	)

	context.prepare()

	if !context.validate() {
		Configuration.Logger.Warning.Println("Mail Validation failed")
		return false
	}

	var err error
	err = sendMails(
		Configuration.Mail.Server+":"+Configuration.Mail.Port,
		auth,
		context.From,
		strings.Split(context.To, ","),
		[]byte(context.content),
	)
	Configuration.Logger.Info.Println("Email.Size", len(context.content))

	if err != nil {
		Configuration.Logger.Warning.Println("Error: ", err)
		context.Success = false
		return false
	}
	context.Success = true
	Configuration.Logger.Info.Println("Done")
	return true
}
