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
	"strconv"
	"time"
)

func EventHandling() {
	var (
		timer = time.NewTicker(1000 * time.Millisecond)
	)
	fmt.Print("Eventhandler starting...\n")
	defer timer.Stop()

	go elevio.PollButtons(elevator.G_Ch_drv_buttons)

	for {
		select {
		case <-timer.C:
			lampChange()
			if len(peers.G_PeersUpdate.Lost) > 0 {
				lost, _ := strconv.Atoi(peers.G_PeersUpdate.Lost[0])
				delete(peers.G_Datamap, lost)
				peers.G_PeersUpdate.Lost = peers.G_PeersUpdate.Lost[:0]
				updateOrders()
			}
			peers.G_Ch_PeersData_Tx <- peers.G_PeersElevator
		case msg := <-peers.G_Ch_PeersData_Rx:
			removeAcknowledgedOrder(msg)
			if newPeersData(msg) {
				updateOrders()
			}
			peers.G_Ch_PeersData_Tx <- peers.G_PeersElevator
		case btnEvent := <-elevator.G_Ch_drv_buttons:
			btnEventHandler(btnEvent)

		case elevData := <-elevator.G_Ch_elevator_update:
			peers.G_PeersElevator.Elevator = elevData
			lampChange()

		case orderComplete := <-elevator.G_Ch_clear_orders:
			orderCompleteHandler(orderComplete)
		}
	}
}

func updateOrders() {
	select {
	case peers.G_PeersElevator.SingleOrdersHall = <-cost.CostFuncChan(peers.G_PeersElevator):
	case <-time.After(30 * time.Millisecond):
		fmt.Println("Cost timeout")
		return
	}
	orderToRequest := OrdersHallToRequest(peers.G_PeersElevator.SingleOrdersHall)
	select {
	case elevator.G_Ch_requests <- orderToRequest:
	default:
		fmt.Print("Channel requests full")
	}
	select {
	case peers.G_Ch_PeersData_Tx <- peers.G_PeersElevator:
	default:
		fmt.Println("Channel transmit filled")
	}
}

func OrdersHallToRequest(order types.OrdersHall) types.Requests {
	req := types.InitRequests()
	for i, ord := range order {
		req.HallUp[i] = ord[0]
		req.HallDown[i] = ord[1]
	}
	req.CabFloor = peers.G_PeersElevator.Elevator.Requests.CabFloor
	return req
}

func removeAcknowledgedOrder(msg peers.PeersData) {
	for i := range msg.GlobalAckOrders {
		for j := 0; j < 2; j++ {
			if msg.GlobalAckOrders[i][j] == peers.G_PeersElevator.GlobalOrderHall[i][j] {
				peers.G_PeersElevator.GlobalOrderHall[i][j] = false
				peers.G_PeersElevator.GlobalAckOrders[i][j] = false
			}
		}
	}
}

func newPeersData(msg peers.PeersData) bool {
	newOrder := false
	peers.G_Datamap[msg.ElevatorId] = msg
	newOrderGlobal := make(types.OrdersHall, config.NumFloors)
	if msg.ElevatorId == peers.G_PeersElevator.ElevatorId {
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
		updateOrders()
	}
}

func orderCompleteHandler(orderComplete []types.ButtonEvent) {
	requests.ClearRequests(&peers.G_PeersElevator.Elevator, orderComplete)
	requests.ClearRequests(&elevator.G_this_Elevator, orderComplete)
	for _, order := range orderComplete {
		if order.Button == types.BT_Cab {
			peers.G_PeersElevator.Elevator.Requests.CabFloor[order.Floor] = false
		} else {
			peers.G_PeersElevator.SingleOrdersHall[order.Floor][order.Button] = false
			peers.G_PeersElevator.GlobalOrderHall[order.Floor][order.Button] = false
			peers.G_PeersElevator.GlobalAckOrders[order.Floor][order.Button] = true
			peers.G_Ch_PeersData_Tx <- peers.G_PeersElevator
		}
	}
}

func lampChange() {
	for floor := range config.NumFloors {
		elevio.SetButtonLamp(types.BT_HallUp, floor, peers.G_PeersElevator.GlobalOrderHall[floor][types.BT_HallUp])
		elevio.SetButtonLamp(types.BT_HallDown, floor, peers.G_PeersElevator.GlobalOrderHall[floor][types.BT_HallDown])
		elevio.SetButtonLamp(types.BT_Cab, floor, elevator.G_this_Elevator.Requests.CabFloor[floor])
	}
}
