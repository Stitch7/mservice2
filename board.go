package main

type Board struct {
	Id          int      `json:"id"`
	Name        string   `json:"name"`
	Topic       string   `json:"topic"`
	LastMessage JSONTime `json:"lastMessage"`
	Mods        []string `json:"mods"`
}

type Boards []Board
