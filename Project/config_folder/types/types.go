package types

import (
	"ProjectHeis/config_folder/config"
)

// -------------------------------- TYPES --------------------------------
type ReqList map[int]bool
type AckList [config.NumElevators]bool
type OrdersAckTable []AckList
type OrdersCab [config.NumFloors]bool
type OrdersHall [][2]bool // [floor][False]: ned [floor][True]: OPP

// -------------------------------- ENUM --------------------------------

type MotorDirection int

const (
	MD_Up   MotorDirection = 1
	MD_Down                = -1
	MD_Stop                = 0
)

type ButtonType int

const (
	BT_HallUp ButtonType = iota
	BT_HallDown
	BT_Cab
)

type ElevatorBehavior int

const (
	BehaviorIdle ElevatorBehavior = iota
	BehaviorMoving
	BehaviorOpen
	BehaviorObst
)

// -------------------------------- STRUCTS --------------------------------
type Requests struct {
	HallUp   []bool
	HallDown []bool
	CabFloor []bool
}

type ButtonEvent struct {
	Floor  int
	Button ButtonType
}

type PeerUpdate struct {
	Peers []string
	New   string
	Lost  []string
}

type Order struct {
	Active   bool
	BtnEvent ButtonEvent
}

type BehaviorAndDirection struct {
	Behavior  ElevatorBehavior
	Direction MotorDirection
}

// -------------------------------- FUNC --------------------------------

func InitEmptyOrder() OrdersHall {
	OrdersNull := make(OrdersHall, config.NumFloors)
	for i := range config.NumFloors {
		OrdersNull[i] = [2]bool{false, false}
	}
	return OrdersNull
}

func InitRequests() Requests {
	return Requests{
		HallUp:   make([]bool, config.NumFloors), // Top floor has no up button
		HallDown: make([]bool, config.NumFloors), // Ground floor has no down button
		CabFloor: make([]bool, config.NumFloors),
	}
}
