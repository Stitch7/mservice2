package main

import (
	"fmt"
	"golang.org/x/net/html"
	"io"
	"net/http"
	"time"
)

var currentBoardId int

var boards Boards

// Give us some seed data
func init() {
	BoardRepoCreate(Board{Id: 1, Name: "Smalltalk", Topic: "Diskussionen rund um die Welt der Videospiele.", LastMessage: parseTime("2016-04-30T14:27:00+02:00"), Mods: []string{"Andi", "Rocco", "Leviathan", "Slapshot"}})
	BoardRepoCreate(Board{Id: 2, Name: "For Sale", Topic: "Private Kleinanzeigen: An- und Verkauf gebrauchter Spiele", LastMessage: parseTime("2016-04-30T14:26:00+02:00"), Mods: []string{"Andi", "Rocco", "Leviathan", "pzykoskinhead", "Slapshot"}})
	BoardRepoCreate(Board{Id: 4, Name: "Retro'n'Tech", Topic: "Retro-Themen, Umbau-Lösungen, Anschluss-Probleme, Computerprobleme, Spielehilfen", LastMessage: parseTime("2016-04-30T09:14:00+02:00"), Mods: []string{"Slapshot", "Rocco", "Leviathan", "Andi"}})
	BoardRepoCreate(Board{Id: 6, Name: "OT", Topic: "Ohne Tiefgang - der tägliche Schwachsinn", LastMessage: parseTime("2016-04-30T14:30:00+02:00"), Mods: []string{"Andi", "Leviathan", "Rocco", "Slapshot", "Florian M."}})
	BoardRepoCreate(Board{Id: 26, Name: "Filme & Serien", Topic: "Alles wofür 24 fps reichen", LastMessage: parseTime("2016-04-30T13:41:00+02:00"), Mods: []string{"Andi", "Rocco", "Leviathan", "Slapshot"}})
	BoardRepoCreate(Board{Id: 8, Name: "Online-Gaming", Topic: "Alles rund um Onlinespiele", LastMessage: parseTime("2016-04-22T10:23:00+02:00"), Mods: []string{"Mod-Team"}})
}

func BoardRepoFindById(id int) Board {
	for _, b := range boards {
		if b.Id == id {
			return b
		}
	}
	// return empty Board if not found
	return Board{}
}

//this is bad, I don't think it passes race condtions
func BoardRepoCreate(b Board) Board {
	currentBoardId += 1
	// b.Id = currentBoardId
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

func fetchBoards(url string) Boards {
	resp, err := http.Get(url)
	if err != nil {
		// handle error
	}
	body := resp.Body
	defer body.Close() // close Body when the function returns

	return parseHtml(body)
}

// Helper function to pull the href attribute from a Token
func isBoardTable(t html.Token) bool {
	// Iterate over all of the Token's attributes until we find an "href"
	correctWith := false
	correctCellspacing := false
	correctCellpadding := false

	for _, a := range t.Attr {
		if a.Key == "width" && a.Val == "775" {
			correctWith = true
		}
		if a.Key == "cellspacing" && a.Val == "2" {
			correctCellspacing = true
		}
		if a.Key == "cellpadding" && a.Val == "5" {
			correctCellpadding = true
		}
	}

	return correctWith && correctCellspacing && correctCellpadding
}

func parseHtml(body io.ReadCloser) Boards {
	z := html.NewTokenizer(body)

	var result Boards

	for {
		tt := z.Next()

		switch {
		case tt == html.ErrorToken:
			// End of the document, we're done
			return result
		case tt == html.StartTagToken:
			t := z.Token()

			// Check if the token is an <table> tag
			if t.Data != "table" {
				continue
			}

			// Extract the href value, if there is one
			if !isBoardTable(t) {
				continue
			}

			b := Board{Id: 1, Name: "Smalltalk", Topic: "Diskussionen rund um die Welt der Videospiele.", LastMessage: parseTime("2016-04-30T14:27:00+02:00"), Mods: []string{"Andi", "Rocco", "Leviathan", "Slapshot"}}

			result = append(result, b)

		}
	}

	return result
}

func parseTime(timeString string) time.Time {
	t, e := time.Parse(time.RFC3339, timeString)
	if e != nil {
		fmt.Println(e)
	}

	return t
}
