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
    g := &Goclock{}
    g.initialize(request, callback)
    return g
}

func (g *Goclock) Initialize(request Request, callback Callback) {
    g.initialize(request, callback)
}

func (g Goclock) Time() time.Time {
    return time.Now().Add(g.Offset)
}

func (g *Goclock) initialize(request Request, callback Callback) {
    offset, reliability := timeOffset(request.Url)
    g.Offset = offset
    g.Reliability = reliability
    callback(g)
}