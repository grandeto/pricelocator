package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func emailNotify(msg []byte) {
	// Sender data.
	from := os.Getenv("PRICELOCATOR_MAIL_FROM")
	password := os.Getenv("PRICELOCATOR_MAIL_PASS")

	// Receiver email address.
	to := []string{
		os.Getenv("PRICELOCATOR_MAIL_TO"),
	}

	// smtp server configuration.
	smtpHost := os.Getenv("PRICELOCATOR_MAIL_HOST")
	smtpPort := "587"
	
	// Authentication.
	auth := smtp.PlainAuth("", from, password, smtpHost)
	
	// Sending email.
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, msg)
	
	if err != nil {
		log.Println(err)
		return
	}

	log.Println("Email Sent Successfully!")
}

func scrapePrices(url string, tag string, prices chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	
	errmsg := ""

	// Request the HTML page.
	res, err := http.Get(url)
	if err != nil {
		log.Println(err)
		errmsg = fmt.Sprintf("http.Get error %s %s\r\n\r\n", err.Error(), url)
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		errmsg = fmt.Sprintf("status code error: %d %s", res.StatusCode, res.Status)
		log.Println(errmsg)
		errmsg += "\r\n\r\n"
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)

	if err != nil {
		log.Print(err)
		errmsg = fmt.Sprintf("goquery error: %s %s\r\n\r\n", err.Error(), url)
	}
	
	if errmsg != "" {
		prices <- errmsg
		return
	}

	prices <- fmt.Sprintf("%s - %s\r\n\r\n", doc.Find(tag).Text(), url)
}

func isExecutedToday() bool {
	executionLogPath := "pricelocator-latest-run"
	t := time.Now()
	currentExecutionTime := fmt.Sprintf("%d-%s-%d\n", t.Year(), t.Month().String(), t.Day())
	logContent := []byte(currentExecutionTime)

	if _, err := os.Stat(executionLogPath); os.IsNotExist(err) {
		ioutil.WriteFile(executionLogPath, logContent, 0644)
		return false
	}

	lastExecutionTime, _ := ioutil.ReadFile(executionLogPath)

	if (string(lastExecutionTime) == currentExecutionTime) {
		return true
	}

	ioutil.WriteFile(executionLogPath, logContent, 0644)

	return false
}

func main() {
	urls := map[string]string{
		"https://sopharmacy.bg/bg/product/000000000030011088": ".price.price--inline.price--l",
		"https://bemore.shop/bg-en/2-br-collanol-s-25-otst-pka": ".product-info-price .price",
		"https://bemore.shop/bg-en/collanol-10": ".product-info-price .price",
		"https://subra.bg/bg/hranitelni-dobavki/vitaslim-collanol-x-20-caps.html": "#sec_discounted_price_12778",
		"https://remedium.bg/collanol-intakten-kolagen-i-kurkumin-za-zdravi-kosti-i-stavi-h20-kapsuli-148074/p": ".Price__PriceLabel-sc-14hy5o8-1",
		"https://befit.bg/collanol-kolagen-i-kurkumin-za-zdravi-stavi-i-kosti-20-kaps": ".price-box .price",
		"https://www.aptekadetelina.bg/collanol-kolanol-680-mg-20-kapsuli?manufacturer_id=575" : "#ProductPriceSystem_DAuHUM6x .price",
		"https://www.afya-pharmacy.bg/productbg/11026" : ".priceRow .productPrice .currPrice",
		"https://bglek.com/product/kolanol-kapsuli-680mg-x20-collanol" : "._product-sidebar ._product-details-price-new",
		"https://aptekamladost.com/product/collanol-kolanol-20-kapsuli/" : ".wd-price-outside .summary-inner .price .woocommerce-Price-amount bdi",
		"https://366.bg/product/kolanol-680mg-kaps" : ".prices-wrapper .prices",
		"https://www.pharmacie.bg/productbg/11772" : ".priceRow .productPrice .currPrice",
		"https://apteka.puls.bg/bg/za-zdravi-kosti-stavi/vitaslim-kolanol-h-20-kaps-67876" : ".product-extra-info .product-prices .price",
		"https://aptekadara.com/vitamini-i-minerali/42798/vitaslim-kolanol-h-20-kaps" : ".taxed-price-value.price-value",
		"https://www.adonis.bg/%D0%B2%D0%B8%D1%82%D0%B0%D1%81%D0%BB%D0%B8%D0%BC-%D0%BA%D0%BE%D0%BB%D0%B0%D0%BD%D0%BE%D0%BB-%D0%BA%D0%B0%D0%BF%D1%81%D1%83%D0%BB%D0%B8-680-%D0%BC%D0%B3-%D1%85-20/14806" : ".visible-sm.visible-xs .product-price",
		"https://www.silabg.com/bg/12655-NOW-UCII-Type-II-Collagen-40-mg--60-Caps-.html" : "#product_price_container .price_1",
	}

	if (isExecutedToday()) {
		log.Println("Today execution has already been processed")
		return
	}

	var wg sync.WaitGroup

	prices := make(chan string, len(urls))

	for url, tag := range urls {
		wg.Add(1)
		go scrapePrices(url, tag, prices, &wg)
	}

	wg.Wait()

	close(prices)

	msg := fmt.Sprintf("Subject: PriceLocator Summary - %d urls\r\n\r\n", len(urls))

	counter := 1

	for price := range prices {
		pricerow := fmt.Sprintf("%d\r\n %s", counter, price)
		counter+=1
		msg += pricerow
	}

	emailNotify([]byte(msg))
}
