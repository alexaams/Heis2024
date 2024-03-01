package requests

import (
	"ProjectHeis/drivers/config"
	"ProjectHeis/drivers/elevator"
	"fmt"
)

func ArrivedRequest(elev elevator.Elevator) bool {

}

func MakeReqList(amountFloors, botFloor int) config.ReqList {
	listFloor := make(map[int]bool)
	for x := 0; x < amountFloors; x++ {
		listFloor[x+botFloor] = false
	}
	return listFloor
}

func (r config.ReqList) SetFloor(floor int) {
	if _, ok := r[floor]; ok {
		r[floor] = true
	} else {
		fmt.Println("Floor does not exist")
	}
}

func (r config.ReqList) ClearFloor(floor int) {
	if _, ok := r[floor]; ok {
		r[floor] = false
	} else {
		fmt.Println("Floor does not exist")
	}
}
