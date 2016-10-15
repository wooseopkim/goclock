package main

import (
    "fmt"
    "time"
    "github.com/linterpreteur/goclock"
)

func main() {
    clock := goclock.New(goclock.Request{
        Url : url(),
        ClientTime : clientTime(),
    }, func (this *goclock.Goclock) {
        fmt.Println("offset is", this.Offset)
        fmt.Println("border is", time.Second - this.Offset)
        fmt.Println("reliability is", this.Reliability)
    })
    tick(*clock)
}

func url() string {
    return "http://m.naver.com/"
}

func clientTime() time.Time {
    return time.Now()
}

func tick(clock goclock.Goclock) {
    for lastlySeen := clock.Time(); /* true */; lastlySeen = clock.Time() {
	    fmt.Println(lastlySeen)
	    time.Sleep(time.Second - time.Duration(lastlySeen.Nanosecond()))
    }
}