package parser

import (
	"currency-parser-mig/pkg/models"
	"currency-parser-mig/pkg/utils"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func ParseCurrencies() {
	log.Println("Start parsing...")

	migUrl := os.Getenv("MIG_BASE_URL")
	if migUrl == "" {
		log.Printf("Error loading MIG_BASE_URL from .env file")
		return
	}

	resp, err := http.Get(migUrl)
	if err != nil {
		log.Printf("Failed to get site: %v", err)
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

	currencies := make([]models.Currency, 0)

	currencyBlock := doc.Find(".informer")

	// date parse
	var timestamp time.Time
	dateBlock := currencyBlock.Find("h4").First()
	dateStr := strings.TrimSpace(dateBlock.Text())

	preposition := "на"

	if !strings.Contains(dateStr, preposition) {
		log.Printf("Failed to parse dateStr, there is no preposition \"%s\": %s", preposition, dateStr)
		return
	}

	dateStr = strings.SplitN(dateStr, preposition, 2)[1] // get date without preposition
	dateStr = strings.TrimSpace(dateStr)
	dateStr, err = utils.ReplaceRussianMonth(dateStr)

	if err != nil {
		log.Printf("Error parsing date: %v", err)
	}

	// get UTC+5
	loc, err := time.LoadLocation("Asia/Almaty")
	if err != nil {
		log.Printf("Error loading timezone:  %v", err)
	}

	if parsedTime, err := time.ParseInLocation(utils.ETALON_DATE, dateStr, loc); err == nil {
		timestamp = parsedTime
	} else {
		log.Printf("Error parsing date: %v", err)
	}

	// currency code, buy, sell rate parse
	currencyBlock.Find("table tr").Each(func(i int, s *goquery.Selection) {
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

		currencies = append(currencies, models.Currency{
			CurrencyCode: strings.TrimSpace(currencyCode),
			BuyRate:      buyRate,
			SellRate:     sellRate,
			Timestamp:    timestamp,
		})
	})

	log.Print(currencies)
}
