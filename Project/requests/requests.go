package requests

import (
	"ProjectHeis/config"
	"ProjectHeis/drivers/elevator"
	"ProjectHeis/drivers/elevio"
)

// Checks current floor to top floor
func IsRequestAbove(elev elevator.Elevator) bool {
	for floor := elev.Floor; floor < config.NumFloors; floor++ {
		for button := 0; button < config.NumButtonTypes; button++ {
			if elev.Requests[floor][button] {
				return true
			}
		}
	}
	return false
}

// Check request from 0 to current floor
func IsRequestBelow(elev elevator.Elevator) bool {
	for floor := 0; floor < elev.Floor; floor++ {
		for button := 0; button < config.NumButtonTypes; button++ {
			if elev.Requests[floor][button] {
				return true
			}
		}
	}
	return false
}

// Checks current floor
func IsRequestArrived(elev elevator.Elevator) bool {
	for button := 0; button < config.NumButtonTypes; button++ {
		if elev.Requests[elev.Floor][button] {
			return true
		}
	}
	return false
}

func ClearOneRequest(elev *elevator.Elevator, button elevio.ButtonEvent) {
	elev.Requests[button.Floor][button.Button] = false
}

func ClearRequests(elev *elevator.Elevator, buttons []elevio.ButtonEvent) {
	for _, button := range buttons {
		elev.Requests[button.Floor][button.Button] = false
	}
}

func ClearAllRequests(elev *elevator.Elevator) {
	var tempEmptyRequests config.Requests // initialized as default false
	elev.Requests = tempEmptyRequests
}

func RequestToElevatorMovement(elev elevator.Elevator) elevator.BehaviorAndDirection {
	switch elev.Direction {
	case elevio.MD_Stop:
		if IsRequestArrived(elev) {
			return elevator.BehaviorAndDirection{Behavior: elevator.BehaviorOpen, Direction: elevio.MD_Stop}
		} else if IsRequestAbove(elev) {
			return elevator.BehaviorAndDirection{Behavior: elevator.BehaviorMoving, Direction: elevio.MD_Up}
		} else if IsRequestBelow(elev) {
			return elevator.BehaviorAndDirection{Behavior: elevator.BehaviorMoving, Direction: elevio.MD_Down}
		} else {
			return elevator.BehaviorAndDirection{Behavior: elevator.BehaviorIdle, Direction: elevio.MD_Stop}
		}
	case elevio.MD_Up:
		if IsRequestArrived(elev) {
			return elevator.BehaviorAndDirection{Behavior: elevator.BehaviorOpen, Direction: elevio.MD_Stop}
		} else if IsRequestAbove(elev) {
			return elevator.BehaviorAndDirection{Behavior: elevator.BehaviorMoving, Direction: elevio.MD_Up}
		} else if IsRequestBelow(elev) {
			return elevator.BehaviorAndDirection{Behavior: elevator.BehaviorMoving, Direction: elevio.MD_Down}
		} else {
			return elevator.BehaviorAndDirection{Behavior: elevator.BehaviorIdle, Direction: elevio.MD_Stop}
		}
	case elevio.MD_Down:
		if IsRequestArrived(elev) {
			return elevator.BehaviorAndDirection{Behavior: elevator.BehaviorOpen, Direction: elevio.MD_Stop}
		} else if IsRequestAbove(elev) {
			return elevator.BehaviorAndDirection{Behavior: elevator.BehaviorMoving, Direction: elevio.MD_Up}
		} else if IsRequestBelow(elev) {
			return elevator.BehaviorAndDirection{Behavior: elevator.BehaviorMoving, Direction: elevio.MD_Down}
		} else {
			return elevator.BehaviorAndDirection{Behavior: elevator.BehaviorIdle, Direction: elevio.MD_Stop}
		}
	default:
		return elevator.BehaviorAndDirection{Behavior: elevator.BehaviorIdle, Direction: elevio.MD_Stop}
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

// func (r config.ReqList) SetFloor(floor int) {
// 	if _, ok := r[floor]; ok {
// 		r[floor] = true
// 	} else {
// 		fmt.Println("Floor does not exist")
// 	}
// }

// func (r config.ReqList) ClearFloor(floor int) {
// 	if _, ok := r[floor]; ok {
// 		r[floor] = false
// 	} else {
// 		fmt.Println("Floor does not exist")
// 	}
// }
