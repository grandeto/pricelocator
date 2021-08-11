package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

func ExampleScrape(prices chan string, wg *sync.WaitGroup) {
	// Request the HTML page.
	res, err := http.Get("https://sopharmacy.bg/bg/product/000000000030011088")
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

	// Find the review items
	prices <- doc.Find(".price.price--inline.price--l").Text()

}

func main() {
	var wg sync.WaitGroup
	prices := make(chan string)

	go ExampleScrape(prices, &wg)
	res := <-prices
	fmt.Printf("Price: %s\n", res)
	wg.Wait()
}
