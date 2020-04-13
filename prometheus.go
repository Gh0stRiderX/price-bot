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
	amazonNLLastSync = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "price_bot_amazon_nl_last_sync",
		Help: "The timestamp on which Amazon NL price was last synced",
	})
	amazonNLLastPrice = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "price_bot_amazon_nl_last_price",
		Help: "Last price observed on Amazon NL",
	})
	amazonFRLastSync = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "price_bot_amazon_fr_last_sync",
		Help: "The timestamp on which Amazon FR price was last synced",
	})
	amazonFRLastPrice = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "price_bot_amazon_fr_last_price",
		Help: "Last price observed on Amazon FR",
	})
	amazonDELastSync = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "price_bot_amazon_de_last_sync",
		Help: "The timestamp on which Amazon FR price was last synced",
	})
	amazonDELastPrice = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "price_bot_amazon_de_last_price",
		Help: "Last price observed on Amazon DE",
	})
)
