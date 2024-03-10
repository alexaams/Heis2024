package ticker

import (
	"time"
)

func TickerStart(waitDuration float64) {
	waitDurations := time.Duration(waitDuration) * time.Second
	startTime := time.Now()
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		<-ticker.C
		if time.Since(startTime) >= waitDurations {
			break
		}
	}
}
