package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/smtp"
)

type Notifier interface {
	Notify(websiteName string, price int)
}

type SmtpCrendentials struct {
	Hostname string `json:"hostname"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type SmtpOptions struct {
	From string           `json:"from"`
	To   string           `json:"to"`
	Auth SmtpCrendentials `json:"auth"`
}

type SmtpNotifier struct {
	options *SmtpOptions
}

func NewSmtpNotifier(stmpOptionsFilepath string) *SmtpNotifier {
	return &SmtpNotifier{
		options: getSmtpOptions(stmpOptionsFilepath),
	}
}

func (n *SmtpNotifier) Notify(websiteName string, price int) {
	infoLogger.Printf("[%s] OMG BUY: %d", websiteName, price)

	recipients := []string{n.options.To}
	subject := fmt.Sprintf("[%s] AVAILABLE AT PRICE %d", websiteName, price)
	mailContent := fmt.Sprintf(
		"To: %s\r\n"+"Subject: %s\r\n"+"\r\n"+"The product is available at price %d.\r\n",
		recipients[0],
		subject,
		price)

	auth := smtp.PlainAuth("", n.options.Auth.Username, n.options.Auth.Password, n.options.Auth.Hostname)
	err := smtp.SendMail(n.options.Auth.Hostname+":25", auth, n.options.From, recipients, []byte(mailContent))
	if err != nil {
		errorLogger.Printf("[%s] failed to send email to %v: %v", websiteName, recipients, err)
	}
}

func getSmtpOptions(stmpOptionsFilepath string) *SmtpOptions {
	var options SmtpOptions
	content, err := ioutil.ReadFile(stmpOptionsFilepath)
	if err != nil {
		errorLogger.Fatalf("failed to retrieve SMTP options from file %q: %v", stmpOptionsFilepath, err)
	}
	err = json.Unmarshal(content, &options)
	if err != nil {
		errorLogger.Fatalf("failed to unmarshal SMTP options from file %q, content %s: %v", stmpOptionsFilepath, content, err)
	}
	return &options
}
