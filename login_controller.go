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

func LoginTest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	username, password, ok := r.BasicAuth()
	if !ok {
		w.WriteHeader(http.StatusForbidden)
	} else {
		urlStr := "http://www.maniac-forum.de/pxmboard.php" // TODO put this in some sort of config or so

		form := url.Values{}
		form.Set("mode", "login")
		form.Set("nick", username)
		form.Add("pass", password)
		payload := form.Encode()

		client := &http.Client{}
		req, _ := http.NewRequest("POST", urlStr, bytes.NewBufferString(payload))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Add("Content-Length", strconv.Itoa(len(payload)))
		resp, _ := client.Do(req)

		if resp.StatusCode == 200 {
			body := resp.Body
			defer body.Close()

			doc, err := goquery.NewDocumentFromReader(body)
			if err != nil {
				log.Fatal(err)
			}

			if doc.Find("form").Length() > 0 {
				w.WriteHeader(http.StatusForbidden)
			} else {
				w.WriteHeader(http.StatusOK)
			}
		}
	}
}
