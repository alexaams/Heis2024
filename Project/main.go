package main

import (
	"ProjectHeis/config"
	"ProjectHeis/drivers/elevio"
	"ProjectHeis/network/localip"
	"ProjectHeis/network/peers"
	"fmt"
	"strconv"
)

func main() {

	//Create and asssign ID
	id := localip.CreateID()

	//ID-channel - updates (new and lost peers)
	peerUpdateCh := make(chan peers.PeerUpdate)
	//Enable-transmit-channel
	peerTxEnable := make(chan bool)
	//Transmit- and receive-threads
	go peers.Transmitter(15647, strconv.Itoa(id), peerTxEnable)
	go peers.Receiver(15647, peerUpdateCh)

	fmt.Println("Started")
	for {
		select {
		case p := <-peerUpdateCh:
			fmt.Printf("Peer update:\n")
			fmt.Printf("  Peers:    %q\n", p.Peers)
			fmt.Printf("  New:      %q\n", p.New)
			fmt.Printf("  Lost:     %q\n", p.Lost)
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
					if !localPeersdata.OrdersHall[floor][a.Button] {
						fmt.Println("Sending message")
					}
				}
			}
		}
	}
}
