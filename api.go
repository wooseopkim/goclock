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

func New(request Request) (*Goclock, error) {
    g := &Goclock{}
    err := g.initialize(request)
    return g, err
}

func (g *Goclock) Initialize(request Request) error {
    return g.initialize(request)
}

func (g Goclock) Time() time.Time {
    return time.Now().Add(g.Offset)
}

func (g *Goclock) initialize(request Request) error {
    g.Source = request.Url
    offset, reliability, err := timeOffset(g.Source)
    if err != nil {
        return err
    }
    g.Offset = offset
    g.Reliability = reliability
    return nil
}