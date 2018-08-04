package main

import (
	"fmt"
	"time"

	"github.com/linterpreteur/goclock"
)

func main() {
	initClock(url())
}

func initClock(url string) {
	clock, err := goclock.New(goclock.Request{
		Url:        url,
		ClientTime: clientTime(),
	})
	if err != nil {
		fmt.Println(err)
	} else {
		tick(*clock)
	}
}

func url() string {
	return "http://www.example.com/"
}

func clientTime() time.Time {
	return time.Now()
}

func tick(clock goclock.Goclock) {
	for lastlySeen := clock.Time(); ; /* true */ lastlySeen = clock.Time() {
		fmt.Println(lastlySeen)
		time.Sleep(time.Second - time.Duration(lastlySeen.Nanosecond()))
	}
}
