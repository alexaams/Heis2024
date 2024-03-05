package requests

import (
	"ProjectHeis/drivers/config"
	"ProjectHeis/drivers/elevator"
	"ProjectHeis/drivers/elevio"
	"fmt"
)

// Checks current floor to top floor
func IsRequestAbove(elev elevator.Elevator) bool {
	for floor := elev.Floor + 1; floor <= config.NumFloors; floor++ {
		for i := 0; i < config.NumButtons; i++ {
			if elev.Requests[elev.Floor][i] {
				return true
			}
		}
	}
	return false
}

// Check request from 0 to current floor
func IsRequestBelow(elev elevator.Elevator) bool {
	for floor := 0; floor < elev.Floor; floor++ {
		for i := 0; i < config.NumButtons; i++ {
			if elev.Requests[elev.Floor][i] {
				return true
			}
		}
	}
	return false
}

// Checks current floor
func IsRequestArrived(elev elevator.Elevator) bool {
	for i := 0; i < config.NumButtons; i++ {
		if elev.Requests[elev.Floor][i] {
			return true
		}
	}
	return false
}

func ClearOneRequest(elev elevator.Elevator, button elevio.ButtonEvent) {
	elev.Requests[button.Floor][button.Button] = false
}

func ClearAllRequests(elev elevator.Elevator) {
	for floor := 0; floor < config.NumFloors; floor++ {
		for button := 0; button < config.NumButtons; button++ {
			elev.Requests[floor][button] = false
		}
	}
}

// ----------------- GIVEN ------------------------
func MakeReqList(amountFloors, botFloor int) config.ReqList {
	listFloor := make(map[int]bool)
	for x := 0; x < amountFloors; x++ {
		listFloor[x+botFloor] = false
	}
	return listFloor
}

func (r config.ReqList) SetFloor(floor int) {
	if _, ok := r[floor]; ok {
		r[floor] = true
	} else {
		fmt.Println("Floor does not exist")
	}
}

func (r config.ReqList) ClearFloor(floor int) {
	if _, ok := r[floor]; ok {
		r[floor] = false
	} else {
		fmt.Println("Floor does not exist")
	}
}
