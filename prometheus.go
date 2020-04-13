package main

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	lastSyncGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "price_bot_website_last_sync",
			Help: "The timestamp on which the website price was last synced",
		},
		[]string{"website"},
	)
	lastPriceObservedGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "price_bot_website_last_price_observed",
			Help: "The last price observed for the website",
		},
		[]string{"website"},
	)
	isAvailableGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "price_bot_website_is_available",
			Help: "Whether or not the product is available in the website",
		},
		[]string{"website"},
	)
)
