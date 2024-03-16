package elevator

import (
	"ProjectHeis/config_folder/config"
	"ProjectHeis/config_folder/types"
	"ProjectHeis/drivers/elevio"
)

// --------------------------------Globals--------------------------------
var (
	G_Ch_clear_orders    = make(chan []types.ButtonEvent)
	G_Ch_drv_buttons     = make(chan types.ButtonEvent)
	G_Ch_drv_floors      = make(chan int)
	G_Ch_drv_obstr       = make(chan bool)
	G_Ch_stop            = make(chan bool)
	G_Ch_requests        = make(chan types.Requests, 1024)
	G_Ch_elevator_update = make(chan Elevator, 1024)
	G_this_Elevator      = InitElevator()
	G_ticks              = config.DoorOpenDuration * 100
	G_door_open_counter  = 0
)

// --------------------------------VALUES--------------------------------

// --------------------------------TYPES------------------

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
		OpenDuration: config.DoorOpenDuration,
		Requests:     types.InitRequests(),
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
		return "doorOpen"
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
