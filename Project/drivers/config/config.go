package config

import (
	"ProjectHeis/drivers/elevator"
	"ProjectHeis/network/localip"
	"fmt"
	"strconv"
	"strings"
)

// --------------------------------GLOBALS--------------------------------
const NumElevators int = 3
const NumFloors int = 4
const NumButtons int = 3
const BackupFile string = "systemBackup.txt"
const doorOpenDuration float64 = 4.0 // [s] open door duration

var ElevatorID int = -1
var Peers PeersConnection

// --------------------------------TYPES--------------------------------

type ReqList map[int]bool
type AckList [NumElevators]bool
type OrdersAckTable []AckList
type OrdersCab [NumFloors]bool

type OrdersHall [NumFloors][2]bool // [floor][False]: ned [floor][True]: OPP
type Requests [NumFloors][NumButtons]bool

// --------------------------------STRUCTS--------------------------------

type PeersConnection struct {
	Peers []string
	New   []string
	Lost  []string
}

type PeersData struct {
	Elevator   elevator.Elevator
	Id         int
	OrdersHall OrdersHall
}

type Order struct {
	Taken bool
	ID    int
}

type GlobalOrders [NumFloors][2]Order

// -------------------------------FUNCTIONS--------------------------------

// Creating ID with local ip and PID
func CreateID() int {
	idStr := ""

	if idStr == "" {
		localIP, err := localip.LocalIP()
		if err != nil {
			fmt.Println(err)
			localIP = "DISCONNECTED"
		}
		idStr = localIP
		temp_arr := strings.Split(idStr, ".")
		idStr = temp_arr[3]

	}
	idInt, err := strconv.Atoi(idStr)
	if err != nil {
		fmt.Println("Error Converting stringID to Int: " + e)
	}
	return idInt
}

func InitPeers() PeersData {
	return PeersData{
		Elevator:   elevator.InitElevator(),
		Id:         CreateID(),
		OrdersHall: OrdersHall{},
	}
}
