package main

import (
	"context"
	"fmt"
	"github.com/chromedp/chromedp"
	"time"
)

type Amazon struct {
	productUrl string
	minPrice   float64
	country    string
}

func (a *Amazon) Name() string {
	return "AMAZON_" + a.country
}

func (a *Amazon) MinPrice() float64 {
	return a.minPrice
}

func (a *Amazon) IsAvailable(_ context.Context) (bool, error) {
	// an Amazon product is technically always in stock as it is not present if not
	return true, nil
}

func (a *Amazon) FetchPrice(ctx context.Context) (float64, error) {
	p, err := a.getPrice(ctx)
	if err != nil {
		return InvalidPrice, fmt.Errorf("could not fetch price, got error %v", err)
	}
	return p, nil
}

func (a *Amazon) getPrice(ctx context.Context) (float64, error) {
	var prices float64
	err := chromedp.Run(ctx,
		chromedp.Navigate(a.productUrl),
		chromedp.Sleep(1*time.Second),
		chromedp.Evaluate(a.getLowestPriceFromListViewJS(), &prices))
	return prices, err
}

func (a *Amazon) getLowestPriceFromListViewJS() string {
	return fmt.Sprintf(`
function convertPriceStringToFloat(price) {
    return parseFloat(price.replace('EUR ', '').replace(',', '.'))
}

function getPriceInClass(e, htmlClass) {
    const priceElement = e.getElementsByClassName(htmlClass)[0];
    return priceElement && priceElement.innerText ? convertPriceStringToFloat(priceElement.innerText) : 0;
}

function getLowestPriceFromListView() {
    const priceLines = Array.from(document.getElementsByClassName("olpPriceColumn"));
    return priceLines.reduce((lowestPrice, l) => {
        const price = getPriceInClass(l, "olpOfferPrice");
        const shippingPrice = getPriceInClass(l, "olpShippingPrice");
        if (!price) {
            return lowestPrice
        }

		const totalPrice = price + shippingPrice;
        return totalPrice < lowestPrice ? totalPrice : lowestPrice
    }, %d)
}

getLowestPriceFromListView();
`, InvalidPrice)
}
