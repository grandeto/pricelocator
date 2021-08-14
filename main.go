package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"strconv"
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
		fmt.Println(err)
		return
	}
	fmt.Println("Email Sent Successfully!")
}

func scrapePrices(url string, tag string, prices chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	
	errmsg := ""

	// Request the HTML page.
	res, err := http.Get(url)
	if err != nil {
		log.Println(err)
		errmsg = "http.Get error " + err.Error() + " " + url + "\r\n\r\n"
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Printf("status code error: %d %s", res.StatusCode, res.Status)
		errmsg = res.Status + " " + url + "\r\n\r\n"
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Print(err)
		errmsg = "goquery error " + err.Error() + " " + url + "\r\n\r\n"
	}
	
	if errmsg != "" {
		prices <- errmsg
		return
	}

	prices <- doc.Find(tag).Text() + " - " + url + "\r\n\r\n"
}

func isExecutedToday() bool {
	executionLogPath := "pricelocator-latest-run"
	t := time.Now()
	currentExecutionTime := strconv.Itoa(t.Year()) + "-" + t.Month().String() + "-" + strconv.Itoa(t.Day()) + "\n"
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
	}

	if (isExecutedToday()) {
		fmt.Println("Today execution has been already processed")
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

	msg := "Subject: PriceLocator\r\n\r\n"

	for price := range prices {
		msg += price
	}

	emailNotify([]byte(msg))
}
