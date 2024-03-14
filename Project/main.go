package main

import (
	"ProjectHeis/drivers/elevator"
	"ProjectHeis/drivers/eventHandler"
	"ProjectHeis/drivers/fsm"
	"ProjectHeis/network/peers"
)

func main() {
	peers.G_PeersElevator = peers.InitPeers()
	peers.SendPeersData_init()
	go peers.PeersHeartBeat()
	go eventHandler.EventHandling()
	go fsm.Fsm(elevator.G_Ch_requests)

	for {
		select {}
	}
}
