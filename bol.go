package main

import (
	"context"
	"fmt"
	"github.com/chromedp/chromedp"
	"time"
)

type Bol struct {
	productUrl string
	minPrice   float64
}

func (b *Bol) Name() string {
	return "BOL"
}

func (b *Bol) MinPrice() float64 {
	return b.minPrice
}

func (b *Bol) FetchPrice(ctx context.Context) (float64, error) {
	p, err := b.getPrice(ctx)
	if err != nil {
		return InvalidPrice, fmt.Errorf("could not fetch price, got error %v", err)
	}
	return p, nil
}

func (b *Bol) getPrice(ctx context.Context) (float64, error) {
	var price float64
	err := chromedp.Run(ctx,
		chromedp.Navigate(b.productUrl),
		chromedp.Sleep(2*time.Second),
		chromedp.Evaluate(b.getPriceJS(), &price))
	return price, err
}

func (b *Bol) getPriceJS() string {
	return fmt.Sprintf(`
function convertPriceStringToFloat(price) {
    return parseFloat(price.replace('\n', '.').replace('-', 0))
}

function getPrice() {
    const priceElement = document.getElementsByClassName("promo-price")[0];
    return priceElement && priceElement.innerText ? convertPriceStringToFloat(priceElement.innerText) : %d;
}

getPrice();
`, InvalidPrice)
}
