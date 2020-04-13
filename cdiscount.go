package main

import (
	"context"
	"fmt"
	"github.com/chromedp/chromedp"
	"strconv"
	"time"
)

type Cdiscount struct {
	productUrl string
	minPrice   float64
}

func (cd *Cdiscount) Name() string {
	return "CDISCOUNT"
}

func (cd *Cdiscount) MinPrice() float64 {
	return cd.minPrice
}

func (cd *Cdiscount) FetchPrice(ctx context.Context) (float64, error) {
	p, err := cd.getPrice(ctx)
	if err != nil {
		return InvalidPrice, fmt.Errorf("could not fetch price, got error %v", err)
	}
	return p, nil
}

func (cd *Cdiscount) getPrice(ctx context.Context) (float64, error) {
	var price string
	err := chromedp.Run(ctx,
		chromedp.Navigate(cd.productUrl),
		chromedp.Sleep(2*time.Second),
		chromedp.Evaluate(cd.getPriceJS(), &price))
	if err != nil {
		return InvalidPrice, err
	}
	return strconv.ParseFloat(price, 64)
}

func (cd *Cdiscount) getPriceJS() string {
	return `document.getElementsByClassName("price")[0].innerText.replace('€', '.')`
}
