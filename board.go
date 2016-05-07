package main

import "time"

type Board struct {
	Id          int       `json:"id"`
	Name        string    `json:"name"`
	Topic       string    `json:"topic"`
	LastMessage time.Time `json:"lastMessage"`
	Mods        []string  `json:"mods"`
}

type Boards []Board
