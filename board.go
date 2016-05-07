// Copyright 2016 Christopher Reitz. All rights reserved.
//
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

type Board struct {
	Id          int      `json:"id"`
	Name        string   `json:"name"`
	Topic       string   `json:"topic"`
	LastMessage JSONTime `json:"lastMessage"`
	Mods        []string `json:"mods"`
}

type Boards []Board
