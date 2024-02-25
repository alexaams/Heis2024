package config

import (
	"ProjectHeis/drivers/elevio"
	"ProjectHeis/network/localip"
	"flag"
	"fmt"
	"os"
)

// ---- Globals----
const NumElevators int = 3
const NumFloors int = 4
const NumButtons int = 3


var ElevatorID int = -1



// ---------TYPES-----------
type ReqList 			map[int]bool
type AckList 			[NumElevators]bool
type OrderList			[NumFloors]bool


// ---------STRUCTS----------
type Elevator struct {
	Floor 				int
	Direction 			elevio.MotorDirection
	Requests			[NumFloors][NumElevators]bool
	Status				int 	// 0:idle, 1:open, 2:moving, 3: obst

}

type PeersConnection struct{
	Peers				[]string
	New					[]string
	Lost				[]string
}

type PeersData struct{
	Elevator 			Elevator
	Id					int
	//Orders
	//
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

//Creating ID with local ip and PID
func CreateID () string {
	id := ""
	flag.StringVar(&id, "id", "", "id of this peer")
	flag.Parse()

	if id == "" {
		localIP, err := localip.LocalIP()
		if err != nil {
			fmt.Println(err)
			localIP = "DISCONNECTED"
		}
		id = fmt.Sprintf("peer-%s-%d", localIP, os.Getpid())
	}

	return id
}