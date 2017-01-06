package goclock

import (
    "errors"
    "fmt"
    "net/http"
    "time"
)

const immediately = time.Duration(0)
const permissibleMargin = time.Second / 50
var ping = make(map[string]time.Duration)

type comparison struct {
    client time.Time
    remote time.Time
    elapsed time.Duration
}

func compareDelayed(url string, delay time.Duration) (comparison, error) {
    start := time.Now()
    httpClient := http.Client{ Timeout: 5 * time.Second }
    clichan, rmtchan := make(chan time.Time), make(chan time.Time)
    errchan := make(chan error)
    
    if delay != immediately {
        if lastPing, ok := ping[url]; ok {
            delay = delay - lastPing
            if lastPing > delay {
                return comparison{}, errors.New(fmt.Sprintf("Poor connection: %s", url))
            }
        }
        time.Sleep(delay)
    }
    go func(clichan chan time.Time, rmtchan chan time.Time, errchan chan error) {
        clichan <- time.Now()
        req, err := http.NewRequest("GET", url, nil)
        if (err != nil) {
            rmtchan <- time.Time{}
            errchan <- err
        }
        req.Header.Set("User-Agent", "Goclock/1.0")
        
        res, err := httpClient.Do(req)
        if (err != nil) {
            rmtchan <- time.Time{}
            errchan <- err
        } else {
            time, _ := http.ParseTime(res.Header["Date"][0])
            rmtchan <- time
            errchan <- nil
        }
    }(clichan, rmtchan, errchan)
    
    client, remote, err := <-clichan, <-rmtchan, <-errchan
    elapsed := time.Since(start)
    if err != nil {
        return comparison{}, err
    }
    
    ping[url] = elapsed - delay
    if delay != immediately && delay - elapsed > permissibleMargin {
        return comparison{}, errors.New(fmt.Sprintf("Poor connection: %v", elapsed))
    }
    
    return comparison{client, remote, elapsed}, nil
}

func compare(url string) (comparison, error) {
    return compareDelayed(url, immediately)
}

func (c comparison) remoteChanged(that comparison) bool {
    return c.remote.Second() != that.remote.Second()
}

func (c comparison) estimatedDifference(estimatedBorderNanos int) time.Duration {
    client := c.client
    remote := c.remote
    if client.Nanosecond() < estimatedBorderNanos {
        remote = remote.Add(time.Second)
    }
    client = client.Add(time.Duration(-client.Nanosecond() + estimatedBorderNanos) % time.Second)
    return remote.Sub(client)
}