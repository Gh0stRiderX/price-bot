package main

import (
	"context"
	"fmt"
	"github.com/chromedp/chromedp"
	"github.com/prometheus/client_golang/prometheus"
	"strconv"
	"strings"
)

type Amazon struct {
	productUrl     string
	minPrice       int
	country        string
	lastPriceGauge prometheus.Gauge
	lastSyncGauge  prometheus.Gauge
}

func (a *Amazon) Name() string {
	return "AMAZON_" + a.country
}

func (a *Amazon) MinPrice() int {
	return a.minPrice
}

func (a *Amazon) FetchPrice(ctx context.Context) (int, error) {
	price, err := a.getPrice(ctx)
	if err != nil {
		return -1, fmt.Errorf("could not fetch price, got error %v", err)
	}
	p, err := a.convertPrice(price)
	if err != nil {
		return -1, fmt.Errorf("could not convert price %q to number, got error %v", price, err)
	}
	a.lastPriceGauge.Set(float64(p))
	a.lastSyncGauge.SetToCurrentTime()
	return p, nil
}

func (a *Amazon) convertPrice(price string) (int, error) {
	roundedUp := strings.Split(price, ",")[0]
	p, err := strconv.Atoi(roundedUp)
	if err != nil {
		return -1, err
	}
	return p, nil
}

func (a *Amazon) getPrice(ctx context.Context) (string, error) {
	var price string
	err := chromedp.Run(ctx, a.getPriceActionList(&price)...)
	if err != nil {
		return "", err
	}
	return price, nil
}

func (a *Amazon) getPriceActionList(price *string) []chromedp.Action {
	return []chromedp.Action{
		chromedp.Navigate(a.productUrl),
		chromedp.Evaluate("(document.getElementById('priceblock_ourprice') || {innerHTML: '999'}).innerHTML", price),
	}
}
