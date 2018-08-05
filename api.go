package goclock

import (
	"time"
)

type Goclock struct {
	Source     string        `json:"source"`
	Offset     time.Duration `json:"offset"`
	ErrorRange time.Duration `json:"error_range"`
}

type Request struct {
	Url        string
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

func (g *Goclock) initialize(request Request) error {
	g.Source = request.Url
	o, e, err := offset(g.Source)
	if err != nil {
		return err
	}
	g.ErrorRange = e
	g.Offset = o
	return nil
}
