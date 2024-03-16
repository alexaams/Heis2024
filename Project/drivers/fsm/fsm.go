package fsm

import (
	"ProjectHeis/config_folder/config"
	"ProjectHeis/config_folder/types"
	"ProjectHeis/drivers/elevator"
	"ProjectHeis/drivers/elevio"
	"ProjectHeis/requests"
	"fmt"
	"time"
)

func requestUpdates() {
	BevAndDir := requests.RequestToElevatorMovement(elevator.G_this_Elevator)
	if elevator.G_this_Elevator.Behavior != types.BehaviorOpen && elevator.G_this_Elevator.Behavior != types.BehaviorObst {
		if requests.IsThisOurStop(elevator.G_this_Elevator) {
			elevator.G_this_Elevator.Stop()
			elevator.G_this_Elevator.SetElevatorBehaviour(types.BehaviorOpen)
		} else {
			elevio.SetMotorDirection(BevAndDir.Direction)
			elevator.G_this_Elevator.SetElevatorBehaviour(BevAndDir.Behavior)
		}
	} else {
		elevio.SetMotorDirection(types.MD_Stop)
	}

}

func CheckFloorCurrent(a int) {
	elevator.G_this_Elevator.Floor = a
	elevio.SetFloorIndicator(elevator.G_this_Elevator.Floor)
	if requests.IsThisOurStop(elevator.G_this_Elevator) {
		elevator.G_this_Elevator.Stop()
		elevio.SetDoorOpenLamp(true)
		elevator.G_this_Elevator.SetElevatorBehaviour(types.BehaviorOpen)
	}
}

func mapNewRequests(reqs types.Requests) {
	for i := 0; i < config.NumFloors; i++ {
		elevator.G_this_Elevator.Requests.HallUp[i] = reqs.HallUp[i]
		elevator.G_this_Elevator.Requests.HallDown[i] = reqs.HallDown[i]
		elevator.G_this_Elevator.Requests.CabFloor[i] = reqs.CabFloor[i]

	}
}

func initFloorReading(drv_floors chan int) {
	firstFloorReading := elevio.GetFloor()
	if firstFloorReading == -1 {
		fmt.Println("Undefined floor")
		elevio.SetMotorDirection(types.MD_Up)
		for {
			select {
			case a := <-drv_floors:
				fmt.Printf("\nFound floor: %d\n", a)
				elevio.SetMotorDirection(types.MD_Stop)
				elevator.G_this_Elevator.Floor = a
				return
			case <-time.After(time.Second * 5):
				panic("Elevator Stuck, shiiit!")
			default:
				fmt.Print(".")
				time.Sleep(30 * time.Millisecond)
			}
		}
	} else {
		fmt.Printf("Started at defined floor: %d\n", firstFloorReading)
	}

}

func Fsm(ch_requests chan types.Requests) {

	elevio.Init("localhost:15657", config.NumFloors)
	fmt.Print("Initiating FSM...")
	drv_floors := make(chan int)
	drv_obstr := make(chan bool)
	drv_stop := make(chan bool)

	go elevio.PollFloorSensor(drv_floors)
	go elevio.PollObstructionSwitch(drv_obstr)
	go elevio.PollStopButton(drv_stop)

	initFloorReading(drv_floors)

	go StateMachineBehavior()

	var timer = time.NewTicker(600 * time.Millisecond)
	defer timer.Stop()

	for {
		select {
		case <-timer.C:
			elevator.G_Ch_elevator_update <- elevator.G_this_Elevator
			if elevator.G_this_Elevator.Behavior != types.BehaviorOpen && elevator.G_this_Elevator.Behavior != types.BehaviorMoving {
				requestUpdates()
			}
		case a := <-drv_floors:
			CheckFloorCurrent(a)

		case a := <-drv_obstr:
			fmt.Printf("%+v\n", a)
			if a {
				fmt.Print("Obstruction\n")
				elevator.G_this_Elevator.SetElevatorBehaviour(types.BehaviorObst)
			} else {
				fmt.Print("Obstruction removed\n")
				elevator.G_this_Elevator.SetElevatorBehaviour(types.BehaviorIdle)
			}

		case a := <-drv_stop:
			fmt.Printf("Stopbutton: %+v\n", a)

			if a {
				elevio.SetStopLamp(true)
				elevio.SetMotorDirection(types.MD_Stop)
			} else {
				elevio.SetStopLamp(false)
			}
		case requests := <-ch_requests:
			fmt.Println(requests)
			mapNewRequests(requests)
			fmt.Println("New orders")

		}
	}
}

func StateMachineBehavior() {
	elevator.G_door_open_counter = 0
	clearOrderFlag := true
	for {
		switch elevator.G_this_Elevator.Behavior {
		case types.BehaviorOpen:
			elevio.SetDoorOpenLamp(true)
			if clearOrderFlag {
				elevator.G_door_open_counter = 0
				requests.ClearOrders(elevator.G_this_Elevator)
				clearOrderFlag = false
			}
			elevator.G_door_open_counter++
			if elevator.G_door_open_counter > elevator.G_ticks {
				elevio.SetDoorOpenLamp(false)
				elevator.G_this_Elevator.SetElevatorBehaviour(types.BehaviorIdle)
				elevator.G_door_open_counter = 0
				clearOrderFlag = true
			}
		case types.BehaviorObst:
			fmt.Print("obstruction - state machine!")
		}
	}
}
