package main

import (
	"ProjectHeis/drivers/elevio"
	"ProjectHeis/drivers/fsm"
	"ProjectHeis/network/peers"
	"fmt"
)

func main() {


	go fsm.Fsm(elevator.G_Ch_requests)

	go eventHandler.EventHandling()

	for {
		select {}
	}
}
