package main

import (
	"ProjectHeis/config"
	"ProjectHeis/drivers/elevio"
	"ProjectHeis/network/bcast"
	"ProjectHeis/network/localip"
	"ProjectHeis/network/peers"
	"fmt"
	"strconv"
	"time"
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

	//Create and asssign ID
	id := localip.CreateID()

	//ID-channel - updates (new and lost peers)
	peerUpdateCh := make(chan peers.PeerUpdate)
	//Enable-transmit-channel
	peerTxEnable := make(chan bool)
	//Transmit- and receive-threads
	go peers.Transmitter(15647, strconv.Itoa(id), peerTxEnable)
	go peers.Receiver(15647, peerUpdateCh)

	// We make channels for sending and receiving our custom data types
	helloTx := make(chan HelloMsg)
	helloRx := make(chan HelloMsg)

	//Start transmitting and receiving
	go bcast.Transmitter(16569, helloTx)
	go bcast.Receiver(16569, helloRx)

	// The example message. We just send one of these every second.
	go func() {
		helloMsg := HelloMsg{"Hello from ", id}
		for {
			helloMsg.Iter++
			helloTx <- helloMsg
			time.Sleep(1 * time.Second)
		}
	}()

	fmt.Println("Started")
	for {
		select {
		case p := <-peerUpdateCh:
			fmt.Printf("Peer update:\n")
			fmt.Printf("  Peers:    %q\n", p.Peers)
			fmt.Printf("  New:      %q\n", p.New)
			fmt.Printf("  Lost:     %q\n", p.Lost)

		case a := <-helloRx:
			fmt.Printf("Received: %#v\n", a)
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

/*
Neste steg (Jørgen):
1. Primary og backup

2. Påse at heartbeat/watchdog kjører i bakgrunnen som tiltenkt

3. Alle heiser skal ha en lokal PeersData. Vi lager en thread som leser sine knapper, sjekker
opp mot Peersdata, og broadcaster eventuelt ut en melding om at nå er det ny ordre ved mismatch mot
Peersdata. Alle skal kvittere på denne meldingen.
*/
