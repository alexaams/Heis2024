// Use `go run foo.go` to run your program

package main

import (
	. "fmt"
	"runtime"
	"time"
)

var i = 0

func incrementing(ch chan int) {
	//TODO: increment i 1000000 times
	for x := 0; x < 1000000; x++ {
		i++
		ch <- i
	}

}

func decrementing(ch chan int) {
	//TODO: decrement i 1000000 times
	for x := 0; x < 1000000; x++ {
		i--
		ch <- i
	}
}

func main() {
	// What does GOMAXPROCS do? What happens if you set it to 1?

	// sets number of max threads
	runtime.GOMAXPROCS(2)

	// TODO: Spawn both functions as goroutines
	ch := make(chan int)

	go incrementing(ch)
	go decrementing(ch)

	// We have no direct way to wait for the completion of a goroutine (without additional synchronization of some sort)
	// We will do it properly with channels soon. For now: Sleep.
	time.Sleep(500 * time.Millisecond)
	Println("The magic number is:", i)
}
