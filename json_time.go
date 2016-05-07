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
	lastMessageString := strings.TrimSpace(str)
	t, err := time.Parse("02.01.06 15:04", lastMessageString)
	if err != nil {
		log.Fatal(err)
	}

	return JSONTime(t)
}
