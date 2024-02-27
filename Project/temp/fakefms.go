package temp

import (
	"ProjectHeis/drivers/elevio"
	"fmt"
)

// ---------TYPES-----------
type ReqList map[int]bool

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

func ElevMoving(ReqFloor int, CurrentFloor int) {
	if ReqFloor == CurrentFloor {
		elevio.SetMotorDirection(elevio.MD_Stop)
		elevio.SetDoorOpenLamp(true)
		for i := elevio.ButtonType(0); i < 3; i++ {
			elevio.SetButtonLamp(i, CurrentFloor, false)
			fmt.Println("Button val", i)
		}
	} else if ReqFloor-CurrentFloor >= 1 {
		elevio.SetDoorOpenLamp(false)
		elevio.SetMotorDirection(elevio.MD_Up)
	} else {
		elevio.SetDoorOpenLamp(false)
		elevio.SetMotorDirection(elevio.MD_Down)
	}
	fmt.Println("Current:  ", CurrentFloor)

}
