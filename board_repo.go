package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

var boards Boards

func BoardReset() {

	boards = boards[:0]
}

func BoardRepoFindById(id int) Board {
	board := Board{}
	for _, b := range boards {
		if b.Id == id {
			board = b
			break
		}
	}

	return board
}

// TODO: race condtions??
func BoardRepoCreate(b Board) Board {
	boards = append(boards, b)
	return b
}

func BoardRepoDeleteById(id int) error {
	for i, b := range boards {
		if b.Id == id {
			boards = append(boards[:i], boards[i+1:]...)
			return nil
		}
	}

	return fmt.Errorf("Could not find Board with id of %d to delete", id)
}

func fetchBoards(url string) {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	body := resp.Body
	defer body.Close()

	parseHtml(body)
}

func parseHtml(body io.ReadCloser) {
	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		log.Fatal(err)
	}

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
