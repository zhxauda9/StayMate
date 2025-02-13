package models

import "time"

type Transaction struct {
	ID            string    `json:"id,omitempty"`
	UserEmail     string    `json:"user_email"`
	Amount        float64   `json:"amount"`
	Products      []string  `json:"products"`
	CreatedAt     time.Time `json:"created_at"`
	PaymentMethod string    `json:"payment_method"`
	CardDetails   string    `json:"card_details"`
}
