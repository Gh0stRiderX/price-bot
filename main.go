package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/chromedp/chromedp"
	"io/ioutil"
	"log"
	"net/smtp"
	"os"
	"time"
)

var (
	debugLogger = log.New(os.Stdout, "[DEBUG  ] ", log.LstdFlags)
	infoLogger  = log.New(os.Stdout, "[INFO   ] ", log.LstdFlags)
	errorLogger = log.New(os.Stderr, "[ERROR  ] ", log.LstdFlags)
)

type SmtpCrendentials struct {
	Hostname string `json:"hostname"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type SmtpOptions struct {
	From string `json:"from"`
	To string `json:"to"`
	Auth SmtpCrendentials `json:"auth"`
}

type Website interface {
	MinPrice() int

	Name() string

	FetchPrice(ctx context.Context) (int, error)
}

var (
	stmpOptionsFilepath string
)

func init() {
	flag.StringVar(&stmpOptionsFilepath, "smtp-filepath", "smtp.json", "filepath to JSON SMTP options used to send the price notifications")
}

func main() {
	flag.Parse()

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.UserAgent("Pricebot/1.0"),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	taskCtx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	go fetch(taskCtx, 30*time.Second, 0*time.Second, &MediaMarkt{
		productUrl: "https://www.mediamarkt.nl/nl/product/_nintendo-switch-rood-en-blauw-2019-revisie-1635020.html",
		minPrice:   330,
	})

	select {} // infinite waiting
}

func fetch(ctx context.Context, frequency time.Duration, initialDelay time.Duration, w Website) {
	infoLogger.Printf("[%s] Starting fetching routine in %v", w.Name(), initialDelay)
	time.Sleep(initialDelay)
	for {
		price, err := w.FetchPrice(ctx)
		if err != nil {
			errorLogger.Printf("[%s] %v", w.Name(), err)
		} else {
			if price <= w.MinPrice() {
				notify(w.Name(), price)
			} else {
				debugLogger.Printf("[%s] don't buy %d\n", w.Name(), price)
			}
		}

		time.Sleep(frequency)
	}
}

func notify(name string, price int) {
	infoLogger.Printf("[%s] OMG BUY: %d", name, price)

	options := getSmtpOptions()

	recipients := []string{options.To}
	subject := fmt.Sprintf("[%s] AVAILABLE AT PRICE %d", name, price)
	mailContent := fmt.Sprintf(
		"To: %s\r\n"+"Subject: %s\r\n"+"\r\n"+"The product is available at price %d.\r\n",
		recipients[0],
		subject,
		price)

	auth := smtp.PlainAuth("", options.Auth.Username, options.Auth.Password, options.Auth.Hostname)
	err := smtp.SendMail(options.Auth.Hostname+":25", auth, options.From, recipients, []byte(mailContent))
	if err != nil {
		errorLogger.Printf("[%s] failed to send email to %v: %v", name, recipients, err)
	}
}

func getSmtpOptions() *SmtpOptions {
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
