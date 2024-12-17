package models

type Room struct {
	id     int     `json:"id"`
	number int     `json:"number"`
	class  string  `json:"type"`
	price  float64 `json:"price"`
}
