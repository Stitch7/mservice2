// Copyright 2016 Christopher Reitz. All rights reserved.
//
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"log"
	"strings"
	"time"
)

type JSONTime time.Time

func (t JSONTime) MarshalJSON() ([]byte, error) {
	timezoneOffset := time.Now().Format("-07:00")
	stamp := fmt.Sprintf("\"%s\"", time.Time(t).Format("2006-01-02T15:04:00")+timezoneOffset)
	return []byte(stamp), nil
}

func JSONTimeParse(str string) JSONTime {
	str = strings.TrimSpace(str)
	format := "02.01.06 15:04"
	if len(str) != len(format) { // TODO find way to to handle null value
		return JSONTime(time.Unix(0, 0))
	}

	t, err := time.Parse("02.01.06 15:04", str)
	if err != nil {
		log.Fatal(err)
	}

	return JSONTime(t)
}
