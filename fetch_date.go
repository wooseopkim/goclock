package goclock

import (
	"net/http"
	"time"
)

const dateHeaderFmt = "Mon, 02 Jan 2006 15:04:05 GMT"

func fetchDate(url string) (time.Time, error) {
	resp, err := http.Get(url)
	if err != nil {
		return time.Time{}, err
	}

	dateHeader := resp.Header.Get("Date")
	date, err := time.Parse(dateHeaderFmt, dateHeader)
	if err != nil {
		return time.Time{}, err
	}

	return date, nil
}
