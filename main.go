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
	MinPrice() int

	Name() string

	FetchPrice(ctx context.Context) (int, error)
}

var (
	stmpOptionsFilepath = flag.String("smtp-filepath", "/opt/config/smtp.json", "filepath to JSON SMTP options used to send the price notifications")
	port                = flag.Int("port", 8091, "port on which Prometheus metrics (/prometheus) are exposed")
)

func main() {
	flag.Parse()

	notifier := NewSmtpNotifier(*stmpOptionsFilepath)

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.UserAgent("Pricebot/1.0"),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	taskCtx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	go fetch(taskCtx, 30*time.Second, 0*time.Second, notifier, &MediaMarkt{
		productUrl: "https://www.mediamarkt.nl/nl/product/_nintendo-switch-rood-en-blauw-2019-revisie-1635020.html",
		minPrice:   330,
	})

	go fetch(taskCtx, 30*time.Second, 1*time.Second, notifier, &Bol{
		productUrl: "https://www.bol.com/nl/p/nintendo-switch-console-met-35-eshop-tegoed-voucher-32gb-rood-blauw/9200000114613417/",
		minPrice:   330,
	})

	go fetch(taskCtx, 30*time.Second, 2*time.Second, notifier, &CoolBlue{
		productUrl: "https://www.coolblue.nl/en/product/838252/nintendo-switch-2019-upgrade-red-blue.html",
		minPrice:   330,
	})

	go fetch(taskCtx, 30*time.Second, 3*time.Second, notifier, &Amazon{
		productUrl:     "https://www.amazon.nl/Nintendo-Switch-Console-Rood-Blauw/dp/B07WKNQ8JT",
		minPrice:       330,
		country:        "NL",
		lastPriceGauge: amazonNLLastPrice,
		lastSyncGauge:  amazonNLLastSync,
	})

	go fetch(taskCtx, 30*time.Second, 4*time.Second, notifier, &Amazon{
		productUrl:     "https://www.amazon.fr/Nintendo-Switch-Console-Rood-Blauw/dp/B07WKNQ8JT",
		minPrice:       330,
		country:        "FR",
		lastPriceGauge: amazonFRLastPrice,
		lastSyncGauge:  amazonFRLastSync,
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
		newTab, cancel := chromedp.NewContext(ctx)
		price, err := w.FetchPrice(newTab)
		if err != nil {
			errorLogger.Printf("[%s] %v", w.Name(), err)
		} else {
			if price <= w.MinPrice() {
				notifier.Notify(w.Name(), price)
			} else {
				debugLogger.Printf("[%s] don't buy %d\n", w.Name(), price)
			}
		}
		cancel()
		time.Sleep(frequency)
	}
}
