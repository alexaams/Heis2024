package fsm2

import (
	"ProjectHeis/config"
	"ProjectHeis/drivers/elevator"
	"ProjectHeis/drivers/elevio"
	"ProjectHeis/network/peers"
)

var newOrderChan = make(chan bool)
var openDoorChan = make(chan bool)
var shouldStopchan = make(chan bool)
var obstChan = make(chan bool)
var buttonEventChan = make(chan bool)

var elevator_obj = elevator.InitElevator()
var peers_obj = peers.InitPeers()

var peersDataMap = make(map[int]peers.PeersData)


func fsmDriver(
	elev *elevator.Elevator,
	newOrderChan chan bool,
	openDoorChan chan bool,
	shouldStopChan chan bool,
){

	switch elev.Behavior {
	case elevator.BehaviorIdle:
			
		
	}
}



func lampChange() {
	for f := range config.NumFloors {
		for b := range config.NumButtonTypes {
			elevio.SetButtonLamp(elevio.ButtonType(b), f, elevator_obj.Requests[f][b])
		}
	}
}
