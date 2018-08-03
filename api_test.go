package goclock

import (
	"math/rand"
	"net/http"
	"testing"
	"time"
)

func TestGoclock(t *testing.T) {
	runServerAnd(t, func(offset time.Duration) {
		g, err := New(Request{
			Url:        "http://localhost:3000",
			ClientTime: time.Now(),
		})
		if err != nil {
			t.Error(err)
		}

		diff := g.Offset - offset
		if diff < 0 {
			diff = -diff
		}
		if diff > 200*time.Millisecond {
			t.Error("Error too big", g.Reliability, g.Offset, offset)
		}
	})
}

func runServerAnd(t *testing.T, do func(time.Duration)) {
	offset := time.Duration(rand.Intn(350)) * time.Millisecond
	delay := time.Duration(rand.Intn(10)) * time.Millisecond
	s := &http.Server{
		Addr: ":3000",
	}
	http.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		time.Sleep(delay)
		headers := w.Header()
		headers["Date"] = []string{
			time.Now().In(time.UTC).Add(offset).Format("Mon, 02 Jan 2006 15:04:05 GMT"),
		}
		w.Write([]byte{})
	})
	defer s.Shutdown(nil)

	go func() {
		if err := s.ListenAndServe(); err != nil {
			t.Error(err)
		}
	}()

	do(offset + delay)
}
