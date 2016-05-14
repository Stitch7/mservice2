// Copyright 2016 Christopher Reitz. All rights reserved.
//
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

var baseUrl = "http://www.maniac-forum.de/pxmboard.php"
var ticker = time.NewTicker(time.Hour / 4)

func init() {
	fetchAll()

	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				fetchAll()
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
	fetchThreads()
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

	BoardRepoReset()
	parseBoardsHtml(body)
}

func fetchThreads() {
	ThreadRepoReset()
	for _, board := range boards {
		boardIdStr := strconv.Itoa(board.Id)
		fmt.Println("fetching threads for boardId " + boardIdStr + "...")
		url := baseUrl + "?mode=threadlist&brdid=" + boardIdStr
		resp, err := http.Get(url)
		if err != nil {
			log.Fatal(err)
		}

		body := resp.Body
		defer body.Close()

		parseThreadsHtml(board.Id, body)
	}
}

func parseBoardsHtml(body io.ReadCloser) {
	threadHtmlDoc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		log.Fatal(err)
	}

	threadHtmlDoc.Find("div > table table:nth-child(3) tr.bg2").Each(func(i int, s *goquery.Selection) {
		idAndTitleA := s.Find("td:nth-child(2) a")
		href, _ := idAndTitleA.Attr("href")

		id, _ := strconv.Atoi(strings.Replace(href, "pxmboard.php?mode=board&brdid=", "", 1))
		name := idAndTitleA.Text()
		topic := s.Find("td:nth-child(3)").Text()
		lastMessage := JSONTimeParse(s.Find("td:nth-child(4)").Text())
		mods := strings.Split(strings.TrimSpace(s.Find("td:nth-child(5)").Text()), "\n")

		BoardRepoCreate(Board{
			Id:          id,
			Name:        name,
			Topic:       topic,
			LastMessage: lastMessage,
			Mods:        mods})

		// fmt.Printf(" * [%s] %s - %s (%s) ### %s\n", id, name, topic, lastMessage, mods)
	})
}

func parseThreadsHtml(boardId int, body io.ReadCloser) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(body)
	bodyStr := buf.String()

	threadEntries := strings.Split(bodyStr, "<br>")
	// Remove first + last (empty) entries
	threadEntries = append(threadEntries[:0], threadEntries[1:]...) // TODO
	lastIndex := len(threadEntries) - 1
	threadEntries = append(threadEntries[:lastIndex], threadEntries[lastIndex+1:]...)

	mainReg, _ := regexp.Compile(`(.+)\s-\s(.+)\sam\s(.+)\(\s.+\s(\d+)\s(?:\|\s[A-Za-z:]+\s(.+)\s|)\)`)

	for _, threadHtml := range threadEntries {
		threadHtmlDoc, err := goquery.NewDocumentFromReader(bytes.NewBufferString(threadHtml))
		if err != nil {
			log.Fatal(err)
		}

		// remove whitespace, tabs and newlines
		threadText := threadHtmlDoc.Text()
		threadText = strings.TrimSpace(threadText)
		threadText = strings.Replace(threadText, "\t", "", -1)
		threadText = strings.Replace(threadText, "\n", "", -1)

		messageHref := threadHtmlDoc.Find("a").First()

		// fishing threadId from ld function call in onclick attribute
		id := 0
		onclick, _ := messageHref.Attr("onclick")
		idReg, _ := regexp.Compile(`ld\((\w.+),0\)`)
		idMatches := idReg.FindStringSubmatch(onclick)
		if len(idMatches) > 1 {
			id = toInt(idMatches[1])
		}

		messageId := 0
		href, _ := messageHref.Attr("href")
		mIdReg, _ := regexp.Compile(`(.+)msgid=(.+)`)
		mIdMatches := mIdReg.FindStringSubmatch(href)
		if len(mIdMatches) > 1 {
			messageId = toInt(mIdMatches[2])
		}

		src, _ := threadHtmlDoc.Find("img").First().Attr("src")
		image := path.Base(src)

		// Sticky threads have pin image
		sticky := image == "fixed.gif"
		// Closed threads have lock image
		closed := image == "closed.gif"

		// // Mods have are marked with the highlight css class
		mod := threadHtmlDoc.Find("span").First().HasClass("highlight")

		mainMatches := mainReg.FindStringSubmatch(threadText)
		if len(mainMatches) > 1 {
			subject := mainMatches[1]
			username := mainMatches[2]
			date := JSONTimeParse(mainMatches[3])
			answerCount := toInt(mainMatches[4])
			answerDate := JSONTimeParse(mainMatches[5])

			ThreadRepoCreate(
				Thread{
					Id:          id,
					BoardId:     boardId,
					MessageId:   messageId,
					Subject:     subject,
					Sticky:      sticky,
					Closed:      closed,
					Username:    username,
					Mod:         mod,
					Date:        date,
					AnswerCount: answerCount,
					AnswerDate:  answerDate})
		}

	}
}

func toInt(str string) int {
	i, err := strconv.Atoi(str)
	if err != nil {
		return 0
	}
	return i
}
