package main

import (
	"flag"
	"fmt"
	"log"
	smtp_client "nikolai/simple-client/internal/smtp-client"
	"strings"
)

var (
	address    = flag.String("address", "", "email server")
	port       = flag.Int("port", 0, "email server port")
	sender     = flag.String("sender", "", "sender email")
	recipients = flag.String("recipient", "", "recipient email")
	text       = flag.String("text", "", "mail text content")
	html       = flag.String("html", "", "mail html content")
	subject    = flag.String("subject", "", "mail subject")
	password   = flag.String("password", "", "password")
)

func main() {
	flag.Parse()

	mail := smtp_client.Mail{
		Sender:  *sender,
		To:      strings.Split(*recipients, ","),
		Subject: *subject,
		Body:    "Empty",
	}

	client := smtp_client.Client{
		Password: *password,
		Address:  *address,
		Port:     *port,
	}

	if len(*html) > 0 {
		mail.Body = *html
		mail.ContentType = "text/html"
	} else {
		mail.Body = *text
		mail.ContentType = "text/plain"
	}

	if err := client.SendMail(mail); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Email sent successfully")
}
