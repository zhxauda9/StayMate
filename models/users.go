package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// Модель User для MongoDB
type User struct {
	ID    primitive.ObjectID `json:"id" bson:"_id"`
	Name  string             `json:"name"`
	Email string             `json:"email"`
}
