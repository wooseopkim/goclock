package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/wooseopkim/goclock"
)

var url string

func main() {
	start(url)
}

func init() {
	flag.StringVar(&url, "url", "http://example.com/", "url to fetch from")
	flag.Parse()
}

func start(url string) {
	clock, err := goclock.New(goclock.Request{
		URL:        url,
		ClientTime: time.Now(),
	})
	if err != nil {
		panic(err)
	} else {
		tick(*clock)
	}
}

func tick(clock goclock.Goclock) {
	for lastlySeen := timeOf(clock); true; lastlySeen = timeOf(clock) {
		fmt.Println(lastlySeen)
		time.Sleep(time.Second - time.Duration(lastlySeen.Nanosecond()))
	}
}

func timeOf(g goclock.Goclock) time.Time {
	return time.Now().Add(g.Offset - g.ErrorRange)
}
