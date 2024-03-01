package elevator

import "fmt"
import "ProjectHeis/drivers/config"

// --------------------------------VALUES--------------------------------

type ElevBehavior int

const (
	BehaviorIdle ElevBehavior = iota
	BehaviorMoving
	BehaviorOpen
	BehavoirObst
)

// --------------------------------TYPES--------------------------------

type Elevator struct {
	Floor        int
	Direction    elevio.MotorDirection
	Requests     config.OrdersCab
	Behavior     ElevBehavior // 0:idle, 1:open, 2:moving, 3: obst
	OpenDuration float64
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
		return "moving"
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

func SetElevatorOpenDuration(elev *Elevator, time_s float) {
	elev.OpenDuration = time_s
}