package ticks

import "time"

func tickerStart(waitDuration time.Duration) {
	startTime := time.Now()
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		<-ticker.C
		if time.Since(startTime) >= waitDuration {
			break
		}
	}
}
