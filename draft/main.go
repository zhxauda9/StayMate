package main

import (
	"github.com/gorilla/mux"
	"github.com/zhxauda9/StayMate/draft/internal/db"
	"github.com/zhxauda9/StayMate/draft/internal/handlers"
	"log"
	"net/http"
)

func main() {
	// Создание нового маршрутизатора
	mux := mux.NewRouter()

	// Подключение к базе данных
	db.Connect()
	defer db.Disconnect()

	// Настройка маршрутов
	mux.HandleFunc("/bookings", handlers.GetBookings).Methods("GET")
	mux.HandleFunc("/bookings", handlers.CreateBooking).Methods("POST")
	mux.HandleFunc("/rooms", handlers.GetRooms).Methods("GET")
	mux.HandleFunc("/rooms", handlers.CreateRoom).Methods("POST")
	mux.HandleFunc("/users", handlers.GetUsers).Methods("GET")
	mux.HandleFunc("/users", handlers.CreateUser).Methods("POST")

	// Запуск HTTP сервера
	log.Fatal(http.ListenAndServe(":8080", mux))
}
