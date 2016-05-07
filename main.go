// Copyright 2016 Christopher Reitz. All rights reserved.
//
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	startServer()
}

func startServer() {
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "3000"
		fmt.Printf("[mservice] listening on port %s", port)
	}

	router := NewRouter()
	log.Fatal(http.ListenAndServe(":"+port, router))
}
