package ticks

import "time"

func tickerStart(waitDuration time.Duration) <-chan bool{
	startSignal := make(chan bool)
	go func(){
		startTime := time.Now()
		ticker := time.NewTicker(1*time.Second)
		defer ticker.Stop()

		for range ticker.C {
			if time.Since(startTime) >= waitDuration {
				startSignal <- true
				close(startSignal)
				return
			}
		}
	}()
	return startSignal
}