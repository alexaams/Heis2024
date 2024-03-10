package types

import (
	"ProjectHeis/config_folder/globals"
	"ProjectHeis/drivers/elevator"
	"ProjectHeis/drivers/elevio"
)

type ReqList map[int]bool
type AckList [globals.NumElevators]bool
type OrdersAckTable []AckList
type OrdersCab [globals.NumFloors]bool

type OrdersHall [][2]bool // [floor][False]: ned [floor][True]: OPP
type Requests [globals.NumFloors][globals.NumButtonTypes]bool

// -------------------------------- ENUM --------------------------------E

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

// --------------------------------STRUCTS--------------------------------
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
	Direction elevio.MotorDirection
}

type PeersData struct {
	Elevator         elevator.Elevator
	Id               int
	SingleOrdersHall OrdersHall
	GlobalOrderHall  OrdersHall
}
