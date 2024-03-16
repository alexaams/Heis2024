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
			elevator.Stop()
			elevator.G_this_Elevator.SetElevatorBehaviour(types.BehaviorOpen)
		} else {
			elevio.SetMotorDirection(BevAndDir.Direction)
			elevator.G_this_Elevator.SetElevatorBehaviour(BevAndDir.Behavior)
		}
	} else {
		elevator.Stop()
	}

}

func checkFloorCurrent(a int) {
	elevator.G_this_Elevator.Floor = a
	elevio.SetFloorIndicator(elevator.G_this_Elevator.Floor)
	if requests.IsThisOurStop(elevator.G_this_Elevator) {
		elevator.Stop()
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

func findInitialFloor(drv_floors chan int) {
	timer := time.NewTimer(4 * time.Second)
	stuck := false
	elevator.MoveDown()
	for {
		select {
		case floor := <-drv_floors:
			fmt.Println("Floor found: ", floor)
			elevator.Stop()
			elevator.G_this_Elevator.Floor = floor
			return
		case <-timer.C:
			if stuck {
				elevator.MoveUp()
				stuck = true
			} else {
				elevator.Stop()
				panic("Elevator is stuck!")
			}
		}
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

	findInitialFloor(drv_floors)

	go doorOpen()

	var timer = time.NewTicker(300 * time.Millisecond)
	defer timer.Stop()

	for {
		select {
		case <-timer.C:
			elevator.G_Ch_elevator_update <- elevator.G_this_Elevator
			if elevator.G_this_Elevator.Behavior != types.BehaviorOpen && elevator.G_this_Elevator.Behavior != types.BehaviorMoving {
				requestUpdates()
			}
		case a := <-drv_floors:
			checkFloorCurrent(a)

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
				elevator.Stop()
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

func doorOpen() {
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
			time.Sleep(10 * time.Millisecond)
		case types.BehaviorObst:
		}
	}
}
