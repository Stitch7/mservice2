// Copyright 2016 Christopher Reitz. All rights reserved.
//
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func ThreadIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	vars := mux.Vars(r)
	boardId := vars["boardId"]
	boardThreads := ThreadRepoFindByBoardId(toInt(boardId))

	if err := json.NewEncoder(w).Encode(boardThreads); err != nil {
		panic(err)
	}
}
