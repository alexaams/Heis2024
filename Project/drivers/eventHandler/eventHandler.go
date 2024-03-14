package eventHandler

import (
	"ProjectHeis/config_folder/config"
	"ProjectHeis/config_folder/types"
	"ProjectHeis/cost"
	"ProjectHeis/drivers/elevator"
	"ProjectHeis/drivers/elevio"
	"ProjectHeis/network/peers"
	"ProjectHeis/requests"
	"fmt"
	"time"
)

func EventHandling() {
	var (
		timer = time.NewTicker(300 * time.Millisecond)
	)
	fmt.Print("Eventhandler starting...\n")
	defer timer.Stop()

	go elevio.PollButtons(elevator.G_Ch_drv_buttons)

	for {
		select {
		case <-timer.C:
			if len(peers.G_PeersUpdate.Lost) > 0 {
				updateOrders(peers.G_PeersElevator) //må lage logikk
			}
		case msg := <-peers.G_Ch_PeersData_Rx:
			removeAcknowledgedOrder(msg)
			if newPeersData(msg) {
				updateOrders(msg)
			}
		case btnEvent := <-elevator.G_Ch_drv_buttons:
			btnEventHandler(btnEvent)

		case elevData := <-elevator.G_Ch_elevator_update:
			peers.G_PeersElevator.Elevator = elevData

		case orderComplete := <-elevator.G_Ch_clear_orders:
			orderCompleteHandler(orderComplete)
		}
		peers.G_Ch_PeersData_Tx <- peers.G_PeersElevator
	}
}

func updateOrders(someElevator peers.PeersData) {
	if peers.G_isMaster {
		someElevator.SingleOrdersHall = cost.CostFunc(someElevator)
		peers.G_Ch_PeersData_Tx <- someElevator
		fmt.Println("runs cost as master")
	}
	if someElevator.Id == peers.G_PeersElevator.Id {
		orderToRequest := OrdersHallToRequest(peers.G_PeersElevator.SingleOrdersHall)
		elevator.G_Ch_requests <- orderToRequest
		fmt.Println("sent request to fsm")
	}
}

func OrdersHallToRequest(order types.OrdersHall) types.Requests {
	req := types.InitRequests()
	for i, ord := range order {
		req.HallUp[i] = ord[0]
		req.HallDown[i] = ord[1]
	}
	return req
}

func removeAcknowledgedOrder(msg peers.PeersData) {
	for i := range msg.GlobalAckOrders {
		for j := 0; j < 2; j++ {
			if msg.GlobalAckOrders[i][j] == peers.G_PeersElevator.GlobalOrderHall[i][j] {
				peers.G_PeersElevator.GlobalOrderHall[i][j] = false
			}
		}
	}
}

func newPeersData(msg peers.PeersData) bool {
	newOrder := false
	peers.G_Datamap[msg.Id] = msg
	newOrderGlobal := make(types.OrdersHall, config.NumFloors)
	if msg.Id == peers.G_PeersElevator.Id {
		return newOrder
	}
	for i := range peers.G_PeersElevator.GlobalOrderHall {
		for j := 0; j < 2; j++ {
			if msg.GlobalOrderHall[i][j] {
				newOrderGlobal[i][j] = true
				if !peers.G_PeersElevator.GlobalOrderHall[i][j] {
					newOrder = true
				}
			} else {
				newOrderGlobal[i][j] = peers.G_PeersElevator.GlobalOrderHall[i][j]
			}
		}
	}
	peers.G_PeersElevator.GlobalOrderHall = newOrderGlobal
	return newOrder
}

func btnEventHandler(btnEvent types.ButtonEvent) {
	if btnEvent.Button == types.BT_Cab {
		peers.G_PeersElevator.Elevator.Requests.CabFloor[btnEvent.Floor] = true
		elevator.G_Ch_requests <- peers.G_PeersElevator.Elevator.Requests
	} else {
		peers.G_PeersElevator.GlobalOrderHall[btnEvent.Floor][btnEvent.Button] = true
		peers.G_PeersElevator.SingleOrdersHall[btnEvent.Floor][btnEvent.Button] = true
		updateOrders(peers.G_PeersElevator)
	}
}

func orderCompleteHandler(orderComplete []types.ButtonEvent) {
	requests.ClearRequests(&peers.G_PeersElevator.Elevator, orderComplete)
	requests.ClearRequests(&elevator.G_this_Elevator, orderComplete)
	for _, order := range orderComplete {
		if order.Button == types.BT_Cab {
			peers.G_PeersElevator.Elevator.Requests.CabFloor[order.Floor] = false
			//skrive til fil
		} else {
			peers.G_PeersElevator.SingleOrdersHall[order.Floor][order.Button] = false
			peers.G_PeersElevator.GlobalOrderHall[order.Floor][order.Button] = false
			peers.G_PeersElevator.GlobalAckOrders[order.Floor][order.Button] = true
			peers.G_Ch_PeersData_Tx <- peers.G_PeersElevator
			peers.G_PeersElevator.GlobalAckOrders[order.Floor][order.Button] = false
		}
	}
}
