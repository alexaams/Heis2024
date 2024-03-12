package fms

import (
	"ProjectHeis/config"
	"ProjectHeis/cost"
	"ProjectHeis/drivers/elevator"
	"ProjectHeis/drivers/elevio"
	"ProjectHeis/network/peers"
	"ProjectHeis/requests"
	"ProjectHeis/ticker"
	"fmt"
	"time"
)

// channels
var sendElevDataChan = make(chan bool)
var obschan = make(chan bool)
var reqFinChan = make(chan elevio.ButtonEvent)
var updateCabChan = make(chan []bool, 1)
var orderCompleteChan = make(chan elevio.ButtonEvent)
var elevUpdateChan = make(chan elevator.Elevator)

var numFloors = config.NumFloors
var cuElevator elevator.Elevator
var peersElevator peers.PeersData
var peersUpdate peers.PeerUpdate
var peersDataMap = make(map[int]peers.PeersData)

func InitFms() {
	peersElevator = peers.InitPeers()
	fmt.Println("Starting FMS")
	eventHandling()
	peers.G_Ch_PeersData_Tx <- peersElevator
}

func requestUpdates() {
	switch cuElevator.Behavior {
	case elevator.BehaviorIdle:
		requests.WhichWay(&cuElevator)
	}
}

func FloorCurrent(a int) {
	cuElevator.Floor = a
	elevio.SetFloorIndicator(cuElevator.Floor)
	StopFlag := requests.IsThisOurStop(&cuElevator)
	switch StopFlag {
	case true:
		elevio.SetMotorDirection(elevio.MD_Stop)
		ticker.TickerStart(cuElevator.OpenDuration)
		elevio.SetDoorOpenLamp(true)
		requests.ClearOrders(cuElevator) //FEIL HER! Klarer å gjennomføre funksjonen, problemet oppstår når vi prøver å sette verdien til config.G_Ch_clear_orders
		cuElevator.Behavior = elevator.BehaviorOpen
		ticker.TickerStart(cuElevator.OpenDuration)
		elevio.SetDoorOpenLamp(false)
		cuElevator.Behavior = elevator.BehaviorIdle

	default:
	}
}

func clearRequestsPeer(variable interface{}) {
	switch types := variable.(type) {
	case elevio.ButtonEvent:
		orderCompleteChan <- types

	case []elevio.ButtonEvent:
		for _, t := range types {
			orderCompleteChan <- t
		}
	}
}

func ObstFound() {
	if cuElevator.Behavior == elevator.BehaviorOpen {
		ticker.TickerStart(cuElevator.OpenDuration)
		obschan <- true
	}
}

func StopFound(a bool) {
	//fmt.Printf("%+v\n", a)
	for f := 0; f < numFloors; f++ {
		for b := elevio.ButtonType(0); b < 3; b++ {
			elevio.SetButtonLamp(b, f, false)
		}
	}
	if a {
		elevio.SetStopLamp(true)
		elevio.SetMotorDirection(elevio.MD_Stop)
	}
}

func cabRequestAssigner(orders []bool) {
	for i, j := range orders {
		cuElevator.Requests[i][2] = j
	}
}

func hallRequestAssigner(orders config.OrdersHall) {
	for i := 0; i < config.NumFloors; i++ {
		for j := 0; j < 2; j++ {
			cuElevator.Requests[i][j] = orders[i][j]
		}
	}
}

func lampChange() {
	for floors := range config.NumFloors {
		for buttons := range config.NumButtonTypes - 1 {
			elevio.SetButtonLamp(elevio.ButtonType(buttons), floors, cuElevator.Requests[floors][buttons])
		}
		elevio.SetButtonLamp(elevio.BT_Cab, floors, cuElevator.CabRequests[floors])
	}
}

func updateOrders() {
	peersElevator.SingleOrdersHall = cost.CostFunc(peersElevator, peersDataMap, peers.G_PeersUpdate)
	config.G_Ch_hall_orders <- peersElevator.SingleOrdersHall
	peers.G_Ch_PeersData_Tx <- peersElevator
}

