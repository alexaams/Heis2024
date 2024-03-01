package main

import (
	"ProjectHeis/drivers/config"
	"ProjectHeis/drivers/elevio"
	"fmt"
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
	elevio.Init("localhost:15657", config.NumFloors)

	//Create and asssign ID
	id := config.CreateID()
	fmt.Printf("ID: %d\n", id)

	//Create global order-table
	globalOrderTable := config.CreateGlobalOrderTable()
	//Create channel for button-event
	ch_HallButton_event := make(chan elevio.ButtonEvent)
	//Create channel for button-polling
	button_channel := make(chan elevio.ButtonEvent)
	//Running thread checking ch_HallButton_event
	go UpdateGlobalData(globalOrderTable, ch_HallButton_event)
	go elevio.PollButtons(button_channel)

	for {
		select {
		case a := <-button_channel:
			switch a.Button {
			case elevio.BT_HallUp:
				ch_HallButton_event <- a
			case elevio.BT_HallDown:
				ch_HallButton_event <- a
			default:
				fmt.Println("Nothing happens")
			}
		}
	}

}

// Må flyttes senere
func UpdateGlobalData(GlobalTable config.GlobalOrderTable, ch_HallBtn chan elevio.ButtonEvent) {
	for {
		select {
		case a := <-ch_HallBtn:
			switch a.Button {
			case elevio.BT_HallUp:
				fmt.Printf("Requested hall up from floor %d\n", a.Floor)
				GlobalTable[a.Floor][0].Active = true
				GlobalTable.PrintGlobalOrderTable()
				time.Sleep(10 * time.Millisecond)
			case elevio.BT_HallDown:
				fmt.Printf("Requested hall down from floor %d\n", a.Floor)
				GlobalTable[a.Floor][1].Active = true
				GlobalTable.PrintGlobalOrderTable()
				time.Sleep(10 * time.Millisecond)
			default:
				fmt.Printf("Button type is N/A\n")
				time.Sleep(10 * time.Millisecond)
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

/* OLD CODE
//ID-channel - updates (new and lost peers)
peerUpdateCh := make(chan peers.PeerUpdate)
//Enable-transmit-channel
peerTxEnable := make(chan bool)
//Transmit- and receive-threads
go peers.Transmitter(15647, id, peerTxEnable)
go peers.Receiver(15647, peerUpdateCh)

// We make channels for sending and receiving our custom data types
helloTx := make(chan HelloMsg)
helloRx := make(chan HelloMsg)

//Start transmitting and receiving
go bcast.Transmitter(16569, helloTx)
go bcast.Receiver(16569, helloRx)

// The example message. We just send one of these every second.
go func() {
	helloMsg := HelloMsg{"Hello from " + id, 0}
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
*/
