package fms

import (
	"ProjectHeis/config"
	"ProjectHeis/drivers/elevator"
	"ProjectHeis/drivers/elevio"
	"ProjectHeis/network/peers"
	"ProjectHeis/requests"
	"ProjectHeis/ticker"
	"fmt"
	"time"
)

// channels
var elevBehaviorChan = make(chan elevator.ElevatorBehavior)
var obschan = make(chan bool)
var peerMsgChan = make(chan peers.PeersData)

// variables
// var d elevio.MotorDirection = elevio.MD_Up
var numFloors = config.NumFloors
var cuElevator elevator.Elevator

// func ButtonSelected(a elevio.ButtonEvent) {
// 	request_list := requests.MakeReqList(4, 0)
// 	elevio.SetButtonLamp(a.Button, a.Floor, true)
// 	//Test
// 	if a.Button == elevio.BT_Cab {
// 		elevio.SetDoorOpenLamp(false)
// 		request_list.SetFloor(a.Floor)
// 	}
// }

func requestUpdates() {
	var buttonpressed elevio.ButtonEvent
	switch cuElevator.Behavior {
	case elevator.BehaviorOpen:
		if floor, buttonType := requests.ClearRequestBtnReturn(cuElevator); floor < -1 {
			ticker.TickerStart(cuElevator.OpenDuration)
			buttonpressed.Button = buttonType
			buttonpressed.Floor = floor
			requests.ClearOneRequest(&cuElevator, buttonpressed)
		}

	case elevator.BehaviorIdle:
		set := requests.RequestToElevatorMovement(cuElevator)
		cuElevator.Behavior = set.Behavior
		cuElevator.Direction = set.Direction
		elevBehaviorChan <- cuElevator.Behavior
		switch set.Behavior {
		case elevator.BehaviorOpen:
			elevio.SetDoorOpenLamp(true)
			ticker.TickerStart(cuElevator.OpenDuration)
			requests.ClearOneRequest(&cuElevator, elevio.CurrentOrder.BtnEvent)

		case elevator.BehaviorMoving:
			elevio.SetMotorDirection(cuElevator.Direction)
		}

	}
}

func FloorCurrent(a int) {
	cuElevator.Floor = a
	elevio.SetFloorIndicator(cuElevator.Floor)
	// elevio.SetMotorDirection(d)
	switch cuElevator.Behavior {
	case elevator.BehaviorMoving:
		if requests.IsRequestArrived(cuElevator) {
			elevio.SetMotorDirection(elevio.MD_Stop)
			ticker.TickerStart(cuElevator.OpenDuration)
			elevio.SetDoorOpenLamp(true)
			requests.ClearOneRequest(&cuElevator, elevio.CurrentOrder.BtnEvent)
			cuElevator.Behavior = elevator.BehaviorOpen
			elevBehaviorChan <- cuElevator.Behavior

		}
	}
}

func ObstFound() {
	if cuElevator.Behavior == elevator.BehaviorOpen {
		ticker.TickerStart(cuElevator.OpenDuration)
		obschan <- true
	}
}

func StopFound(a bool) {
	fmt.Printf("%+v\n", a)
	for f := 0; f < numFloors; f++ {
		for b := elevio.ButtonType(0); b < 3; b++ {
			elevio.SetButtonLamp(b, f, false)
		}
	}
	if a {
		elevio.SetStopLamp(true)
		elevio.SetMotorDirection(elevio.MD_Stop)
	}
}

func fms(hallOrderChan chan config.OrdersHall, orderChan chan []bool) {

	elevio.Init("localhost:15657", numFloors)

	//elevio.SetMotorDirection(d)

	elevio.CurrentOrder.Active = false

	drv_floors := make(chan int)
	drv_obstr := make(chan bool)
	drv_stop := make(chan bool)
	//awaiting_orders := make(chan elevio.Order)
	//Channel receives all buttonevents on every floor
	go elevio.PollFloorSensor(drv_floors)      //Channel receives which floor you are at
	go elevio.PollObstructionSwitch(drv_obstr) //Channel receives state for obstruction switch when changed
	go elevio.PollStopButton(drv_stop)         //Channel receives state of stop switch when changed

	for {
		select {
		case a := <-drv_floors:
			FloorCurrent(a)

		case a := <-drv_obstr:
			fmt.Printf("%+v\n", a)
			if a {
				ObstFound()
			}

		case a := <-drv_stop:
			fmt.Printf("%+v\n", a)
			for f := 0; f < numFloors; f++ {
				for b := elevio.ButtonType(0); b < 3; b++ {
					elevio.SetButtonLamp(b, f, false)
				}
			}
			if a {
				elevio.SetStopLamp(true)
				elevio.SetMotorDirection(elevio.MD_Stop)
			}
		}
	}
}

func lampChange() {
	for floors := range config.NumFloors {
		for buttons := range config.NumButtonTypes - 1 {
			elevio.SetButtonLamp(elevio.ButtonType(buttons), floors, cuElevator.Requests[floors][buttons])
		}
		elevio.SetButtonLamp(elevio.BT_Cab, floors, cuElevator.CabRequests[floors])
	}
}

func eventHandling(orderChan chan []bool) {
	var (
		hallOrderChan  = make(chan config.OrdersHall)
		elevUpdateChan = make(chan elevator.Elevator)
		//bcastReadChan     = make(chan peers.PeersData)
		orderCompleteChan = make(chan elevio.ButtonEvent)
		drv_buttons       = make(chan elevio.ButtonEvent)
		timer             = time.NewTicker(300 * time.Millisecond)
	)
	defer timer.Stop()

	go elevio.PollButtons(drv_buttons)
	//go bcast.Transmitter(peerMsgChan)
	//go bcast.Receiver(bcastReadChan)
	go fms(hallOrderChan, orderChan)

	for {
		select {
		case <-timer.C:
			if len(peers.PeerUpdate.Lost) > 0 {

			}
		}
	}

}
