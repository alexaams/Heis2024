package debugprint

import (
	"ProjectHeis/config_folder/types"
	"ProjectHeis/drivers/elevator"
	"fmt"
)

func PrintElevatorbehaviour(elev elevator.Elevator) {
	switch elev.Behavior {
	case types.BehaviorIdle:
		fmt.Println("Elevator Behaviour: IDLE")
	case types.BehaviorMoving:
		fmt.Println("Elevator Behaviour: MOVING")
	case types.BehaviorOpen:
		fmt.Println("Elevator Behaviour: OPEN")
	case types.BehaviorObst:
		fmt.Println("Elevator Behaviour: OBST")
	}
}

func PrintElevatorDirection(elev elevator.Elevator) {
	switch elev.Direction {
	case types.MD_Up:
		fmt.Println("Elevator Direction: UP")
	case types.MD_Down:
		fmt.Println("Elevator Direction: UP")
	case types.MD_Stop:
		fmt.Println("Elevator Direction: UP")
	}
}

func PrintElevatorRequests(elev elevator.Elevator) {
	fmt.Println("Hall Up Requests: ", elev.Requests.HallUp)
	fmt.Println("Hall Down Requests: ", elev.Requests.HallDown)
	fmt.Println("Cab Requests: ", elev.Requests.CabFloor)
}

func PrintElevatorFloor(elev elevator.Elevator) {
	fmt.Println("Current floor: ", elev.Floor)
}

func PrintElevator(elev elevator.Elevator) {
	PrintElevatorFloor(elev)
	PrintElevatorDirection(elev)
	PrintElevatorbehaviour(elev)
	PrintElevatorRequests(elev)
}
