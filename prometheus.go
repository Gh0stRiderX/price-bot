package main

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	lastSync = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "price_bot_website_last_sync",
			Help: "The timestamp on which the website price was last synced",
		},
		[]string{"website"},
	)
	lastPriceObserved = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "price_bot_website_last_price_observed",
			Help: "The last price observed for the website",
		},
		[]string{"website"},
	)
)
