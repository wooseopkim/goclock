package goclock

import (
	"errors"
	"fmt"
	"time"
)

type timeRecord struct {
	local  time.Time
	remote time.Time
}

type possibleRange struct {
	offset time.Duration
	length time.Duration
}
type request struct {
	previous      time.Time
	start         time.Time
	delay         time.Duration
	possibleRange possibleRange
}

const minSleep = 50 * time.Millisecond
const maxSleep = time.Second - minSleep

func offset(url string) (time.Duration, time.Duration, error) {
	t, r, err := secondBorder(url, time.Now())
	if err != nil {
		return time.Duration(0), time.Duration(0), err
	}

	adjusted := t.local.Add(r.offset)
	expected := t.remote.Add(time.Second)
	offset := expected.Sub(adjusted)

	return offset, r.length, nil
}

func secondBorder(url string, start time.Time) (timeRecord, possibleRange, error) {
	params := request{
		start: start,
		delay: 0,
		possibleRange: possibleRange{
			offset: 0,
			length: time.Second,
		},
	}
	t, r, err := loop(url, params)
	if err != nil {
		return timeRecord{}, possibleRange{}, err
	}

	return t, r, err
}

func loop(url string, r request) (timeRecord, possibleRange, error) {
	if r.delay > 0 {
		time.Sleep(r.delay)
	}

	d, err := fetchDate(url)
	if err != nil {
		return timeRecord{}, possibleRange{}, err
	}
	if r.delay == 0 {
		r.previous = d
	}

	var changed bool
	switch diff := d.Sub(r.previous); diff {
	case 0 * time.Second:
		changed = false
	case 1 * time.Second:
		changed = true
	default:
		msg := fmt.Sprintf("Date header changed from %v to %v", r.previous, d)
		return timeRecord{}, possibleRange{}, errors.New(msg)
	}

	var newSleep time.Duration
	var newOffset time.Duration
	if changed {
		newOffset = r.possibleRange.offset
		newSleep = time.Second - r.possibleRange.length/2
	} else {
		newSleep = r.possibleRange.length / 2
		newOffset = (r.possibleRange.offset + r.possibleRange.length) % time.Second
	}
	newLength := r.possibleRange.length / 2
	newRange := possibleRange{
		offset: newOffset,
		length: newLength,
	}

	sleepError := time.Now().Sub(r.start) - r.delay
	newSleep = newSleep - sleepError
	if newSleep < minSleep || maxSleep < newSleep {
		return timeRecord{}, newRange, nil
	}

	_, pr, err := loop(url, request{
		previous:      d,
		start:         time.Now(),
		delay:         newSleep,
		possibleRange: newRange,
	})
	if err != nil {
		return timeRecord{}, newRange, nil
	}

	return timeRecord{
		local:  r.start,
		remote: r.previous,
	}, pr, nil
}
