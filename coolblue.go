package main

import (
	"context"
	"fmt"
	"github.com/chromedp/chromedp"
)

type CoolBlue struct {
	productUrl string
	minPrice   float64
}

func (cb *CoolBlue) Name() string {
	return "COOLBLUE"
}

func (cb *CoolBlue) MinPrice() float64 {
	return cb.minPrice
}

func (cb *CoolBlue) FetchPrice(ctx context.Context) (float64, error) {
	price, err := cb.getPrice(ctx)
	if err != nil {
		return INVALID_PRICE, fmt.Errorf("could not fetch price, got error %v", err)
	}
	coolblueLastPrice.Set(price)
	coolblueLastSync.SetToCurrentTime()
	return price, nil
}

func (cb *CoolBlue) getPrice(ctx context.Context) (float64, error) {
	var price float64
	err := chromedp.Run(ctx,
		chromedp.Navigate(cb.productUrl),
		chromedp.Evaluate(cb.getPriceJS(), &price))
	return price, err
}

func (cb *CoolBlue) getPriceJS() string {
	return fmt.Sprintf(`
function convertPriceStringToFloat(price) {
    return parseFloat(price.replace(',', '.').replace('-', 0))
}

function getPrice() {
    const priceElement = document.getElementsByClassName("sales-price__current")[0]
    return priceElement && priceElement.innerText ? parseFloat(priceElement.innerText) : %d;
}

getPrice();
`, INVALID_PRICE)
}
