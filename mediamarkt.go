package main

import (
	"context"
	"fmt"
	"github.com/chromedp/chromedp"
	"strconv"
)

type MediaMarkt struct {
	productUrl string
	minPrice   float64
}

func (m *MediaMarkt) Name() string {
	return "MEDIAMARKT"
}

func (m *MediaMarkt) MinPrice() float64 {
	return m.minPrice
}

func (m *MediaMarkt) FetchPrice(ctx context.Context) (float64, error) {
	price, err := m.getPrice(ctx)
	if err != nil {
		return -1, fmt.Errorf("could not fetch price, got error %v", err)
	}
	p, err := m.convertPrice(price)
	if err != nil {
		return -1, fmt.Errorf("could not convert price %q to number, got error %v", price, err)
	}
	mediamarktLastPrice.Set(p)
	mediamarktLastSync.SetToCurrentTime()
	return p, nil
}

func (m *MediaMarkt) convertPrice(price string) (float64, error) {
	p, err := strconv.ParseFloat(price, 64)
	if err != nil {
		return -1, err
	}
	return p, nil
}

func (m *MediaMarkt) getPrice(ctx context.Context) (string, error) {
	var price string
	var ok bool
	err := chromedp.Run(ctx, m.getPriceActionList(&price, &ok)...)
	if err != nil {
		return "", err
	}
	if !ok {
		return "", fmt.Errorf("failed to retrieve attribute value")
	}
	return price, nil
}

func (m *MediaMarkt) getPriceActionList(price *string, ok *bool) []chromedp.Action {
	return []chromedp.Action{
		chromedp.Navigate(m.productUrl),
		chromedp.AttributeValue("[itemprop='price']", "content", price, ok),
	}
}
