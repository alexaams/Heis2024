// Use `go run foo.go` to run your program

package main

import (
	. "fmt"
	"runtime"
	"sync"
	"time"
)

var i = 0
var wg sync.WaitGroup

//func incrementing(ch, quit chan int) {
//	//TODO: increment i 1000000 times
//	for x := 0; x < 1000; x++ {
//		i = <-ch
//		i++
//		ch <- i
//	}
//    quit <-0
//	Println("inc finito")
//}
//
//func decrementing(ch chan int) {
//	//TODO: decrement i 1000000 times
//	for x := 0; x < 1000; x++ {
//		i = <-ch
//		i--
//		ch <- i
//
//	}
//	Println("dec finito")
//}

func number_server(inc, dec, get chan int) {

	for {
		select {
		case <-inc:
			i++
		case <-dec:
			i--
		case <-get:
			return
		}
	}
}

func main() {
	// What does GOMAXPROCS do? What happens if you set it to 1?

	// sets number of max threads
	runtime.GOMAXPROCS(9)
	wg.Add(2)

	// TODO: Spawn both functions as goroutines
	inc := make(chan int)
	dec := make(chan int)
	get := make(chan int)

	go func() {
		for j := 0; j < 100000; j++ {

			inc <- 0
		}
		wg.Done()
	}()

	go func() {
		for j := 0; j < 100000; j++ {

			dec <- 0
		}
		wg.Done()
	}()

	go func() {
		wg.Wait()
		get <- 0
	}()

	number_server(inc, dec, get)

	// We have no direct way to wait for the completion of a goroutine (without additional synchronization of some sort)
	// We will do it properly with channels soon. For now: Sleep.

	time.Sleep(500 * time.Millisecond)
	Println("The magic number is:", i)
}
