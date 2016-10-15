package goclock

import (
    "time"
)

const trial = 4

type Goclock struct {
    Source string
    Offset time.Duration
    // 0 is the worst; bigger means better
    Reliability int
}

type Request struct {
    Url string
    ClientTime time.Time
}

type Callback func (*Goclock)

func New(request Request, callback Callback) *Goclock {
    this := &Goclock{}
    this.initialize(request, callback)
    return this
}

func (this *Goclock) Initialize(request Request, callback Callback) {
    this.initialize(request, callback)
}

func (this Goclock) Time() time.Time {
    return time.Now().Add(this.Offset)
}

func (this *Goclock) initialize(request Request, callback Callback) {
    offset, reliability := timeOffset(request.Url)
    this.Offset = offset
    this.Reliability = reliability
    callback(this)
}