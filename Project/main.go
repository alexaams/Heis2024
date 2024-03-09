package main

import (
	"ProjectHeis/drivers/config"
	"ProjectHeis/drivers/elevio"
	"ProjectHeis/network/bcast"
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
	//Set all button lamps to off
	//REMOVE THIS LATER!
	//__________________________
	elevio.SetButtonLamp(elevio.BT_HallUp, 0, false)
	elevio.SetButtonLamp(elevio.BT_HallDown, 1, false)
	elevio.SetButtonLamp(elevio.BT_HallUp, 1, false)
	elevio.SetButtonLamp(elevio.BT_HallDown, 2, false)
	elevio.SetButtonLamp(elevio.BT_HallUp, 2, false)
	elevio.SetButtonLamp(elevio.BT_HallDown, 3, false)
	//__________________________

	//Create and asssign ID
	id := config.CreateID()
	fmt.Printf("ID: %d\n", id)

	//Create global order-table
	globalOrderTable := config.CreateGlobalOrderTable()
	//Create channel for button-event
	ch_HallButton_event := make(chan elevio.ButtonEvent)
	//Create channel for button-polling
	button_channel := make(chan elevio.ButtonEvent)
	//Create channel for transmitting globalordertable over UDP
	udp_GlobalOrder_Tx := make(chan elevio.ButtonEvent)
	//Create channel for receiving globalordertable over UDP
	udp_GLobalOrder_Rx := make(chan elevio.ButtonEvent)

	//Startup-procedure--------------------------------

	//Create channel for sensor-polling
	sensor_channel := make(chan int)
	//Running thread checking ch_HallButton_event
	go UpdateGlobalData(&globalOrderTable, ch_HallButton_event, udp_GlobalOrder_Tx)
	go elevio.PollButtons(button_channel)
	go elevio.PollFloorSensor(sensor_channel)
	go bcast.Transmitter(16569, udp_GlobalOrder_Tx)
	go bcast.Receiver(16569, udp_GLobalOrder_Rx)
	//go elevio.PollFloorSensor(sensor_channel)

	for {
		select {
		case a := <-button_channel:
			switch a.Button {
			case elevio.BT_HallUp, elevio.BT_HallDown:
				ch_HallButton_event <- a

			default:
				fmt.Println("Nothing happens")
			}
		case a := <-udp_GLobalOrder_Rx:
			fmt.Println("Order received over UDP")
			a.PrintButtonEvent()
			if IsOrderNew(&globalOrderTable, a) {
				elevio.SetButtonLamp(a.Button, a.Floor, true)
				fmt.Println("Order is new to this node")
				elevio.SetButtonLamp(a.Button, a.Floor, true)
			}
		case a := <-sensor_channel:

			globalOrderTable[a][0].Active = false
			globalOrderTable[a][0].ElevatorID = -1

		}
	}

}

// ____________________________________________________________
// Må flyttes senere
func UpdateGlobalData(GlobalTable *config.GlobalOrderTable, ch_HallBtn chan elevio.ButtonEvent, udp_GlobalOrder chan elevio.ButtonEvent) {
	for {
		select {
		case a := <-ch_HallBtn:
			switch a.Button {
			case elevio.BT_HallUp, elevio.BT_HallDown:
				if IsOrderNew(GlobalTable, a) {
					udp_GlobalOrder <- a
					elevio.SetButtonLamp(a.Button, a.Floor, true)
				}
				time.Sleep(10 * time.Millisecond)
			default:
				fmt.Printf("Button type is N/A\n")
				time.Sleep(10 * time.Millisecond)
			}
		default:
			continue
		}
	}

}

func IsOrderNew(GlobalTable *config.GlobalOrderTable, Button elevio.ButtonEvent) bool {
	if !GlobalTable[Button.Floor][Button.Button].Active {
		GlobalTable[Button.Floor][Button.Button].Active = true
		return true
	} else {
		return false
	}
}

func UDP_SendRead_GlobalOrder(udp_GlobalOrder chan config.GlobalOrderTable) {
	for {
		select {
		case a := <-udp_GlobalOrder:
			fmt.Println("Order sent for sending over UDP")
			a.PrintGlobalOrderTable()
		default:
			continue
		}
	}
}

//NOT TESTED_______________________________________________________

func BackupRequestHandler(isMaster bool, reqChanTx chan config.GlobalOrderTable, reqChanRx chan bool, GlobalTable config.GlobalOrderTable) {
	if isMaster {
		for {
			select {
			case <-reqChanRx:
				reqChanTx <- GlobalTable
			default:
				continue
			}
		}
	}
}

func ReceiveBackup(askedForBackup bool, reqChanRx chan config.GlobalOrderTable, GlobalTable *config.GlobalOrderTable) {
	if askedForBackup {
		for {
			select {
			case a := <-reqChanRx:
				fmt.Println("Backup received")
				*GlobalTable = a
			default:
				continue
			}
		}
	}
}

func SendBackupRequest(askedForBackup *bool, reqChanTx chan bool) {
	reqChanTx <- true
	*askedForBackup = true
	time.Sleep(10 * time.Millisecond)
}

func InitiatingProcedure(id int) {

}

/*
ORDNE:
* Lag en initieringsprosess som henter ut liste over alle noder i nettverket.

* Deretter skal den sende en request om å få ordreliste.
* Lag og test funksjoner for å sende og motta ordreliste.
* Se om du kan gjøre main-koden mer kompakt ved å deklarere kanaler inne i funksjonene.
* Få et overblikk over andre løsninger, for å se hvor mange kanaler over UDP som er i bruk der.
* Finn ut hva slags kø-system som burde brukes (ikke bruke kø-system?)
* Busy-waiting, risikerer vi det når vi har default-case? ChatGpT sier så.
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