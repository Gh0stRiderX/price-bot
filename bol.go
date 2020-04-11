package main

import (
	"context"
	"fmt"
	"github.com/chromedp/chromedp"
	"strconv"
)

type Bol struct {
	productUrl string
	minPrice   int
}

func (b *Bol) Name() string {
	return "BOL"
}

func (b *Bol) MinPrice() int {
	return b.minPrice
}

func (b *Bol) FetchPrice(ctx context.Context) (int, error) {
	price, err := b.getPrice(ctx)
	if err != nil {
		return -1, fmt.Errorf("could not fetch price, got error %v", err)
	}
	p, err := b.convertPrice(price)
	if err != nil {
		return -1, fmt.Errorf("could not convert price %q to number, got error %v", price, err)
	}
	bolLastPrice.Set(float64(p))
	bolLastSync.SetToCurrentTime()
	return p, nil
}

func (b *Bol) convertPrice(price string) (int, error) {
	p, err := strconv.Atoi(price)
	if err != nil {
		return -1, err
	}
	return p, nil
}

func (b *Bol) getPrice(ctx context.Context) (string, error) {
	var price string
	err := chromedp.Run(ctx, b.getBolPriceActionList(&price)...)
	if err != nil {
		return "", err
	}
	return price, nil
}

func (b *Bol) getBolPriceActionList(price *string) []chromedp.Action {
	return []chromedp.Action{
		chromedp.Navigate(b.productUrl),
		chromedp.InnerHTML("[data-test='price'] font font", price),
	}
}
