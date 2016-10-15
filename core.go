package goclock

import (
    "time"
)

func timeOffset(url string) (time.Duration, int) {
    const start = 1
    records := make([]comparison, start)
    intervals := [trial]time.Duration{
        time.Duration(float64(time.Second) * 1.000) % time.Second,
        time.Duration(float64(time.Second) * 0.500) % time.Second,
        time.Duration(float64(time.Second) * 0.250) % time.Second,
        time.Duration(float64(time.Second) * 0.125) % time.Second,
    }
    margin := int(intervals[len(intervals) - 1]) / 2
    
    reliability := 0
    for i, _ := range intervals {
        timeToSleepFor := intervals[i]
        if i > start && records[i].remoteChanged(records[i - 1]) {
            for _, interval := range intervals[:i] {
                timeToSleepFor += interval
            }
        }
        // fmt.Println("Sleep", timeToSleepFor)
        
        cmp, err := compareDelayed(url, timeToSleepFor)
        if err != nil {
            margin = int(timeToSleepFor)
            break
        }
        reliability++;
        records = append(records, cmp)
        
        /*
        fmt.Printf("%02d:%02d:%02d:%09d",
            cmp.client.Hour(), cmp.client.Minute(),
            cmp.client.Second(), cmp.client.Nanosecond())
        fmt.Printf(" => ")
        fmt.Printf("%02d:%02d:%02d:%09d\n",
                cmp.remote.Hour(), cmp.remote.Minute(),
                cmp.remote.Second(), cmp.remote.Nanosecond())
        */
    }
    
    if reliability == 0 {
        return time.Duration(0), 0
    }
    
    nanosecOffset := records[start].client.Nanosecond()
    // fmt.Println("(nanosecOffset + margin) equals", time.Duration(nanosecOffset + margin))
    for i, _ := range records {
	    if 0 < i && i < len(records) - 1 && !records[i].remoteChanged(records[i + 1]) {
	        offset := int(time.Second)
	        for j := 0; j < i; j++ {
	            offset = offset / 2
	        }
	        nanosecOffset += offset
	        // fmt.Println("+=", time.Duration(offset))
            // fmt.Println("(nanosecOffset + margin) equals", time.Duration(nanosecOffset + margin))
	    }
    }
    _ = margin
    return records[start].estimatedDifference(nanosecOffset /* + margin */), reliability
}