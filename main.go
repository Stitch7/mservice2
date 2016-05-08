// Copyright 2016 Christopher Reitz. All rights reserved.
//
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"log"
	"net/http"
	"os"
	"os/user"
)

func main() {
	startServer()
}

func isRoot() bool {
	usr, _ := user.Current()
	return usr.Uid == "0"
}

func certsExists(certFile string, keyFile string) bool {
	for _, file := range []string{certFile, keyFile} {
		_, err := os.Stat(file)
		if os.IsNotExist(err) {
			return false
		}
	}

	return true
}

func findPort(ssl bool) string {
	port := os.Getenv("PORT")
	if len(port) == 0 {
		user := !isRoot()
		if ssl {
			if user {
				port = "4443"
			} else {
				port = "443"
			}
		} else {
			if user {
				port = "8080"
			} else {
				port = "80"
			}
		}
	}

	return port
}

func startServer() {
	certFile := "certs/server.pem"
	keyFile := "certs/server.key"
	ssl := certsExists(certFile, keyFile)
	port := findPort(ssl)
	router := NewRouter()

	log.Printf("start listening on port %s", port)

	var err error
	if ssl {
		err = http.ListenAndServeTLS(":"+port, certFile, keyFile, router)
	} else {
		err = http.ListenAndServe(":"+port, router)
	}

	log.Fatal(err)
}
