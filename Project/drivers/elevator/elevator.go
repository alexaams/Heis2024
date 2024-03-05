package elevator

import (
	"ProjectHeis/drivers/config"
	"ProjectHeis/drivers/elevio"
)

// --------------------------------VALUES--------------------------------

type ElevBehavior int

const (
	BehaviorIdle ElevBehavior = iota
	BehaviorMoving
	BehaviorOpen
	BehavoirObst
)

const doorOpenDuration = 3.0

// --------------------------------TYPES--------------------------------

type Elevator struct {
	Floor        int
	Direction    elevio.MotorDirection
	Behavior     ElevBehavior // 0:idle, 1:open, 2:moving, 3: obst
	OpenDuration float64
	Requests     config.Requests
}

// --------------------------------FUNCTIONS--------------------------------

func ElevatorBehaviorToString(elev Elevator) string {
	behavior := elev.Behavior
	switch behavior {
	case BehaviorIdle:
		return "idle"
	case BehaviorMoving:
		return "moving"
	case BehaviorOpen:
		return "open"
	case BehavoirObst:
		return "obst"
	default:
		return "undefined"
	}
}

func ElevatorDirectionToString(elev Elevator) string {
	dir := elev.Direction
	switch dir {
	case elevio.MD_Stop:
		return "stop"
	case elevio.MD_Up:
		return "up"
	case elevio.MD_Down:
		return "down"
	default:
		return "undefined"
	}
}

func InitElevator() Elevator {
	return Elevator{
		Direction:    elevio.MD_Stop,
		Floor:        -1,
		Behavior:     BehaviorIdle,
		OpenDuration: doorOpenDuration,
	}
}

func SetElevatorFloor(elev *Elevator, floor int) {
	elev.Floor = floor
}

func SetElevatorDir(elev *Elevator, dir elevio.MotorDirection) {
	elev.Direction = dir
}

func SetElevatorBehaviour(elev *Elevator, behavior ElevBehavior) {
	elev.Behavior = behavior
}

func SetElevatorOpenDuration(elev *Elevator, time_s float64) {
	elev.OpenDuration = time_s
}
