package main

import (
	"context"
	"fmt"
	"github.com/chromedp/chromedp"
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
	var price float64
	err := chromedp.Run(ctx,
		chromedp.Navigate(cd.productUrl),
		chromedp.Sleep(2*time.Second),
		chromedp.Evaluate(cd.getPriceJS(), &price))
	return price, err
}

func (cd *Cdiscount) getPriceJS() string {
	return fmt.Sprintf(`
function convertPriceStringToFloat(price) {
    return parseFloat(price.replace('€', '.'))
}

function getPrice() {
    const priceElement = document.getElementsByClassName("price")[0]
    return priceElement && priceElement.innerText ? parseFloat(priceElement.innerText) : %d;
}

getPrice();
`, InvalidPrice)
}
