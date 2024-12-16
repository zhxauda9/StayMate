package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// Модель Room для MongoDB
type Room struct {
	ID     primitive.ObjectID `json:"id" bson:"_id"`
	Number int                `json:"number"`
	Class  string             `json:"type"`
	Price  float64            `json:"price"`
}
