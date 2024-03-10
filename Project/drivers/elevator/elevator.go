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
	Behavior     ElevatorBehavior // 0:Idle, 1:Moving, 2:Open, 3: Obst
	OpenDuration float64
	CabRequests  config.OrdersCab
	Requests     config.Requests // list default as false
}

type BehaviorAndDirection struct {
	Behavior  ElevatorBehavior
	Direction elevio.MotorDirection
}

// --------------------------------FUNCTIONS--------------------------------

func InitElevator() Elevator {
	return Elevator{
		Direction:    elevio.MD_Stop,
		Floor:        0,
		Behavior:     BehaviorIdle,
		OpenDuration: 3.0,
	}
}

func (elev *Elevator) HasCabRequests() bool {
	for _, hasRequest := range elev.CabRequests {
		if hasRequest {
			return true
		}
	}
	return false
}

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

func SetElevatorFloor(elev *Elevator, floor int) {
	elev.Floor = floor
}

func SetElevatorDir(elev *Elevator, dir elevio.MotorDirection) {
	elev.Direction = dir
}

func SetElevatorBehaviour(elev *Elevator, behavior ElevatorBehavior) {
	elev.Behavior = behavior
}

//func SetElevatorOpenDuration(elev *Elevator, time_s float64) {
//elev.OpenDuration = time_s
//}
