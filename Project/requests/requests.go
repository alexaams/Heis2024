package requests

import (
	"ProjectHeis/config_folder/config"
	"ProjectHeis/config_folder/types"
	"ProjectHeis/drivers/elevator"
)

// Checks current floor to top floor
func IsRequestAbove(elev elevator.Elevator) bool {
	requests := elev.Requests
	for floor := elev.Floor + 1; floor < config.NumFloors; floor++ {
		if requests.HallUp[floor] || requests.HallDown[floor] || requests.CabFloor[floor] {
			return true
		}
	}
	return false
}

// Check request from 0 to current floor
func IsRequestBelow(elev elevator.Elevator) bool {
	requests := elev.Requests
	for floor := 0; floor < elev.Floor; floor++ {
		if requests.HallUp[floor] || requests.HallDown[floor] || requests.CabFloor[floor] {
			return true
		}
	}
	return false
}

// Checks current floor
func IsThisOurStop(elev elevator.Elevator) bool {
	return elev.Requests.HallUp[elev.Floor] || elev.Requests.HallDown[elev.Floor] || elev.Requests.CabFloor[elev.Floor]
}

func ClearOneRequest(elev *elevator.Elevator, button types.ButtonEvent) {
	switch button.Button {
	case types.BT_HallUp:
		elev.Requests.HallUp[button.Floor] = false
	case types.BT_HallDown:
		elev.Requests.HallDown[button.Floor] = false
	case types.BT_Cab:
		elev.Requests.CabFloor[button.Floor] = false
	}
}

func ClearRequests(elev *elevator.Elevator, buttons []types.ButtonEvent) {
	for _, button := range buttons {
		switch button.Button {
		case types.BT_HallUp:
			elev.Requests.HallUp[button.Floor] = false
		case types.BT_HallDown:
			elev.Requests.HallDown[button.Floor] = false
		case types.BT_Cab:
			elev.Requests.CabFloor[button.Floor] = false
		}
	}
}

func ClearAllRequests(elev *elevator.Elevator) {
	for floor := range config.NumFloors {
		elev.Requests.CabFloor[floor] = false
		elev.Requests.HallUp[floor] = false
		elev.Requests.HallDown[floor] = false
	}
}

func RequestsShouldStop(elev elevator.Elevator) bool {
	if elev.Requests.CabFloor[elev.Floor] {
		return true
	}
	switch elev.Direction {
	case types.MD_Down:
		if elev.Requests.HallDown[elev.Floor] || !IsRequestBelow(elev) {
			return true
		}
	case types.MD_Up:
		if elev.Requests.HallUp[elev.Floor] || !IsRequestAbove(elev) {
			return true
		}
	case types.MD_Stop:
		return true
	}
	return false
}

func GiveButtonToClear(elev elevator.Elevator) types.ButtonEvent {
	var isRequested bool
	var isDirectionMatch bool
	var isStopped bool
	var isCabButton bool

	for buttonType := range config.NumButtonTypes {
		switch types.ButtonType(buttonType) {
		case types.BT_HallUp:
			isRequested = elev.Requests.HallUp[elev.Floor]
			isDirectionMatch = elev.Direction == types.MD_Up
		case types.BT_HallDown:
			isRequested = elev.Requests.HallDown[elev.Floor]
			isDirectionMatch = elev.Direction == types.MD_Down
		case types.BT_Cab:
			isRequested = elev.Requests.CabFloor[elev.Floor]
			isDirectionMatch = elev.Direction == types.MD_Down
			isCabButton = true
		}
		isStopped = elev.Direction == types.MD_Stop
		if isRequested && (isDirectionMatch || isCabButton || isStopped) {
			return types.ButtonEvent{Floor: elev.Floor, Button: types.ButtonType(buttonType)}
		}
	}
	return types.ButtonEvent{Floor: -1, Button: types.BT_HallUp}
}

