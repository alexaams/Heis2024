package config

import "ProjectHeis/drivers/elevio"

// --------------------------------GLOBALS--------------------------------
const NumElevators int = 3
const NumFloors int = 4
const NumButtonTypes int = 3
const DoorOpenDuration float64 = 4.0 // [s] open door duration

const BackupFile string = "SystemBackup.txt"
const BackupDir string = "BackupFiles"

var ElevatorID int = -1

var G_Ch_clear_orders = make(chan []elevio.ButtonEvent)
var G_Ch_cab_orders = make(chan []bool)
var G_Ch_hall_orders = make(chan OrdersHall)
// --------------------------------TYPES--------------------------------

type ReqList map[int]bool
type AckList [NumElevators]bool
type OrdersAckTable []AckList
type OrdersCab [NumFloors]bool

type OrdersHall [][2]bool // [floor][False]: ned [floor][True]: OPP
type Requests [NumFloors][NumButtonTypes]bool

// --------------------------------STRUCTS--------------------------------

type Order struct {
	Taken bool
	ID    int
}

// -------------------------------FUNCTIONS--------------------------------
func InitEmptyOrder() OrdersHall {
	OrdersNull := make(OrdersHall, NumFloors)
	for i := range NumFloors {
		OrdersNull[i] = [2]bool{false, false}
	}
	return OrdersNull
}
