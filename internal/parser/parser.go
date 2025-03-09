package parser

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/PuerkitoBio/goquery"
)

func ParseCurrencies() {
	log.Println("Start parsing...")

	migUrl := os.Getenv("MIG_BASE_URL")
	if migUrl=="" {
		log.Printf("Error loading MIG_BASE_URL from .env file")
		return
	}

	resp, err := http.Get(migUrl)
	if err != nil {
		log.Printf("Failed to fetch webpage: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Bad response status: %d", resp.StatusCode)
		return
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Printf("Failed to parse HTML: %v", err)
		return
	}

	defer resp.Body.Close()

	doc.Find(".informer table tr").Each(func(i int, s *goquery.Selection) {
		currencyCode := s.Find("td.currency").Text()
		if currencyCode == "" {
			return
		}

		buyStr := s.Find("td.buy").Text()
		sellStr := s.Find("td.sell").Text()

		buyRate, err := strconv.ParseFloat(buyStr, 64)
		if err != nil {
			log.Printf("Error parsing buy rate for %s: %v", currencyCode, err)
			return
		}

		sellRate, err := strconv.ParseFloat(sellStr, 64)
		if err != nil {
			log.Printf("Error parsing sell rate for %s: %v", currencyCode, err)
			return
		}

		log.Println("buyRate value:", buyRate)
		log.Println("sellRate value:", sellRate)
	})
}