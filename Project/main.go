package main

import (
	"ProjectHeis/drivers/config"
	"ProjectHeis/network/bcast"
	"ProjectHeis/network/peers"
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
type OtherHelloMsg struct{
	Message string
	Iter int
}

func main() {
	//Create and asssign ID
	id := config.CreateID()

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
	otherHelloTx := make(chan OtherHelloMsg)
	otherHelloRx := make(chan OtherHelloMsg)
	//Start transmitting and receiving
	go bcast.Transmitter(16569, helloTx)
	go bcast.Receiver(16569, helloRx)

	go bcast.Transmitter(16569, otherHelloTx)
	go bcast.Receiver(16569, otherHelloRx)

	go func() {
		otherMessage := OtherHelloMsg{"HEI" , 0}
		for{
			otherMessage.Iter++
			otherHelloTx <- otherMessage
			time.Sleep(3*time.Second)
		}
	}()

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
		case othermsg := <-otherHelloRx:
			fmt.Printf("Other message: ")
			fmt.Printf("msg: %q", othermsg.Message)
			fmt.Printf("iteration: %d\n", othermsg.Iter)
		}
	}
}

func HallOrderUpdate ()