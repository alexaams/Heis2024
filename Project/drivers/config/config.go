package config

import (
	"ProjectHeis/drivers/elevio"
	"ProjectHeis/network/localip"
	"fmt"
	"strings"
)

// ---- GLOBALS----
const NumElevators int = 3
const NumFloors int = 4
const NumButtons int = 3
const BackupFile string = "systemBackup.txt"
const doorOpenDuration float32 = 4.0 // [s] open door duration

const (
	BehaviorIdle = iota
	BehaviorMoving
	BehaviorOpen
	BehavoirObst
)

var ElevatorID int = -1

// ---------TYPES-----------
type ReqList map[int]bool
type AckList [NumElevators]bool
type OrdersAckTable []AckList
type OrdersCab [NumElevators][NumFloors]bool
type OrdersHall [NumFloors][2]bool // [floor][0]: ned [floor][1]: OPP

// ---------STRUCTS----------
type Elevator struct {
	Floor        int
	Direction    elevio.MotorDirection
	Requests     [NumFloors][NumElevators]bool
	Behavior     int // 0:idle, 1:open, 2:moving, 3: obst
	OpenDuration float32
}

type PeersConnection struct {
	Peers []string
	New   []string
	Lost  []string
}

type PeersData struct {
	Elevator       Elevator
	Id             int
	OrdersCab      []bool
	OrdersHall     OrdersHall
	GlobalAckTable OrdersAckTable
}

type HallEvent struct {
	Floor     int
	Direction int
	Id        int
}

// -------FUNCTIONS--------
func MakeReqList(amountFloors, botFloor int) ReqList {
	listFloor := make(map[int]bool)
	for x := 0; x < amountFloors; x++ {
		listFloor[x+botFloor] = false
	}
	return listFloor
}

func (r ReqList) SetFloor(floor int) {
	if _, ok := r[floor]; ok {
		r[floor] = true
	} else {
		fmt.Println("Floor does not exist")
	}
}

func (r ReqList) ClearFloor(floor int) {
	if _, ok := r[floor]; ok {
		r[floor] = false
	} else {
		fmt.Println("Floor does not exist")
	}
}

// Creating ID with local ip and PID
func CreateID() string {
	id := ""

	if id == "" {
		localIP, err := localip.LocalIP()
		if err != nil {
			fmt.Println(err)
			localIP = "DISCONNECTED"
		}
		id = localIP
		temp_arr := strings.Split(id,".")
		id = temp_arr[3]
	}

	return id
}

func ElevatorBehaviorToString(elev Elevator) string {
	behavior := elev.Behavior
	switch behavior {
	case 0:
		return "idle"
	case 1:
		return "moving"
	case 2:
		return "obst"
	default:
		return "undefined"
	}

}

func NewElevator() Elevator {
	return Elevator{
		Direction:    elevio.MD_Stop,
		Floor:        -1,
		Behavior:     BehaviorIdle,
		OpenDuration: doorOpenDuration,
	}
}
