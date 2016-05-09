// Copyright 2016 Christopher Reitz. All rights reserved.
//
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
)

var threads Threads

func ThreadRepoReset() {
	threads = threads[:0]
}

func ThreadRepoFindById(id int) Thread {
	thread := Thread{}
	for _, t := range threads {
		if t.Id == id {
			thread = t
			break
		}
	}
	return thread
}

func ThreadRepoFindByBoardId(boardId int) Threads {
	var found Threads
	fmt.Println(len(threads))
	for _, t := range threads {
		if t.BoardId == boardId {
			found = append(found, t)
		}
	}
	return found
}

// TODO: race condtions??
func ThreadRepoCreate(t Thread) Thread {
	threads = append(threads, t)
	return t
}

func ThreadRepoDeleteById(id int) error {
	for i, b := range threads {
		if b.Id == id {
			threads = append(threads[:i], threads[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("Could not find Thread with id of %d to delete", id)
}
