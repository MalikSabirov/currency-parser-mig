package models

import "time"

type Currency struct {
	ID           int       `json:"id,omitempty"`
	CurrencyCode string    `json:"currency_code"`
	BuyRate      float64   `json:"buy_rate"`
	SellRate     float64   `json:"sell_rate"`
	Timestamp    time.Time `json:"timestamp"`
}

type AverageCurrency struct {
	CurrencyCode string  `json:"currency_code"`
	AverageBuy   float64 `json:"average_buy"`
	AverageSell  float64 `json:"average_sell"`
}
