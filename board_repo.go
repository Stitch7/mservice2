// Copyright 2016 Christopher Reitz. All rights reserved.
//
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
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
