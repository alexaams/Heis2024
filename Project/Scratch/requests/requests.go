package requests

import (
	"ProjectHeis/Scratch/config_folder/globals"
	"ProjectHeis/Scratch/config_folder/types"
	"ProjectHeis/Scratch/drivers/elevator"
)

// Checks current floor to top floor
func IsRequestAbove(elev elevator.Elevator) bool {
	requests := elev.Requests
	for floor := elev.Floor + 1; floor < globals.NumFloors; floor++ {
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
	return elev.Requests.HallUp[elev.Floor] || elev.Requests.HallDown[elev.Floor] || elev.Requests.HallUp[elev.Floor]
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
	for floor := range globals.NumFloors {
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

// func ClearRequestBtnReturn(elev elevator.Elevator) (int, types.ButtonType) {
// 	for buttonType := elevio.BT_HallUp; buttonType < types.ButtonType(globals.NumButtonTypes); buttonType++ {
// 		isRequested := elev.Requests[elev.Floor][buttonType]
// 		isDirectionMatch := (elev.Direction == elevio.MD_Up && buttonType == elevio.BT_HallUp) ||
// 			(elev.Direction == elevio.MD_Down && buttonType == elevio.BT_HallDown)
// 		isCabButton := buttonType == elevio.BT_Cab
// 		isStopped := elev.Direction == elevio.MD_Stop
// 		if isRequested && (isDirectionMatch || isCabButton || isStopped) {
// 			return elev.Floor, buttonType
// 		}
// 	}
// 	// defaults to this as an error indicating that it is stuck at the bottom
// 	fmt.Println("clear btn values")
// 	return -1, elevio.BT_HallUp
// }

func GiveButtonToClear(elev elevator.Elevator) types.ButtonEvent {
	var isRequested bool
	var isDirectionMatch bool
	var isStopped bool
	var isCabButton bool

	for buttonType := range globals.NumButtonTypes {
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

// // Decides where to
// func WhichWay(cuElevator *elevator.Elevator) {
// 	ReqestsAbove := IsRequestAbove(*cuElevator)
// 	RequestsBelow := IsRequestBelow(*cuElevator)

// 	switch {
// 	case ReqestsAbove:
// 		elevio.SetMotorDirection(elevio.MD_Up)
// 		cuElevator.Direction = elevio.MD_Up
// 	case RequestsBelow:
// 		elevio.SetMotorDirection(elevio.MD_Down)
// 		cuElevator.Direction = elevio.MD_Down
// 	default:
// 		elevio.SetMotorDirection(elevio.MD_Stop)
// 		cuElevator.Direction = elevio.MD_Stop
// 	}
// }

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

// func addBtnIfRequested(btnType types.ButtonType, cuElevator elevator.Elevator, btnToClear []types.ButtonEvent) {
// 	if cuElevator.Requests[cuElevator.Floor][btnType] {
// 		btnToClear = append(btnToClear, types.ButtonEvent{Floor: cuElevator.Floor, Button: btnType})
// 	}
// }

func ClearOrders(cuElevator elevator.Elevator) []types.ButtonEvent {
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

	globals.G_Ch_clear_orders <- btnToClear
	return btnToClear
}

// func RequestReadyForClear(elev elevator.Elevator) []types.ButtonEvent {
// 	btnToClear := make([]types.ButtonEvent, 0)

// 	if elev.Requests.CabFloor[elev.Floor] {
// 		btnToClear = append(btnToClear, types.ButtonEvent{Floor: elev.Floor, Button: types.BT_Cab})
// 	}

// 	addBtnIfRequested := func(btnType types.ButtonType) {
// 		if elev.Requests[elev.Floor][btnType] {
// 			btnToClear = append(btnToClear, types.ButtonEvent{Floor: elev.Floor, Button: btnType})
// 		}
// 	}

// 	switch elev.Direction {
// 	case types.MD_Up:
// 		if !IsRequestAbove(elev) {
// 			addBtnIfRequested(types.BT_HallDown)
// 		}
// 		addBtnIfRequested(types.BT_HallUp)
// 	case types.MD_Down:
// 		if !IsRequestBelow(elev) {
// 			addBtnIfRequested(types.BT_HallUp)
// 		}
// 		addBtnIfRequested(types.BT_HallDown)
// 	case types.MD_Stop:
// 		fallthrough
// 	default:
// 		addBtnIfRequested(types.BT_HallUp)
// 		addBtnIfRequested(types.BT_HallDown)
// 	}
// 	return btnToClear
// }

// ----------------- GIVEN ------------------------
// func MakeReqList(amountFloors, botFloor int) config.ReqList {
// 	listFloor := make(map[int]bool)
// 	for x := 0; x < amountFloors; x++ {
// 		listFloor[x+botFloor] = false
// 	}
// 	return listFloor
// }
