package main

import (
	"context"
	"flag"
	"github.com/chromedp/chromedp"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/cors"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

var (
	debugLogger = log.New(os.Stdout, "[DEBUG  ] ", log.LstdFlags)
	infoLogger  = log.New(os.Stdout, "[INFO   ] ", log.LstdFlags)
	errorLogger = log.New(os.Stderr, "[ERROR  ] ", log.LstdFlags)
)

type Website interface {
	MinPrice() float64

	Name() string

	FetchPrice(ctx context.Context) (float64, error)
}

const (
	ExpectedPrice = 330
	InvalidPrice  = 9999999
)

var (
	stmpOptionsFilepath = flag.String("smtp-filepath", "/opt/config/smtp.json", "filepath to JSON SMTP options used to send the price notifications")
	port                = flag.Int("port", 8091, "port on which Prometheus metrics (/prometheus) are exposed")
)

func main() {
	flag.Parse()

	notifier := NewSmtpNotifier(*stmpOptionsFilepath)

	taskCtx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	go fetch(taskCtx, 30*time.Second, 0*time.Second, notifier, &MediaMarkt{
		productUrl: "https://www.mediamarkt.nl/nl/product/_nintendo-switch-rood-en-blauw-2019-revisie-1635020.html",
		minPrice:   ExpectedPrice,
	})

	go fetch(taskCtx, 30*time.Second, 1*time.Second, notifier, &Bol{
		productUrl: "https://www.bol.com/nl/p/nintendo-switch-console-met-35-eshop-tegoed-voucher-32gb-rood-blauw/9200000114613417/",
		minPrice:   ExpectedPrice,
	})

	go fetch(taskCtx, 30*time.Second, 2*time.Second, notifier, &CoolBlue{
		productUrl: "https://www.coolblue.nl/en/product/838252/nintendo-switch-2019-upgrade-red-blue.html",
		minPrice:   ExpectedPrice,
	})

	go fetch(taskCtx, 30*time.Second, 3*time.Second, notifier, &Amazon{
		productUrl:     "http://amazon.nl/gp/offer-listing/B07WKNQ8JT",
		minPrice:       ExpectedPrice,
		country:        "NL",
		lastPriceGauge: amazonNLLastPrice,
		lastSyncGauge:  amazonNLLastSync,
	})

	go fetch(taskCtx, 30*time.Second, 4*time.Second, notifier, &Amazon{
		productUrl:     "https://www.amazon.fr/gp/offer-listing/B07WKNQ8JT",
		minPrice:       ExpectedPrice,
		country:        "FR",
		lastPriceGauge: amazonFRLastPrice,
		lastSyncGauge:  amazonFRLastSync,
	})

	go fetch(taskCtx, 30*time.Second, 5*time.Second, notifier, &Amazon{
		productUrl:     "https://www.amazon.de/gp/offer-listing/B07WKNQ8JT",
		minPrice:       ExpectedPrice,
		country:        "DE",
		lastPriceGauge: amazonDELastPrice,
		lastSyncGauge:  amazonDELastSync,
	})

	mux := http.NewServeMux()
	mux.Handle("/prometheus", promhttp.Handler())

	handler := cors.Default().Handler(mux)
	errorLogger.Fatal(http.ListenAndServe(":"+strconv.Itoa(*port), handler))
}

func fetch(ctx context.Context, frequency time.Duration, initialDelay time.Duration, notifier Notifier, w Website) {
	infoLogger.Printf("[%s] Starting fetching routine in %v", w.Name(), initialDelay)
	time.Sleep(initialDelay)
	for {
		infoLogger.Printf("[%s] Checking price...", w.Name())
		newTab, closeTab := chromedp.NewContext(ctx)
		price, err := w.FetchPrice(newTab)
		if err != nil {
			errorLogger.Printf("[%s] %v", w.Name(), err)
		} else {
			if price <= w.MinPrice() {
				notifier.Notify(w.Name(), price)
			} else {
				debugLogger.Printf("[%s] don't buy %.2f\n", w.Name(), price)
			}
		}
		closeTab()
		time.Sleep(frequency)
	}
}
