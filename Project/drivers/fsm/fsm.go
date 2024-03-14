package fsm

import (
	"ProjectHeis/config_folder/config"
	"ProjectHeis/config_folder/globals"
	"ProjectHeis/config_folder/types"
	"ProjectHeis/drivers/elevator"
	"ProjectHeis/drivers/elevio"
	"ProjectHeis/network/peers"
	"ProjectHeis/requests"
	"fmt"
)

// channels
var sendElevDataChan = make(chan bool)
var obschan = make(chan bool)
var reqFinChan = make(chan types.ButtonEvent)
var updateCabChan = make(chan []bool, 1)
var orderCompleteChan = make(chan types.ButtonEvent)
var elevUpdateChan = make(chan elevator.Elevator)

// variables
var cuElevator elevator.Elevator
var peersElevator peers.PeersData
var peersUpdate peers.PeerUpdate
var peersDataMap = make(map[int]peers.PeersData)

func InitFms() {
	peersElevator = peers.InitPeers()
	fmt.Println("Starting FMS")
	peers.G_Ch_PeersData_Tx <- peersElevator
}

func requestUpdates() {
	var buttonpressed types.ButtonEvent
	switch elevator.G_this_Elevator.Behavior {
	case types.BehaviorOpen:
		fmt.Println("before if opendoor")
		if floor, buttonType := requests.ClearRequestBtnReturn(cuElevator); floor > -1 {
			fmt.Println("if was initiated")
			ticker.TickerStart(cuElevator.OpenDuration)
			buttonpressed.Button = buttonType
			buttonpressed.Floor = floor
			requests.ClearOneRequest(&cuElevator, buttonpressed)
			clearElevator := requests.RequestReadyForClear(cuElevator)
			clearRequestsPeer(clearElevator)
		}

	case elevator.BehaviorIdle:
		set := requests.RequestToElevatorMovement(cuElevator)
		cuElevator.Behavior = set.Behavior
		cuElevator.Direction = set.Direction
		fmt.Println("in idle not moving forward")
		switch set.Behavior {
		case elevator.BehaviorOpen:
			elevio.SetDoorOpenLamp(true)
			ticker.TickerStart(cuElevator.OpenDuration)
			requests.ClearOneRequest(&cuElevator, buttonpressed)
			clearElevator := requests.RequestReadyForClear(cuElevator)
			clearRequestsPeer(clearElevator)
			fmt.Println("stuck here?")

		case elevator.BehaviorMoving:
			fmt.Println("behavior moving", cuElevator.Direction)
			elevio.SetDoorOpenLamp(false)
			elevio.SetMotorDirection(cuElevator.Direction)
		}

	}
}

func FloorCurrent(a int) {
	cuElevator.Floor = a
	elevio.SetFloorIndicator(cuElevator.Floor)
	if requests.IsThisOurStop(elevator.G_this_Elevator) {
		elevator.G_this_Elevator.Stop()
		elevio.SetDoorOpenLamp(true)
		elevator.G_this_Elevator.SetElevatorBehaviour(types.BehaviorOpen)
	}
}

func mapNewRequests(reqs types.Requests) {
	for i := 0; i < config.NumFloors; i++ {
		elevator.G_this_Elevator.Requests.HallUp[i] = reqs.HallUp[i]
		elevator.G_this_Elevator.Requests.HallDown[i] = reqs.HallDown[i]
		elevator.G_this_Elevator.Requests.CabFloor[i] = reqs.CabFloor[i]

	}
}

func lampChange() {
	for floor := range config.NumFloors {
		elevio.SetButtonLamp(types.BT_Cab, floor, elevator.G_this_Elevator.Requests.CabFloor[floor])
		elevio.SetButtonLamp(types.BT_HallUp, floor, elevator.G_this_Elevator.Requests.HallUp[floor])
		elevio.SetButtonLamp(types.BT_HallDown, floor, elevator.G_this_Elevator.Requests.HallDown[floor])
	}
}

func Fsm(ch_requests chan types.Requests) {

	elevio.Init("localhost:15657", globals.NumFloors) //Kan vi legge inn portnumber som en variabel fra config i stedet? God kodeskikk

	drv_floors := make(chan int)
	drv_obstr := make(chan bool)
	drv_stop := make(chan bool)
	//Initiate elevator IO (buttons are read in the event-handler)
	go elevio.PollFloorSensor(drv_floors)
	go elevio.PollObstructionSwitch(drv_obstr)
	go elevio.PollStopButton(drv_stop)
	//State-machine for elevator-behavior
	go StateMachineBehavior()

	for {
		select {
		case a := <-drv_floors:
			FloorCurrent(a)

		case a := <-drv_obstr:
			fmt.Printf("%+v\n", a)
			if a {
				fmt.Print("Obstruction\n")
			}

		case a := <-drv_stop:
			fmt.Printf("Stopbutton: %+v\n", a)

			if a {
				elevio.SetStopLamp(true)
				elevio.SetMotorDirection(types.MD_Stop)
			}
		case requests := <-ch_requests:
			mapNewRequests(requests)
			requestUpdates()
			FloorCurrent(elevator.G_this_Elevator.Floor)
		}
	}
}

func StateMachineBehavior() {
	switch elevator.G_this_Elevator.Behavior {
	case types.BehaviorOpen:
		//Clear orders
		//lampChange()
		//Hold the door (3 seconds)
		//Close the door
		//Set door-lamp to false
		//Set behavior to idle
		//Continue
	case types.BehaviorIdle:
		//requestUpdates
		//lampChange
		//set behavior to moving if we have orders to move to, check requestUpdates()
		//time.sleep(10ms)
	case types.BehaviorMoving:
		//litt usikker på hva som skal skje her egentlig
		//time.sleep(10ms)
	case types.BehaviorObst:
		//stay here, there is an obstruction!
		//Return hall-orders
		//Inform that you have a problem
	}
}
