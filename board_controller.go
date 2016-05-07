package main

import (
	"encoding/json"
	"net/http"
)

func BoardIndex(w http.ResponseWriter, r *http.Request) {
	if len(boards) == 0 {
		fetchBoards("http://www.maniac-forum.de/pxmboard.php")
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(boards); err != nil {
		panic(err)
	}
}
