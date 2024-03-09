package main

import (
	"ProjectHeis/config"
	"ProjectHeis/drivers/elevio"
	"ProjectHeis/drivers/fms"
	"ProjectHeis/network/peers"
	"fmt"
)

// We define some custom struct to send over the network.
// Note that all members we want to transmit must be public. Any private members
//
//	will be received as zero-values.
type HelloMsg struct {
	Message string
	Iter    int
}

func main() {
	go peers.SendPeersData_init()
	go fms.InitFms()

	for {
		select {
		case a := <-peers.G_Ch_PeersData_Rx:
			fmt.Printf("Message from ID: %d\n", a.Id)
		}
	}
}

// Må flyttes senere
func UpdatePeersdata(localPeersdata peers.PeersData) {
	ch_hallBtn := make(chan elevio.ButtonEvent)
	go elevio.PollButtons(ch_hallBtn)

	for {
		select {
		case a := <-ch_hallBtn:
			if a.Button == elevio.BT_HallUp || a.Button == elevio.BT_HallDown {
				for floor := 0; floor < config.NumFloors; floor++ {
					if !localPeersdata.SingleOrdersHall[floor][a.Button] {
						fmt.Println("Sending message")
					}
				}
			}
		}
	}
}
