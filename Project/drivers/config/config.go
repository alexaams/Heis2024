package config

// --------------------------------GLOBALS--------------------------------
const NumElevators int = 3
const NumFloors int = 4
const NumButtons int = 3
const BackupFile string = "systemBackup.txt"
const DoorOpenDuration float64 = 4.0 // [s] open door duration

var ElevatorID int = -1

// --------------------------------TYPES--------------------------------

type ReqList map[int]bool
type AckList [NumElevators]bool
type OrdersAckTable []AckList
type OrdersCab [NumFloors]bool

type OrdersHall [NumFloors][2]bool // [floor][False]: ned [floor][True]: OPP
type Requests [NumFloors][NumButtons]bool

// --------------------------------STRUCTS--------------------------------

type Order struct {
	Taken bool
	ID    int
}

type GlobalOrders [NumFloors][2]Order

// -------------------------------FUNCTIONS--------------------------------
