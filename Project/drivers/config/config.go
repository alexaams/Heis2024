package config

import (
	"ProjectHeis/drivers/elevio"
	"ProjectHeis/network/localip"
	"fmt"
	"strings"
)

// -------------------------------- GLOBALS--------------------------------
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

const (
	DirectionStop = iota
	DirectionUp
	DirectionDown
)

var ElevatorID int = -1

// --------------------------------TYPES--------------------------------

type ReqList map[int]bool
type AckList [NumElevators]bool
type OrdersAckTable []AckList
type OrdersCab [NumFloors]bool
type OrdersHall [NumFloors][2]bool // [floor][False]: ned [floor][True]: OPP

// --------------------------------STRUCTS--------------------------------

// ---------STRUCTS----------
// Elevator-type is current states (not how many floors it has etc.)
type Elevator struct {
	Floor        int
	Direction    elevio.MotorDirection
	Requests     OrdersCab
	Behavior     int // 0:idle, 1:open, 2:moving, 3: obst
	OpenDuration float32
}

type PeersConnection struct {
	Peers []string
	New   []string
	Lost  []string
}

type PeersData struct {
	Elevator Elevator
	Id       int
	//OrdersCab      OrdersCab
	OrdersHall     OrdersHall
	GlobalAckTable OrdersAckTable
}

type HallEvent struct {
	Floor     int
	Direction int
	Id        int
}

// True if active request, int represent elevator-ID that took the order
type RequestActive struct {
	Active     bool //True if there is an active order
	ElevatorID int  //-1 if no elevators are on it
}

// Two-dimensional array (matrix) containing all hall-request for all floors - up and down
type GlobalOrderTable [NumFloors][2]RequestActive

// -------------------------------FUNCTIONS--------------------------------

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
		temp_arr := strings.Split(id, ".")
		id = temp_arr[3]
	}

	return id
}

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
	case DirectionStop:
		return "stop"
	case DirectionUp:
		return "up"
	case DirectionDown:
		return "down"
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

func CreatePeersData(ID int) PeersData {
	return PeersData{
		Elevator:       NewElevator(),
		Id:             ID,
		OrdersHall:     OrdersHall{},
		GlobalAckTable: make(OrdersAckTable, NumElevators),
	}
}

func CreateGlobalOrderTable() GlobalOrderTable {
	var table GlobalOrderTable
	for x := 0; x < NumFloors; x++ {
		table[x][0].Active = false
		table[x][0].ElevatorID = -1
		table[x][1].Active = false
		table[x][1].ElevatorID = -1
	}
	return table
}

func (r GlobalOrderTable) PrintGlobalOrderTable() {
	for x := 0; x < NumFloors; x++ {
		fmt.Printf("%d:\n", x+1)
		fmt.Printf("Up: %t , ID: %d\n", r[x][0].Active, r[x][0].ElevatorID)
		fmt.Printf("Down: %t , ID: %d\n", r[x][1].Active, r[x][1].ElevatorID)
	}
}
