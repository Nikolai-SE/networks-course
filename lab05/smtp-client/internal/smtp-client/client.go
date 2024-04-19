package client

import (
	"bufio"
	"crypto/tls"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
)

////////////////////////////////////////////////////////////////////////////////

type Client struct {
	Address  string
	Password string
	Port     int
}

type Mail struct {
	Sender      string
	To          []string
	Subject     string
	Body        string
	ContentType string
}

func (c *Client) BuildMessage(mail Mail) string {
	msg := fmt.Sprintf("MIME-version: 1.0;\nContent-Type: %s; charset=\"UTF-8\";\r\n", mail.ContentType)
	msg += fmt.Sprintf("From: %s\r\n", mail.Sender)
	msg += fmt.Sprintf("To: %s\r\n", strings.Join(mail.To, ";"))
	msg += fmt.Sprintf("Subject: %s\r\n", mail.Subject)
	msg += fmt.Sprintf("\r\n%s\r\n", mail.Body)
	msg += ".\r\n"

	return msg
}

func (c *Client) SendMail(mail Mail) error {
	readerWriter, conn, err := createReaderWriter(c.Address, c.Port)
	if err != nil {
		return err
	}
	defer conn.Close()

	if err != readerWriter.expect("220") {
		return err
	}

	if err := readerWriter.send("HELO domain.local", "250"); err != nil { // TODO remove %s
		return err
	}

	if err := readerWriter.send("AUTH LOGIN", "334"); err != nil {
		return err
	}

	if err := readerWriter.send(base64.StdEncoding.EncodeToString([]byte(mail.Sender)), "334"); err != nil {
		return err
	}

	if err := readerWriter.send(base64.StdEncoding.EncodeToString([]byte(c.Password)), "235"); err != nil {
		return err
	}

	if err := readerWriter.send(fmt.Sprintf("MAIL FROM: <%s>", mail.Sender), "250"); err != nil {
		return err
	}

	for _, to := range mail.To {
		if err := readerWriter.send(fmt.Sprintf("RCPT TO: <%s>", to), "250"); err != nil {
			return err
		}
	}

	if err := readerWriter.send("DATA", "354"); err != nil {
		return err
	}

	if err := readerWriter.send(c.BuildMessage(mail), "250"); err != nil {
		return err
	}

	readerWriter.send("QUIT", "221")

	return nil
}

////////////////////////////////////////////////////////////////////////////////

type SmtpReaderWriter bufio.ReadWriter

func createReaderWriter(address string, port int) (*SmtpReaderWriter, *tls.Conn, error) {
	cfg := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         address,
	}

	conn, err := tls.Dial("tcp", fmt.Sprintf("%s:%d", address, port), cfg)
	if err != nil {
		return nil, nil, err
	}
	readerWriter := SmtpReaderWriter(bufio.ReadWriter{Reader: bufio.NewReader(conn), Writer: bufio.NewWriter(conn)})
	return &readerWriter, conn, nil
}

func (c *SmtpReaderWriter) expect(expected string) error {
	str, err := c.Reader.ReadString('\n')

	if err != nil {
		return err
	}
	if !strings.HasPrefix(str, expected) {
		return errors.New(str)
	}
	return nil
}

func (c *SmtpReaderWriter) send(data, expected string) error {
	_, err := c.Writer.WriteString(fmt.Sprintf("%s\r\n", data))
	if err != nil {
		return err
	}
	err = c.Writer.Flush()
	if err != nil {
		return err
	}
	return c.expect(expected)
}

////////////////////////////////////////////////////////////////////////////////
