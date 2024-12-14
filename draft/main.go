package main

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v4"
)

func main() {
	conn, err := pgx.Connect(context.Background(), "postgres://postgres:2005@localhost:1195/StayMate")
	if err != nil {
		log.Fatal("Unable to connect to database:", err)
	}
	defer conn.Close(context.Background())

	fmt.Println("Successfully connected to the database")

	rows, err := conn.Query(context.Background(), "SELECT * FROM hotels")
	if err != nil {
		log.Fatal("Query execution failed:", err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var name string
		var city string
		var price float64
		var roomsAvailable int
		if err := rows.Scan(&id, &name, &city, &price, &roomsAvailable); err != nil {
			log.Fatal("Row scan failed:", err)
		}
		fmt.Printf("Hotel ID: %d, Name: %s, City: %s, Price: %.2f, Rooms Available: %d\n", id, name, city, price, roomsAvailable)
	}

	if err := rows.Err(); err != nil {
		log.Fatal("Row iteration failed:", err)
	}

}
