package elevator

import (
	"ProjectHeis/config_folder/types"
	"ProjectHeis/drivers/elevio"
	"time"
)

// --------------------------------VALUES--------------------------------

// type ElevatorBehavior int

// // const (
// // 	BehaviorIdle ElevatorBehavior = iota
// // 	BehaviorMoving
// // 	BehaviorOpen
// // 	BehaviorObst
// // )

// --------------------------------TYPES--------------------------------

// USE THIS WRAPPER TO CREATE METHODS FROM TYPES
type Elevator struct {
	types.Elevator
}

// --------------------------------FUNCTIONS--------------------------------

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
	case types.BehaviorIdle:
		return "idle"
	case types.BehaviorMoving:
		return "moving"
	case types.BehaviorOpen:
		return "open"
	case types.BehaviorObst:
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

func (elev *Elevator) SetElevatorBehaviour(behavior types.ElevatorBehavior) {
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
