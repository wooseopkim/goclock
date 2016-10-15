package goclock

import (
    "time"
)

const trial = 4

type Goclock struct {
    Source string `json:"source"`
    Offset time.Duration `json:"offset"`
    // 0 is the worst; bigger means better
    Reliability int `json:"reliaibility"`
}

type Request struct {
    Url string
    ClientTime time.Time
}

func New(request Request) *Goclock {
    g := &Goclock{}
    g.initialize(request)
    return g
}

func (g *Goclock) Initialize(request Request) {
    g.initialize(request)
}

func (g Goclock) Time() time.Time {
    return time.Now().Add(g.Offset)
}

func (g *Goclock) initialize(request Request) {
    g.Source = request.Url
    offset, reliability := timeOffset(g.Source)
    g.Offset = offset
    g.Reliability = reliability
}