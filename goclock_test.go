package goclock

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"testing"
	"time"
)

const port = 3000

func TestGoclock(t *testing.T) {
	const repeat = 500
	const threshold = 250 * time.Millisecond
	const avgThreshold = 125 * time.Millisecond

	records := []time.Duration{}

	test := func(offset time.Duration) error {
		g, err := New(Request{
			URL:        fmt.Sprintf("http://localhost:%d", port),
			ClientTime: time.Now(),
		})
		if err != nil {
			return err
		}

		diff := g.Offset - g.ErrorRange - offset
		records = append(records, diff)
		if diff < 0 {
			diff = -diff
		}
		if diff > threshold {
			msg := fmt.Sprintf("Error too big: %v", diff)
			return errors.New(msg)
		}

		return nil
	}

	runServerAnd(t, func(offset time.Duration) {
		c := make(chan error)
		for i := 0; i < repeat; i++ {
			go func(i int) {
				time.Sleep(time.Duration(i) * 50 * time.Millisecond)
				c <- test(offset)
			}(i)
		}
		for i := 0; i < repeat; i++ {
			if err := <-c; err != nil {
				t.Error(err)
			}
		}
		close(c)

		sum := 0
		for _, v := range records {
			sum = sum + int(v)
		}
		size := len(records)
		if size > 0 {
			avg := time.Duration(sum / size)
			if avg > avgThreshold || (avg < 0 && -avg > avgThreshold) {
				msg := fmt.Sprintf("Average too high: %v\n", avg)
				t.Error(errors.New(msg))
			}
		}
	})
}

func runServerAnd(t *testing.T, do func(time.Duration)) {
	rand.Seed(int64(time.Now().Nanosecond()))
	offset := time.Duration(rand.Intn(350)) * time.Millisecond
	delay := time.Duration(rand.Intn(50)) * time.Millisecond

	s := &http.Server{
		Addr: fmt.Sprintf(":%d", port),
	}
	http.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		time.Sleep(delay)

		now := time.Now()
		w.Header()["Date"] = []string{
			now.In(now.Location()).Add(offset).Format("Mon, 02 Jan 2006 15:04:05 GMT"),
		}
		w.Write([]byte{})
	})

	go func() {
		if err := s.ListenAndServe(); err != http.ErrServerClosed {
			t.Error(err)
		}
	}()
	defer s.Shutdown(context.Background())

	do(offset)
}
