package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// Модель Booking для MongoDB
type Booking struct {
	ID       primitive.ObjectID `json:"id" bson:"_id"`
	UserID   int                `json:"user_id"`
	RoomID   int                `json:"room_id"`
	CheckIn  string             `json:"check_in"`
	CheckOut string             `json:"check_out"`
}
