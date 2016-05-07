package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

var baseUrl = "http://www.maniac-forum.de/pxmboard.php"
var ticker = time.NewTicker(10 * time.Second)

func init() {
	fetchAll()

	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				fetchBoards()
			case <-quit:
				fmt.Println("TICKER STOPPED")
				ticker.Stop()
				return
			}
		}
	}()
}

func fetchAll() {
	fetchBoards()
}

func fetchBoards() {
	fmt.Println("fetching boards ...")
	url := baseUrl
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	body := resp.Body
	defer body.Close()

	parseBoardsHtml(body)
}

func parseBoardsHtml(body io.ReadCloser) {
	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		log.Fatal(err)
	}

	BoardReset()
	doc.Find("div > table table:nth-child(3) tr.bg2").Each(func(i int, s *goquery.Selection) {
		idAndTitleA := s.Find("td:nth-child(2) a")
		href, _ := idAndTitleA.Attr("href")

		id, _ := strconv.Atoi(strings.Replace(href, "pxmboard.php?mode=board&brdid=", "", 1))
		name := idAndTitleA.Text()
		topic := s.Find("td:nth-child(3)").Text()
		lastMessage := JSONTimeParse(s.Find("td:nth-child(4)").Text())
		mods := strings.Split(strings.TrimSpace(s.Find("td:nth-child(5)").Text()), "\n")

		BoardRepoCreate(Board{Id: id, Name: name, Topic: topic, LastMessage: lastMessage, Mods: mods})
		// fmt.Printf(" * [%s] %s - %s (%s) ### %s\n", id, name, topic, lastMessage, mods)
	})
}
