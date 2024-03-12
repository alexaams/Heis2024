package requests

import (
	"ProjectHeis/config"
	"ProjectHeis/config_folder/globals"
	"ProjectHeis/config_folder/types"
	"ProjectHeis/drivers/elevator"
)

// Checks current floor to top floor
func IsRequestAbove(elev elevator.Elevator) bool {
	for floor := elev.Floor; floor < globals.NumFloors; floor++ {
		for buttonType := 0; buttonType < globals.NumButtonTypes; buttonType++ {
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
		for buttonType := 0; buttonType < globals.NumButtonTypes; buttonType++ {
			if elev.Requests[floor][buttonType] {
				return true
			}
		}
	}
	return false
}

// Checks current floor
func IsRequestArrived(elev elevator.Elevator) bool {
	for buttonType := 0; buttonType < globals.NumButtonTypes; buttonType++ {
		if elev.Requests[elev.Floor][buttonType] {
			return true
		}
	}
	return false
}

func ClearOneRequest(elev *elevator.Elevator, button types.ButtonEvent) {
	elev.Requests[button.Floor][button.Button] = false
}

func ClearRequests(elev *elevator.Elevator, buttons []types.ButtonEvent) {
	for _, button := range buttons {
		elev.Requests[elev.Floor][button.Button] = false
	}
}

func ClearAllRequests(elev *elevator.Elevator) {
	var tempEmptyRequests types.Requests // initialized as default false
	elev.Requests = tempEmptyRequests
}

func ClearRequestBtnReturn(elev elevator.Elevator) (int, types.ButtonType) {
	for buttonType := types.BT_HallUp; buttonType < types.ButtonType(globals.NumButtonTypes); buttonType++ {
		isRequested := elev.Requests[elev.Floor][buttonType]
		isDirectionMatch := (elev.Direction == types.MD_Up && buttonType == types.BT_HallUp) ||
			(elev.Direction == types.MD_Down && buttonType == types.BT_HallDown)
		isCabButton := buttonType == types.BT_Cab
		isStopped := elev.Direction == types.MD_Stop

		if isRequested && (isDirectionMatch || isCabButton || isStopped) {
			return elev.Floor, buttonType
		}
	}
	// defaults to this as an error indicating that it is stuck at the bottom
	return -1, types.BT_HallUp
}

func RequestToElevatorMovement(elev elevator.Elevator) types.BehaviorAndDirection {
	// Determine request locations relative to the elevator once.
	requestArrived := IsRequestArrived(elev)
	requestAbove := IsRequestAbove(elev)
	requestBelow := IsRequestBelow(elev)

	if requestArrived {
		return types.BehaviorAndDirection{Behavior: types.BehaviorOpen, Direction: types.MD_Stop}
	}

	switch elev.Direction {
	case types.MD_Stop:
		if requestAbove {
			return types.BehaviorAndDirection{Behavior: types.BehaviorMoving, Direction: types.MD_Up}
		} else if requestBelow {
			return types.BehaviorAndDirection{Behavior: types.BehaviorMoving, Direction: types.MD_Down}
		}
	case types.MD_Up, types.MD_Down:

		if requestAbove {
			return types.BehaviorAndDirection{Behavior: types.BehaviorMoving, Direction: types.MD_Up}
		} else if requestBelow {
			return types.BehaviorAndDirection{Behavior: types.BehaviorMoving, Direction: types.MD_Down}
		}
	}
	// Default case when no specific requests dictate movement.
	return types.BehaviorAndDirection{Behavior: types.BehaviorIdle, Direction: types.MD_Stop}
}

func RequestReadyForClear(elev elevator.Elevator) []types.ButtonEvent {
	btnToClear := make([]types.ButtonEvent, 0)

	if elev.Requests[elev.Floor][types.BT_Cab] {
		btnToClear = append(btnToClear, types.ButtonEvent{Floor: elev.Floor, Button: types.BT_Cab})
	}

	addBtnIfRequested := func(btnType types.ButtonType) {
		if elev.Requests[elev.Floor][btnType] {
			btnToClear = append(btnToClear, types.ButtonEvent{Floor: elev.Floor, Button: btnType})
		}
	}
	switch elev.Direction {
	case types.MD_Up:
		if !IsRequestAbove(elev) {
			addBtnIfRequested(types.BT_HallDown)
		}
		addBtnIfRequested(types.BT_HallUp)
	case types.MD_Down:
		if !IsRequestBelow(elev) {
			addBtnIfRequested(types.BT_HallUp)
		}
		addBtnIfRequested(types.BT_HallDown)
	case types.MD_Stop:
		fallthrough
	default:
		addBtnIfRequested(types.BT_HallUp)
		addBtnIfRequested(types.BT_HallDown)
	}
	return btnToClear
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
