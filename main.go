package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "3000"
		fmt.Printf("[mservice] listening on port %s", port)
	}

	router := NewRouter()
	log.Fatal(http.ListenAndServe(":"+port, router))
}
