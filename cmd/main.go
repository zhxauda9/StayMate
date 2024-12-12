package main

import (
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("Hotel managment system")

	mux := http.NewServeMux()
	http.ListenAndServe("127.0.0.1:8080", mux)
}
