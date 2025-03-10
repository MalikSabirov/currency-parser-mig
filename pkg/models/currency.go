package models

import "time"

type Currency struct {
	ID           int       `json:"id,omitempty"`
	CurrencyCode string    `json:"currency_code"`
	BuyRate      float64   `json:"buy_rate"`
	SellRate     float64   `json:"sell_rate"`
	Timestamp    time.Time `json:"timestamp"`
}