package ticker

import (
	"fmt"
	"time"
)

func TickerStart(waitDuration float64) {
	waitDurations := time.Duration(waitDuration) * time.Millisecond
	startTime := time.Now()
	ticker := time.NewTicker(1 * time.Millisecond)
	defer ticker.Stop()

	for {
		<-ticker.C
		fmt.Println("waiting")
		if time.Since(startTime) >= waitDurations {
			break
		}
	}
}
