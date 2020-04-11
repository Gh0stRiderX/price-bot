package main

import (
	"context"
	"fmt"
	"github.com/chromedp/chromedp"
	"strconv"
)

type CoolBlue struct {
	productUrl string
	minPrice   int
}

func (cb *CoolBlue) Name() string {
	return "COOLBLUE"
}

func (cb *CoolBlue) MinPrice() int {
	return cb.minPrice
}

func (cb *CoolBlue) FetchPrice(ctx context.Context) (int, error) {
	price, err := cb.getPrice(ctx)
	if err != nil {
		return -1, fmt.Errorf("could not fetch price, got error %v", err)
	}
	p, err := cb.convertPrice(price)
	if err != nil {
		return -1, fmt.Errorf("could not convert price %q to number, got error %v", price, err)
	}
	return p, nil
}

func (cb *CoolBlue) convertPrice(price string) (int, error) {
	p, err := strconv.Atoi(price)
	if err != nil {
		return -1, err
	}
	return p, nil
}

func (cb *CoolBlue) getPrice(ctx context.Context) (string, error) {
	var price string
	err := chromedp.Run(ctx, cb.getPriceActionList(&price)...)
	if err != nil {
		return "", err
	}
	return price, nil
}

func (cb *CoolBlue) getPriceActionList(price *string) []chromedp.Action {
	return []chromedp.Action{
		chromedp.Navigate(cb.productUrl),
		chromedp.InnerHTML(".sales-price__current", price),
	}
}
