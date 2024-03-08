package fms

import (
	"ProjectHeis/drivers/config"
	"ProjectHeis/drivers/elevator"
	"ProjectHeis/drivers/elevio"
	"ProjectHeis/requests"
	"fmt"
)

// channels
var elevBehaviorChan = make(chan elevator.ElevatorBehavior)
var obschan = make(chan bool)

// variables
var d elevio.MotorDirection = elevio.MD_Up
var numFloors config.NumFloors
var cuElevator elevator.Elevator

/*func ButtonSelected(a elevio.ButtonEvent) {
	request_list := config.MakeReqList(4, 0)
	elevio.SetButtonLamp(a.Button, a.Floor, true)
	//Test
	if a.Button == elevio.BT_Cab {
		elevio.SetDoorOpenLamp(false)
		request_list.SetFloor(a.Floor)
	}
}*/
func requestUpdates() {
	var buttonpressed elevio.ButtonEvent
	switch cuElevator.Behavior {
	case elevator.BehaviorMoving:
		elevio.SetMotorDirection(cuElevator.Direction)

	case elevator.BehaviorIdle:
		set := requests.RequestToElevatorMovement(cuElevator)
		cuElevator.Behavior = set.Behavior
		cuElevator.Direction = set.Direction
		elevBehaviorChan <- cuElevator.Behavior
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
			ticks.tickerStart(cuElevator.OpenDuration)
			elevio.SetDoorOpenLamp(true)
			requests.ClearOneRequest(&cuElevator, elevio.CurrentOrder.BtnEvent)
			cuElevator.Behavior = elevator.BehaviorOpen
			elevBehaviorChan <- cuElevator.Behavior

		}
	}
}

func ObstFound() {
	if cuElevator.Behavior == elevator.BehaviorOpen {
		ticks.tickerStart(cuElevator.OpenDuration)
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

func fms() {

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
