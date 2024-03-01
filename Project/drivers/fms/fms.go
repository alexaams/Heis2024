package fms

import (
	"ProjectHeis/drivers/config"
	"ProjectHeis/drivers/elevio"
	"ProjectHeis/ticks"
	"fmt"
)

// variables
var d elevio.MotorDirection = elevio.MD_Up
var numFloors int = 4
var cuElevator config.Elevator



func ButtonSelected(a elevio.ButtonEvent) {
	request_list := config.MakeReqList(4, 0)
	elevio.SetButtonLamp(a.Button, a.Floor, true)
	//Test
	if a.Button == elevio.BT_Cab {
		elevio.SetDoorOpenLamp(false)
		request_list.SetFloor(a.Floor)
	}
}

func FloorCurrent(a int) {
	cuElevator.Floor = a
	elevio.SetFloorIndicator(cuElevator.Floor)
	// elevio.SetMotorDirection(d)
	switch cuElevator.Behavior{
		case config.BehaviorMoving:
			if /*requested should be handled to stop*/{
				elevio.SetMotorDirection(elevio.MD_Stop)
				ticks.tickerStart(cuElevator.OpenDuration)
				
			}
	}
	if elevio.CurrentOrder.BtnEvent.Floor == elevio.GetFloor() {
		elevio.SetButtonLamp(elevio.CurrentOrder.BtnEvent.Button, elevio.CurrentOrder.BtnEvent.Floor, false)
		elevio.SetMotorDirection(elevio.MD_Stop)
		elevio.SetDoorOpenLamp(true)
		elevio.CurrentOrder.Active = false
	}
}

func ObstFound(a bool) {
	fmt.Printf("%+v\n", a)
	if a {
		elevio.SetMotorDirection(elevio.MD_Stop)
	} else {
		elevio.SetMotorDirection(d)
		elevio.SetStopLamp(false)
		elevio.SetMotorDirection(elevio.MD_Stop)
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

/*
func fms() {


	elevio.Init("localhost:15657", numFloors)


	//elevio.SetMotorDirection(d)

	elevio.CurrentOrder.Active = false

	drv_buttons := make(chan elevio.ButtonEvent)
	drv_floors := make(chan int)
	drv_obstr := make(chan bool)
	drv_stop := make(chan bool)
	//awaiting_orders := make(chan elevio.Order)

	go elevio.PollButtons(drv_buttons)         //Channel receives all buttonevents on every floor
	go elevio.PollFloorSensor(drv_floors)      //Channel receives which floor you are at
	go elevio.PollObstructionSwitch(drv_obstr) //Channel receives state for obstruction switch when changed
	go elevio.PollStopButton(drv_stop)         //Channel receives state of stop switch when changed


	for {
		select {
		case a := <-drv_buttons:
			fmt.Printf("%+v\n", a)
			elevio.SetButtonLamp(a.Button, a.Floor, true)

			//Test

			if a.Button == elevio.BT_Cab {
				elevio.SetDoorOpenLamp(false)
				request_list.SetFloor(a.Floor)
			}
			// if (elevio.CurrentOrder.Active) && (a.Button == elevio.BT_Cab) {
			// 	var pending_order elevio.Order
			// 	pending_order.BtnEvent = a
			// 	pending_order.Active = false
			// 	awaiting_orders <- pending_order
			// }

		case a := <-drv_floors:
			fmt.Printf("%+v\n", a)
			// if a == numFloors-1 {
			// 	d = elevio.MD_Down
			// } else if a == 0 {
			// 	d = elevio.MD_Up
			// }
			elevio.SetFloorIndicator(elevio.GetFloor())
			// elevio.SetMotorDirection(d)
			if elevio.CurrentOrder.BtnEvent.Floor == elevio.GetFloor() {
				elevio.SetButtonLamp(elevio.CurrentOrder.BtnEvent.Button, elevio.CurrentOrder.BtnEvent.Floor, false)
				elevio.SetMotorDirection(elevio.MD_Stop)
				elevio.SetDoorOpenLamp(true)
				elevio.CurrentOrder.Active = false
			}

		case a := <-drv_obstr:
			fmt.Printf("%+v\n", a)
			if a {
				elevio.SetMotorDirection(elevio.MD_Stop)
			} else {
				elevio.SetMotorDirection(d)
				elevio.SetStopLamp(false)
				elevio.SetMotorDirection(elevio.MD_Stop)
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
*/
