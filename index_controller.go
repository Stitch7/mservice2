package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/russross/blackfriday"
)

var readMe string

func Index(w http.ResponseWriter, r *http.Request) {
	if len(readMe) == 0 {
		initReadMe()
	}

	fmt.Fprint(w, readMe)
}

func initReadMe() {
	rawHtml, err := ioutil.ReadFile("index.html")
	if err != nil {
		panic(err)
	}

	rawMarkdown, err2 := ioutil.ReadFile("README.md")
	if err2 != nil {
		panic(err2)
	}

	parsedMarkdown := blackfriday.MarkdownBasic(rawMarkdown)
	html := string(rawHtml[:])
	markdown := string(parsedMarkdown[:])
	readMe = strings.Replace(html, "<!-- README.md -->", markdown, 1)
}
