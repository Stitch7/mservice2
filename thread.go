// Copyright 2016 Christopher Reitz. All rights reserved.
//
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

type Thread struct {
	Id          int      `json:"id"`
	BoardId     int      `json:"boardId"`
	MessageId   int      `json:"messageId"`
	Sticky      bool     `json:"sticky"`
	Closed      bool     `json:"closed"`
	Username    string   `json:"username"`
	Mod         bool     `json:"mod"`
	Subject     string   `json:"subject"`
	Date        JSONTime `json:"date"`
	AnswerCount int      `json:"answerCount"`
	AnswerDate  JSONTime `json:"answerDate"`
}

type Threads []Thread