func newPeersData(msg peers.PeersData) bool {
	newOrder := false
	peersDataMap[msg.Id] = msg
	newOrderGlobal := make(config.OrdersHall, config.NumFloors)
	if msg.Id == peersElevator.Id {
		return newOrder
	}
	for i := range peersElevator.GlobalOrderHall {
		for j := 0; j < 2; j++ {
			if msg.GlobalOrderHall[i][j] {
				newOrderGlobal[i][j] = true
				if !peersElevator.GlobalOrderHall[i][j] {
					newOrder = true
				}
			} else {
				newOrderGlobal[i][j] = peersElevator.GlobalOrderHall[i][j]
			}
		}
	}
	peersElevator.GlobalOrderHall = newOrderGlobal
	return newOrder
}

func btnEventHandler(btnEvent elevio.ButtonEvent) {
	if btnEvent.Button == elevio.BT_Cab {
		cuElevator.CabRequests[btnEvent.Floor] = true
		config.G_Ch_cab_orders <- cuElevator.CabRequests[:]
		//Oppdatere cab-orders
	} else {
		cuElevator.Requests[btnEvent.Floor][btnEvent.Button] = true
		peersElevator.GlobalOrderHall[btnEvent.Floor][btnEvent.Button] = true
	}
}

func clearRequestHandler(btnToClear []elevio.ButtonEvent) {
	for _, y := range btnToClear {
		switch y.Button {
		case elevio.BT_Cab:
			cuElevator.CabRequests[y.Floor] = false
			cuElevator.Requests[y.Floor][y.Button] = false
		default:
			cuElevator.Requests[y.Floor][y.Button] = false
		}
	}
	config.G_Ch_cab_orders <- cuElevator.CabRequests[:]
	updateOrders() //updateOrders må endres, kan den bare ta inn peersData, altså det som skal sendes til kost?
	//updateORders trenger egentlig bare peersData som input-argument. Får vi fikset dette, så er ting litt mer isolert.
}

func orderCompleteHandler(orderComplete elevio.ButtonEvent) {
	if orderComplete.Button == elevio.BT_Cab {
		peersElevator.Elevator.CabRequests[orderComplete.Floor] = false
		//skrive til fil
	} else {
		peersElevator.SingleOrdersHall[orderComplete.Floor][orderComplete.Button] = false
		peersElevator.GlobalOrderHall[orderComplete.Floor][orderComplete.Button] = false
	}
}

func eventHandling() {
	var (
		timer = time.NewTicker(300 * time.Millisecond)
	)

	defer timer.Stop()

	go Fms()

	for {
		select {
		case <-timer.C:
			if len(peersUpdate.Lost) > 0 {
				updateOrders()
			}
		case msg := <-peers.G_Ch_PeersData_Rx:
			if newPeersData(msg) {
				updateOrders()
			}
		case btnEvent := <-elevio.G_Ch_drv_buttons:
			btnEventHandler(btnEvent)

		case elevData := <-elevUpdateChan:
			peersElevator.Elevator = elevData

		case orderComplete := <-config.G_Ch_clear_orders:
			clearRequestHandler(orderComplete)
		}
		lampChange()
	}

}

func Fms() {

	elevio.Init("localhost:15657", numFloors)
	elevio.Init_elevator_IO()

	for {
		select {
		case a := <-elevio.G_Ch_drv_floors:
			FloorCurrent(a)

		case a := <-elevio.G_Ch_drv_obstr:
			fmt.Printf("Obstruction: %+v\n", a)
			if a {
				ObstFound()
			}

		case a := <-elevio.G_Ch_stop:
			fmt.Printf("%+v\n", a)
			for f := 0; f < numFloors; f++ {
				for b := elevio.ButtonType(0); b < 3; b++ {
					elevio.SetButtonLamp(b, f, false)
				}
			}
			if a {
				elevio.SetStopLamp(true)
				elevio.SetMotorDirection(elevio.MD_Stop)
			}
		case hallorders := <-config.G_Ch_hall_orders:
			fmt.Print(hallorders)
			hallRequestAssigner(hallorders)
			requestUpdates()

		case cabOrders := <-config.G_Ch_cab_orders:
			fmt.Print(cabOrders)
			cabRequestAssigner(cabOrders)
			requestUpdates()
		}
	}
}

func OrdersFirstStep() {
	for {
		select {
		case button := <-elevio.G_Ch_drv_buttons:
			switch button.Button {
			case elevio.BT_Cab:
				cuElevator.CabRequests[button.Floor] = true
				config.G_Ch_cab_orders <- cuElevator.CabRequests[:]
				elevio.SetButtonLamp(button.Button, button.Floor, true)
				//Må også oppdatere for kostfunksjon
			}

		}
	}
}
