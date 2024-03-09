package main

import (
	"ProjectHeis/config"
	"ProjectHeis/drivers/elevio"
	"ProjectHeis/network/peers"
)

func main() {
	//Initiating elevator
	go elevio.Init("localhost:15657", config.NumFloors)
	//Creating ID and initiating heartbeat
	go peers.PeersHeartBeat()
	//Initiate PeersData
	go peers.SendPeersData_init()

	//Få i gang polling av alle knapper, etc

	//Disse knappene skal polles, og alt som endrer på status til en heis, skal sendes til peers.G_Ch_PeersData_Tx

	//peers.G_Ch_PeersData_Rx skal inn i kost-funksjon. Kostfunksjon skal da kjøre algoritmen sin, og returnere ordre,
	//disse ordre sendes i annen kanal? samme kanal?

	//alle heiser tar inn ordre fra over inn i sin heis-modul, og derfra kjøres en lokal algoritme?: finn ut av denne.

	for {
		select {}
	}
}

// Må flyttes senere
/*
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
*/
