package requests

import (
	"ProjectHeis/config"
	"ProjectHeis/drivers/elevator"
	"ProjectHeis/drivers/elevio"
)

// Checks current floor to top floor
func IsRequestAbove(elev elevator.Elevator) bool {
	for floor := elev.Floor; floor < config.NumFloors; floor++ {
		for buttonType := 0; buttonType < config.NumButtonTypes; buttonType++ {
			if elev.Requests[floor][buttonType] {
				return true
			}
		}
	}
	return false
}

// Check request from 0 to current floor
func IsRequestBelow(elev elevator.Elevator) bool {
	for floor := 0; floor < elev.Floor; floor++ {
		for buttonType := 0; buttonType < config.NumButtonTypes; buttonType++ {
			if elev.Requests[floor][buttonType] {
				return true
			}
		}
	}
	return false
}

// Checks current floor
func IsRequestArrived(elev elevator.Elevator) bool {
	for buttonType := 0; buttonType < config.NumButtonTypes; buttonType++ {
		if elev.Requests[elev.Floor][buttonType] {
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

func ClearRequestBtnReturn(elev elevator.Elevator) (int, elevio.ButtonType) {
	for buttonType := elevio.BT_HallUp; buttonType < elevio.ButtonType(config.NumButtonTypes); buttonType++ {
		isRequested := elev.Requests[elev.Floor][buttonType]
		isDirectionMatch := (elev.Direction == elevio.MD_Up && buttonType == elevio.BT_HallUp) ||
			(elev.Direction == elevio.MD_Down && buttonType == elevio.BT_HallDown)
		isCabButton := buttonType == elevio.BT_Cab
		isStopped := elev.Direction == elevio.MD_Stop

		if isRequested && (isDirectionMatch || isCabButton || isStopped) {
			return elev.Floor, buttonType
		}
	}
	// defaults to this as an error indicating that it is stuck at the bottom
	return -1, elevio.BT_HallUp
}

func RequestToElevatorMovement(elev elevator.Elevator) elevator.BehaviorAndDirection {
	// Determine request locations relative to the elevator once.
	requestArrived := IsRequestArrived(elev)
	requestAbove := IsRequestAbove(elev)
	requestBelow := IsRequestBelow(elev)

	if requestArrived {
		return elevator.BehaviorAndDirection{Behavior: elevator.BehaviorOpen, Direction: elevio.MD_Stop}
	}

	switch elev.Direction {
	case elevio.MD_Stop:
		if requestAbove {
			return elevator.BehaviorAndDirection{Behavior: elevator.BehaviorMoving, Direction: elevio.MD_Up}
		} else if requestBelow {
			return elevator.BehaviorAndDirection{Behavior: elevator.BehaviorMoving, Direction: elevio.MD_Down}
		}
	case elevio.MD_Up, elevio.MD_Down:

		if requestAbove {
			return elevator.BehaviorAndDirection{Behavior: elevator.BehaviorMoving, Direction: elevio.MD_Up}
		} else if requestBelow {
			return elevator.BehaviorAndDirection{Behavior: elevator.BehaviorMoving, Direction: elevio.MD_Down}
		}
	}
	// Default case when no specific requests dictate movement.
	return elevator.BehaviorAndDirection{Behavior: elevator.BehaviorIdle, Direction: elevio.MD_Stop}
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
