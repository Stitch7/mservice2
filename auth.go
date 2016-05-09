// Copyright 2016 Christopher Reitz. All rights reserved.
//
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/PuerkitoBio/goquery"
)

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, _, _ := r.BasicAuth()
		if authorize(r) {
			log.Printf("user %s authorized", username)
			next.ServeHTTP(w, r)
		} else {
			log.Printf("authorization failed for user %s", username)
			w.WriteHeader(http.StatusForbidden)
		}
	})
}

func authorize(r *http.Request) bool {
	username, password, ok := r.BasicAuth()

	if ok {
		form := url.Values{}
		form.Set("mode", "login")
		form.Set("nick", username)
		form.Add("pass", password)
		payload := form.Encode()

		client := &http.Client{}
		req, _ := http.NewRequest("POST", baseUrl, bytes.NewBufferString(payload))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Add("Content-Length", strconv.Itoa(len(payload)))
		resp, _ := client.Do(req)

		if resp.StatusCode == 200 {
			body := resp.Body
			defer body.Close()

			doc, err := goquery.NewDocumentFromReader(body)
			if err != nil {
				log.Fatal(err)
				return false
			}

			if doc.Find("form").Length() == 0 {
				return true
			}
		}
	}

	return false
}
