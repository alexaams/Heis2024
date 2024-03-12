package requests

import (
	"ProjectHeis/config"
	"ProjectHeis/drivers/elevator"
	"ProjectHeis/drivers/elevio"
	"fmt"
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

func IsThisOurStop(elev *elevator.Elevator) bool {
	return elev.Requests[elev.Floor][0] || elev.Requests[elev.Floor][1] || elev.Requests[elev.Floor][2]
}

func ClearOneRequest(elev *elevator.Elevator, button elevio.ButtonEvent) elevator.Elevator {
	elev.Requests[button.Floor][button.Button] = false
	return *elev
}

func ClearRequests(elev *elevator.Elevator, buttons []elevio.ButtonEvent) {
	for _, button := range buttons {
		elev.Requests[elev.Floor][button.Button] = false
	}
}

func ClearAllRequests(elev *elevator.Elevator) {
	var tempEmptyRequests config.Requests // initialized as default false
	elev.Requests = tempEmptyRequests
}

func RequestsShouldStop(e elevator.Elevator) bool {
	if e.Requests[e.Floor][elevio.BT_Cab] {
		return true
	}
	switch e.Direction {
	case elevio.MD_Down:
		if e.Requests[e.Floor][elevio.BT_HallDown] || !IsRequestBelow(e) {
			return true
		}
	case elevio.MD_Up:
		if e.Requests[e.Floor][elevio.BT_HallUp] || !IsRequestAbove(e) {
			return true
		}
	case elevio.MD_Stop:
		return true
	}
	return false
}

func ClearRequestBtnReturn(elev elevator.Elevator) (int, elevio.ButtonType) {
	for buttonType := elevio.BT_HallUp; buttonType < elevio.ButtonType(config.NumButtonTypes); buttonType++ {
		isRequested := elev.Requests[elev.Floor][buttonType]
		isDirectionMatch := (elev.Direction == elevio.MD_Up && buttonType == elevio.BT_HallUp) ||
			(elev.Direction == elevio.MD_Down && buttonType == elevio.BT_HallDown)
		isCabButton := buttonType == elevio.BT_Cab
		isStopped := elev.Direction == elevio.MD_Stop
		fmt.Println("is requested value: ", isRequested)
		fmt.Println("is direction value: ", isDirectionMatch)
		fmt.Println("is cab value: ", isCabButton)
		fmt.Println("is stopped value", isStopped)

		if isRequested && (isDirectionMatch || isCabButton || isStopped) {
			fmt.Println("clear btn values", elev.Floor, buttonType)
			return elev.Floor, buttonType
		}
	}
	// defaults to this as an error indicating that it is stuck at the bottom
	fmt.Println("clear btn values")
	return -1, elevio.BT_HallUp
}

// Decides where to
func WhichWay(cuElevator *elevator.Elevator) {
	ReqestsAbove := IsRequestAbove(*cuElevator)
	RequestsBelow := IsRequestBelow(*cuElevator)

	switch {
	case ReqestsAbove:
		elevio.SetMotorDirection(elevio.MD_Up)
		cuElevator.Direction = elevio.MD_Up
	case RequestsBelow:
		elevio.SetMotorDirection(elevio.MD_Down)
		cuElevator.Direction = elevio.MD_Down
	default:
		elevio.SetMotorDirection(elevio.MD_Stop)
		cuElevator.Direction = elevio.MD_Stop
	}
}

func RequestToElevatorMovement(elev elevator.Elevator) elevator.BehaviorAndDirection {
	// Determine request locations relative to the elevator once.
	requestArrived := IsRequestArrived(elev)
	requestAbove := IsRequestAbove(elev)
	requestBelow := IsRequestBelow(elev)

	switch elev.Direction {
	case elevio.MD_Stop:
		if requestAbove {
			return elevator.BehaviorAndDirection{Behavior: elevator.BehaviorMoving, Direction: elevio.MD_Up}
		} else if requestBelow {
			return elevator.BehaviorAndDirection{Behavior: elevator.BehaviorMoving, Direction: elevio.MD_Down}
		} else if requestArrived {
			return elevator.BehaviorAndDirection{Behavior: elevator.BehaviorOpen, Direction: elevio.MD_Down}
		} else {
			return elevator.BehaviorAndDirection{Behavior: elevator.BehaviorIdle, Direction: elevio.MD_Stop}
		}
	case elevio.MD_Down:

		if requestAbove {
			return elevator.BehaviorAndDirection{Behavior: elevator.BehaviorMoving, Direction: elevio.MD_Up}
		} else if requestBelow {
			return elevator.BehaviorAndDirection{Behavior: elevator.BehaviorMoving, Direction: elevio.MD_Down}
		} else if requestArrived {
			return elevator.BehaviorAndDirection{Behavior: elevator.BehaviorOpen, Direction: elevio.MD_Up}
		} else {
			return elevator.BehaviorAndDirection{Behavior: elevator.BehaviorIdle, Direction: elevio.MD_Stop}
		}
	case elevio.MD_Up:

		if requestAbove {
			return elevator.BehaviorAndDirection{Behavior: elevator.BehaviorMoving, Direction: elevio.MD_Up}
		} else if requestBelow {
			return elevator.BehaviorAndDirection{Behavior: elevator.BehaviorMoving, Direction: elevio.MD_Down}
		} else if requestArrived {
			return elevator.BehaviorAndDirection{Behavior: elevator.BehaviorOpen, Direction: elevio.MD_Stop}
		} else {
			return elevator.BehaviorAndDirection{Behavior: elevator.BehaviorIdle, Direction: elevio.MD_Stop}
		}
	}
	// Default case when no specific requests dictate movement.
	return elevator.BehaviorAndDirection{Behavior: elevator.BehaviorIdle, Direction: elevio.MD_Stop}
}

func addBtnIfRequested(btnType elevio.ButtonType, cuElevator elevator.Elevator, btnToClear []elevio.ButtonEvent) {
	if cuElevator.Requests[cuElevator.Floor][btnType] {
		btnToClear = append(btnToClear, elevio.ButtonEvent{Floor: cuElevator.Floor, Button: btnType})
	}
}
func ClearOrders(cuElevator elevator.Elevator) {
	btnToClear := make([]elevio.ButtonEvent, 0)

	if cuElevator.Requests[cuElevator.Floor][elevio.BT_Cab] {
		btnToClear = append(btnToClear, elevio.ButtonEvent{Floor: cuElevator.Floor, Button: elevio.BT_Cab})
	}

	switch cuElevator.Direction {
	case elevio.MD_Up:
		if !IsRequestAbove(cuElevator) {
			if cuElevator.Requests[cuElevator.Floor][elevio.BT_HallDown] {
				btnToClear = append(btnToClear, elevio.ButtonEvent{Floor: cuElevator.Floor, Button: elevio.BT_HallDown})
			}
		}

		if cuElevator.Requests[cuElevator.Floor][elevio.BT_HallUp] {
			btnToClear = append(btnToClear, elevio.ButtonEvent{Floor: cuElevator.Floor, Button: elevio.BT_HallUp})
		}
	case elevio.MD_Down:
		if !IsRequestBelow(cuElevator) {
			if cuElevator.Requests[cuElevator.Floor][elevio.BT_HallUp] {
				btnToClear = append(btnToClear, elevio.ButtonEvent{Floor: cuElevator.Floor, Button: elevio.BT_HallUp})
			}

		}
		if cuElevator.Requests[cuElevator.Floor][elevio.BT_HallDown] {
			btnToClear = append(btnToClear, elevio.ButtonEvent{Floor: cuElevator.Floor, Button: elevio.BT_HallDown})
		}

	case elevio.MD_Stop:
		fallthrough
	default:
		if cuElevator.Requests[cuElevator.Floor][elevio.BT_HallUp] {
			btnToClear = append(btnToClear, elevio.ButtonEvent{Floor: cuElevator.Floor, Button: elevio.BT_HallUp})
		}
		if cuElevator.Requests[cuElevator.Floor][elevio.BT_HallDown] {
			btnToClear = append(btnToClear, elevio.ButtonEvent{Floor: cuElevator.Floor, Button: elevio.BT_HallDown})
		}

	}

	config.G_Ch_clear_orders <- btnToClear
}

func RequestReadyForClear(elev elevator.Elevator) []elevio.ButtonEvent {
	btnToClear := make([]elevio.ButtonEvent, 0)

	if elev.Requests[elev.Floor][elevio.BT_Cab] {
		btnToClear = append(btnToClear, elevio.ButtonEvent{Floor: elev.Floor, Button: elevio.BT_Cab})
	}

	addBtnIfRequested := func(btnType elevio.ButtonType) {
		if elev.Requests[elev.Floor][btnType] {
			btnToClear = append(btnToClear, elevio.ButtonEvent{Floor: elev.Floor, Button: btnType})
		}
	}

	switch elev.Direction {
	case elevio.MD_Up:
		if !IsRequestAbove(elev) {
			addBtnIfRequested(elevio.BT_HallDown)
		}
		addBtnIfRequested(elevio.BT_HallUp)
	case elevio.MD_Down:
		if !IsRequestBelow(elev) {
			addBtnIfRequested(elevio.BT_HallUp)
		}
		addBtnIfRequested(elevio.BT_HallDown)
	case elevio.MD_Stop:
		fallthrough
	default:
		addBtnIfRequested(elevio.BT_HallUp)
		addBtnIfRequested(elevio.BT_HallDown)
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
