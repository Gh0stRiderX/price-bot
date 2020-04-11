package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	mediamarktLastSync = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "price_bot_mediamarkt_last_sync",
		Help: "The timestamp on which Mediamarkt price was last synced",
	})
	mediamarktLastPrice = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "price_bot_mediamarkt_last_price",
		Help: "Last price observed on Mediamarkt",
	})
	bolLastSync = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "price_bot_bol_last_sync",
		Help: "The timestamp on which Bol price was last synced",
	})
	bolLastPrice = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "price_bot_bol_last_price",
		Help: "Last price observed on Bol",
	})
	coolblueLastSync = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "price_bot_coolblue_last_sync",
		Help: "The timestamp on which CoolBlue price was last synced",
	})
	coolblueLastPrice = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "price_bot_coolblue_last_price",
		Help: "Last price observed on CoolBlue",
	})
)
