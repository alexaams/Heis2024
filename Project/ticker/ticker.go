package ticker

import (
	"time"
)

var timerActive bool
var timerEndTime float64

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

func getWallTime() float64 {
	return float64(time.Now().UnixNano()) / float64(time.Second)
}

func TimerStart(duration float64) {
	timerEndTime = getWallTime() + duration
	timerActive = true
}

func TimerStop() {
	timerActive = false
}

func TimerTimedOut() bool {
	return ((timerActive) && (getWallTime() > timerEndTime))
}
