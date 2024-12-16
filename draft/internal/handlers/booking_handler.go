package handlers

import (
	"encoding/json"
	"github.com/zhxauda9/StayMate/draft/internal/db"
	"github.com/zhxauda9/StayMate/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

// Создание новой брони
func CreateBooking(w http.ResponseWriter, r *http.Request) {
	var booking models.Booking
	err := json.NewDecoder(r.Body).Decode(&booking)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	collection := db.GetCollection("bookings")
	booking.ID = primitive.NewObjectID()
	_, err = collection.InsertOne(r.Context(), booking)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(booking)
}

func GetBookings(w http.ResponseWriter, r *http.Request) {
	var bookings []models.Booking
	collection := db.GetCollection("bookings")
	cursor, err := collection.Find(r.Context(), bson.D{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer cursor.Close(r.Context())

	for cursor.Next(r.Context()) {
		var booking models.Booking
		if err := cursor.Decode(&booking); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		bookings = append(bookings, booking)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(bookings)
}
