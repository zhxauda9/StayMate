package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/zhxauda9/StayMate/draft/internal/db"
	"github.com/zhxauda9/StayMate/draft/internal/handlers"
)

func main() {
	mux := mux.NewRouter()
	db.Connect()
	defer db.Disconnect()
	mux.HandleFunc("/bookings", handlers.GetBookings).Methods("GET")
	mux.HandleFunc("/bookings", handlers.CreateBooking).Methods("POST")
	mux.HandleFunc("/rooms", handlers.GetRooms).Methods("GET")
	mux.HandleFunc("/rooms", handlers.CreateRoom).Methods("POST")
	mux.HandleFunc("/users", handlers.GetUsers).Methods("GET")
	mux.HandleFunc("/users", handlers.CreateUser).Methods("POST")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
