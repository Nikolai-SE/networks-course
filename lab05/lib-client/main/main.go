package main

import (
	"flag"
	"fmt"
	"log"
	"net/smtp"
	"strings"
)

var (
	address   = flag.String("address", "", "email server")
	port      = flag.Int("port", 0, "email server port")
	sender    = flag.String("sender", "", "sender email")
	recipient = flag.String("recipient", "", "recipient email")
	text      = flag.String("text", "", "mail text content")
	html      = flag.String("html", "", "mail html content")
	subject   = flag.String("subject", "", "mail subject")
	user      = flag.String("user", "", "user")
	password  = flag.String("password", "", "password")
)

type Mail struct {
	Sender  string
	To      []string
	Subject string
	Body    string
}

func main() {
	flag.Parse()

	to := []string{
		*recipient,
	}

	mail := Mail{
		Sender:  *sender,
		To:      to,
		Subject: *subject,
		Body:    "Empty",
	}

	addr := fmt.Sprintf("%s:%d", *address, *port)
	host := *address

	var msg string

	if len(*html) > 0 {
		mail.Body = *html
		msg = BuildHtmlMessage(mail)
	} else {
		mail.Body = *text
		msg = BuildTextMessage(mail)
	}

	auth := smtp.PlainAuth("", *user, *password, host)
	err := smtp.SendMail(addr, auth, *sender, to, []byte(msg))

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Email sent successfully")
}

func BuildTextMessage(mail Mail) string {
	msg := "MIME-version: 1.0;\nContent-Type: text/plain; charset=\"UTF-8\";\r\n"
	msg += fmt.Sprintf("From: %s\r\n", mail.Sender)
	msg += fmt.Sprintf("To: %s\r\n", strings.Join(mail.To, ";"))
	msg += fmt.Sprintf("Subject: %s\r\n", mail.Subject)
	msg += fmt.Sprintf("\r\n%s\r\n", mail.Body)

	return msg
}

func BuildHtmlMessage(mail Mail) string {
	msg := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\r\n"
	msg += fmt.Sprintf("From: %s\r\n", mail.Sender)
	msg += fmt.Sprintf("To: %s\r\n", strings.Join(mail.To, ";"))
	msg += fmt.Sprintf("Subject: %s\r\n", mail.Subject)
	msg += fmt.Sprintf("\r\n%s\r\n", mail.Body)

	return msg
}
