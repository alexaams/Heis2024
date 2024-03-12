package elevator

import (
	"ProjectHeis/config"
	"ProjectHeis/drivers/elevio"
	"time"
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

// USE THIS WRAPPER TO CREATE METHODS FROM TYPES

// type Elevator struct{
// 	types.Elevator
// }

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

func (elev *Elevator) ElevatorBehaviorToString() string {
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

func (elev *Elevator) ElevatorDirectionToString() string {
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

func (elev *Elevator) MoveUp() {
	elevio.SetMotorDirection(elevio.MD_Up)
}

func (elev *Elevator) MoveDown() {
	elevio.SetMotorDirection(elevio.MD_Down)
}

func (elev *Elevator) Stop() {
	elevio.SetMotorDirection(elevio.MD_Stop)
}

func (elev *Elevator) SetElevatorFloor(floor int) {
	elev.Floor = floor
}

func (elev *Elevator) SetElevatorDir(dir elevio.MotorDirection) {
	elev.Direction = dir
}

func (elev *Elevator) SetElevatorBehaviour(behavior ElevatorBehavior) {
	elev.Behavior = behavior
}

func (elev *Elevator) OpenDoor(doorOpenCh, obstCh chan bool) {
	timer := time.NewTicker(500 * time.Millisecond)
	timerCounter := 0
	for {
		select {
		case obst := <-obstCh:
			if obst {

				timer.Stop()
			} else {
				timer.Reset(1 * time.Second)
			}
		case <-timer.C:
			timerCounter++
			if timerCounter <= 6 {
				elevio.SetDoorOpenLamp(false)
				doorOpenCh <- true
				return
			}
		}
	}
}
