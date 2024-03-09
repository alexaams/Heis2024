package config

// --------------------------------GLOBALS--------------------------------
const NumElevators int = 3
const NumFloors int = 4
const NumButtonTypes int = 3
const DoorOpenDuration float64 = 4.0 // [s] open door duration

const BackupFile string = "SystemBackup.txt"
const BackupDir string = "BackupFiles"

var ElevatorID int = -1

// --------------------------------TYPES--------------------------------

type ReqList map[int]bool
type AckList [NumElevators]bool
type OrdersAckTable []AckList
type OrdersCab [NumFloors]bool

type  OrdersHall [][2]bool // [floor][False]: ned [floor][True]: OPP
type Requests [NumFloors][NumButtonTypes]bool

// --------------------------------STRUCTS--------------------------------

type Order struct {
	Taken bool
	ID    int
}

type GlobalOrders [NumFloors][2]Order

// -------------------------------FUNCTIONS--------------------------------
