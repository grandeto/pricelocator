package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

func scrapePrices(url string, tag string, prices chan string, wg *sync.WaitGroup) {
	defer wg.Done()

	// Request the HTML page.
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	prices <- doc.Find(tag).Text() + " - " + url
}

func main() {
	urls := map[string]string{
		"https://sopharmacy.bg/bg/product/000000000030011088": ".price.price--inline.price--l",
		"https://bemore.shop/bg/2-br-collanol-s-25-otst-pka": ".product-info-price .price",
		"https://bemore.shop/bg/collanol-10": ".product-info-price .price",
		"https://subra.bg/bg/hranitelni-dobavki/vitaslim-collanol-x-20-caps.html": "#sec_discounted_price_12778",
		"https://remedium.bg/collanol-intakten-kolagen-i-kurkumin-za-zdravi-kosti-i-stavi-h20-kapsuli-148074/p": ".Price__PriceLabel-sc-14hy5o8-1",
		"https://befit.bg/collanol-kolagen-i-kurkumin-za-zdravi-stavi-i-kosti-20-kaps": ".price-box .price",
		"https://www.aptekadetelina.bg/collanol-kolanol-680-mg-20-kapsuli?manufacturer_id=575" : "#ProductPriceSystem_DAuHUM6x .price",
	}

	var wg sync.WaitGroup

	prices := make(chan string, len(urls))

	for url, tag := range urls {
		wg.Add(1)
		go scrapePrices(url, tag, prices, &wg)
	}

	wg.Wait()
	close(prices)

	for price := range prices {
		fmt.Printf("%s\n\n", price)
	}
}