func RequestToElevatorMovement(elev elevator.Elevator) types.BehaviorAndDirection {
	// Determine request locations relative to the elevator once.
	requestArrived := IsThisOurStop(elev)
	requestAbove := IsRequestAbove(elev)
	requestBelow := IsRequestBelow(elev)
	

	switch elev.Direction {
	case types.MD_Stop:
		if requestAbove {
			return types.BehaviorAndDirection{Behavior: types.BehaviorMoving, Direction: types.MD_Up}
		} else if requestBelow {
			return types.BehaviorAndDirection{Behavior: types.BehaviorMoving, Direction: types.MD_Down}
		} else if requestArrived {
			return types.BehaviorAndDirection{Behavior: types.BehaviorOpen, Direction: types.MD_Down}
		} else {
			return types.BehaviorAndDirection{Behavior: types.BehaviorIdle, Direction: types.MD_Stop}
		}
	case types.MD_Down:

		if requestAbove {
			return types.BehaviorAndDirection{Behavior: types.BehaviorMoving, Direction: types.MD_Up}
		} else if requestBelow {
			return types.BehaviorAndDirection{Behavior: types.BehaviorMoving, Direction: types.MD_Down}
		} else if requestArrived {
			return types.BehaviorAndDirection{Behavior: types.BehaviorOpen, Direction: types.MD_Up}
		} else {
			return types.BehaviorAndDirection{Behavior: types.BehaviorIdle, Direction: types.MD_Stop}
		}
	case types.MD_Up:

		if requestAbove {
			return types.BehaviorAndDirection{Behavior: types.BehaviorMoving, Direction: types.MD_Up}
		} else if requestBelow {
			return types.BehaviorAndDirection{Behavior: types.BehaviorMoving, Direction: types.MD_Down}
		} else if requestArrived {
			return types.BehaviorAndDirection{Behavior: types.BehaviorOpen, Direction: types.MD_Stop}
		} else {
			return types.BehaviorAndDirection{Behavior: types.BehaviorIdle, Direction: types.MD_Stop}
		}
	}
	// Default case when no specific requests dictate movement.
	return types.BehaviorAndDirection{Behavior: types.BehaviorIdle, Direction: types.MD_Stop}
}

func ClearOrders(cuElevator elevator.Elevator) {
	btnToClear := make([]types.ButtonEvent, 0)

	if cuElevator.Requests.CabFloor[cuElevator.Floor] {
		btnToClear = append(btnToClear, types.ButtonEvent{Floor: cuElevator.Floor, Button: types.BT_Cab})
	}

	switch cuElevator.Direction {
	case types.MD_Up:
		if !IsRequestAbove(cuElevator) {
			if cuElevator.Requests.HallDown[cuElevator.Floor] {
				btnToClear = append(btnToClear, types.ButtonEvent{Floor: cuElevator.Floor, Button: types.BT_HallDown})
			}
		}

		if cuElevator.Requests.HallUp[cuElevator.Floor] {
			btnToClear = append(btnToClear, types.ButtonEvent{Floor: cuElevator.Floor, Button: types.BT_HallUp})
		}
	case types.MD_Down:
		if !IsRequestBelow(cuElevator) {
			if cuElevator.Requests.HallUp[cuElevator.Floor] {
				btnToClear = append(btnToClear, types.ButtonEvent{Floor: cuElevator.Floor, Button: types.BT_HallUp})
			}
		}
		if cuElevator.Requests.HallDown[cuElevator.Floor] {
			btnToClear = append(btnToClear, types.ButtonEvent{Floor: cuElevator.Floor, Button: types.BT_HallDown})
		}

	case types.MD_Stop:
		fallthrough
	default:
		if cuElevator.Requests.HallUp[cuElevator.Floor] {
			btnToClear = append(btnToClear, types.ButtonEvent{Floor: cuElevator.Floor, Button: types.BT_HallUp})
		}
		if cuElevator.Requests.HallDown[cuElevator.Floor] {
			btnToClear = append(btnToClear, types.ButtonEvent{Floor: cuElevator.Floor, Button: types.BT_HallDown})
		}
	}

	elevator.G_Ch_clear_orders <- btnToClear
}
