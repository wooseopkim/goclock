package goclock

import (
	"time"
)

// Goclock is a struct to represent calculated difference in time
// from the given Source Url. It actually does not work as a clock
// by itself, so you'll need to retrieve remote time by adding offset
// and error time value to current time.
type Goclock struct {
	// Source is Url whose Date header was parsed and processed.
	Source string `json:"source"`

	// Offset is the difference in datetime from the client time in original request.
	Offset time.Duration `json:"offset"`

	// ErrorRange is an estimation of error in offset which could have been
	// caused the fact that Date header is note precise more than a second.
	ErrorRange time.Duration `json:"error_range"`
}

// Request data is used for calculating remote time.
type Request struct {
	// Url is the Url from which Date header would be read.
	Url string

	// ClientTime is the time that will be the baseline of time comparison.
	ClientTime time.Time
}

// New returns a Goclock instance initialized by given Request data.
func New(request Request) (*Goclock, error) {
	g := &Goclock{}
	err := g.initialize(request)
	return g, err
}

// Initialize resets the Goclock with given Request.
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
