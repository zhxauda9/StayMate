package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/zhxauda9/StayMate/internal/server"
)

func main() {
	fmt.Println("Hotel managment system")

	mux, err := server.InitServer()
	if err != nil {
		os.Exit(1)
	}
	http.ListenAndServe(":8080", mux)
}
