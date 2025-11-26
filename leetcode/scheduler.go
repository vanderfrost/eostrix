package leetcode

import (
	"log"
	"time"
)

// schedule task to be run at midnight UTC with a 5 minute delay
// restart safe as it calculates the duration until next post on startup
func ScheduleMidnightUTCEvent(event func()) {
	go func() {
		for {
			now := time.Now().UTC()
			nextRun := time.Date(
				now.Year(), now.Month(), now.Day()+1,
				0, 5, 0, 0, time.UTC,
			)

			durationUntilNext := time.Until(nextRun)
			log.Printf("event scheduled for %v (in %v)\n", nextRun, durationUntilNext)

			time.Sleep(durationUntilNext)

			func() {
				defer func() {
					if r := recover(); r != nil {
						log.Printf("panic in event: %v\n", r)
					}
				}()
				event()
			}()
		}
	}()
}
