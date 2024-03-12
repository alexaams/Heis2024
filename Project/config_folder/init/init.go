package init

import (
	"ProjectHeis/config_folder/globals"
	"ProjectHeis/config_folder/types"
	"ProjectHeis/drivers/elevator"
	"ProjectHeis/drivers/elevio"
	"ProjectHeis/network/localip"
)

func InitEmptyOrder() types.OrdersHall {
	OrdersNull := make(types.OrdersHall, globals.NumFloors)
	for i := range globals.NumFloors {
		OrdersNull[i] = [2]bool{false, false}
	}
	return OrdersNull
}

func InitPeers() types.PeersData {
	return types.PeersData{
		Elevator:         InitElevator(),
		Id:               localip.CreateID(),
		SingleOrdersHall: InitEmptyOrder(),
		GlobalOrderHall:  InitEmptyOrder(),
	}
}

func InitElevator() elevator.Elevator {
	return elevator.Elevator{
		Floor:        -1,
		Direction:    types.MD_Stop,
		Behavior:     types.BehaviorIdle,
		OpenDuration: globals.DoorOpenDuration,
	}
}
