package elevator

import (
	"ProjectHeis/Scratch/config_folder/globals"
	"ProjectHeis/Scratch/config_folder/types"
	"ProjectHeis/Scratch/drivers/elevio"
	"time"
)

// --------------------------------VALUES--------------------------------

// --------------------------------TYPES--------------------------------

type Elevator struct {
	Floor        int
	Direction    types.MotorDirection
	Behavior     types.ElevatorBehavior // 0:Idle, 1:Moving, 2:Open, 3: Obst
	OpenDuration int
	Requests     types.Requests // list default as false
}

// --------------------------------FUNCTIONS--------------------------------

func InitElevator() Elevator {
	return Elevator{
		Floor:        -1,
		Direction:    types.MD_Stop,
		Behavior:     types.BehaviorIdle,
		OpenDuration: globals.DoorOpenDuration,
	}
}

func (elev *Elevator) HasCabRequests() bool {
	for _, hasRequest := range elev.Requests.CabFloor {
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
	case types.MD_Stop:
		return "stop"
	case types.MD_Up:
		return "up"
	case types.MD_Down:
		return "down"
	default:
		return "undefined"
	}
}

func (elev *Elevator) MoveUp() {
	elevio.SetMotorDirection(types.MD_Up)
}

func (elev *Elevator) MoveDown() {
	elevio.SetMotorDirection(types.MD_Down)
}

func (elev *Elevator) Stop() {
	elevio.SetMotorDirection(types.MD_Stop)
}

func (elev *Elevator) SetElevatorFloor(floor int) {
	elev.Floor = floor
}

func (elev *Elevator) SetElevatorDir(dir types.MotorDirection) {
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