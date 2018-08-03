package goclock

import (
	"time"
)

func timeOffset(url string) (time.Duration, int, error) {
	const start = 1
	records := make([]comparison, start)
	intervals := [trial]time.Duration{
		time.Duration(float64(time.Second)*1.000) % time.Second,
		time.Duration(float64(time.Second)*0.500) % time.Second,
		time.Duration(float64(time.Second)*0.250) % time.Second,
		time.Duration(float64(time.Second)*0.125) % time.Second,
	}
	margin := int(intervals[len(intervals)-1]) / 2

	reliability := 0
	for i, _ := range intervals {
		timeToSleepFor := intervals[i]
		if i > start && records[i].remoteChanged(records[i-1]) {
			for _, interval := range intervals[:i] {
				timeToSleepFor += interval
			}
		}

		cmp, err := compareDelayed(url, timeToSleepFor)
		if err != nil {
			margin = int(timeToSleepFor)
			return time.Duration(0), 0, err
		}
		reliability++
		records = append(records, cmp)
	}

	nanosecOffset := records[start].client.Nanosecond()
	for i, _ := range records {
		if 0 < i && i < len(records)-1 && !records[i].remoteChanged(records[i+1]) {
			offset := int(time.Second)
			for j := 0; j < i; j++ {
				offset = offset / 2
			}
			nanosecOffset += offset
		}
	}
	_ = margin
	return records[start].estimatedDifference(nanosecOffset /* + margin */), reliability, nil
}
