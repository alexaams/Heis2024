package elevator

import (
	"ProjectHeis/config"
	"ProjectHeis/drivers/elevio"
)

// --------------------------------VALUES--------------------------------

type ElevatorBehavior int

const (
	BehaviorIdle ElevatorBehavior = iota
	BehaviorMoving
	BehaviorOpen
	BehaviorObst
)

// --------------------------------TYPES--------------------------------

type Elevator struct {
	Floor        int
	Direction    elevio.MotorDirection
	Behavior     ElevatorBehavior // 0:idle, 1:open, 2:moving, 3: obst
	OpenDuration float64
	CabRequests  config.OrdersCab
	Requests     config.Requests // list default as false
}

type BehaviorAndDirection struct {
	Behavior  ElevatorBehavior
	Direction elevio.MotorDirection
}

// --------------------------------FUNCTIONS--------------------------------

func ElevatorBehaviorToString(elev Elevator) string {
	switch elev.Behavior {
	case BehaviorIdle:
		return "idle"
	case BehaviorMoving:
		return "moving"
	case BehaviorOpen:
		return "open"
	case BehaviorObst:
		return "obst"
	default:
		return "undefined"
	}
}

func ElevatorDirectionToString(elev Elevator) string {
	switch elev.Direction {
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
		OpenDuration: config.DoorOpenDuration,
	}
}

func SetElevatorFloor(elev *Elevator, floor int) {
	elev.Floor = floor
}

func SetElevatorDir(elev *Elevator, dir elevio.MotorDirection) {
	elev.Direction = dir
}

func SetElevatorBehaviour(elev *Elevator, behavior ElevatorBehavior) {
	elev.Behavior = behavior
}

func SetElevatorOpenDuration(elev *Elevator, time_s float64) {
	elev.OpenDuration = time_s
}
