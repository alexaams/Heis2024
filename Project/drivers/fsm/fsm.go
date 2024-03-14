package fsm

import (
	"ProjectHeis/config_folder/config"
	"ProjectHeis/config_folder/types"
	"ProjectHeis/drivers/elevator"
	"ProjectHeis/drivers/elevio"
	"ProjectHeis/requests"
	"fmt"
	"time"
)

func requestUpdates() {
	BevAndDir := requests.RequestToElevatorMovement(elevator.G_this_Elevator)
	elevio.SetMotorDirection(BevAndDir.Direction)
	elevator.G_this_Elevator.SetElevatorBehaviour(BevAndDir.Behavior)
}

func CheckFloorCurrent(a int) {
	elevator.G_this_Elevator.Floor = a
	elevio.SetFloorIndicator(elevator.G_this_Elevator.Floor)
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

	elevio.Init("localhost:15657", config.NumFloors) //Kan vi legge inn portnumber som en variabel fra config i stedet? God kodeskikk
	fmt.Print("Initiating FSM...")
	drv_floors := make(chan int)
	drv_obstr := make(chan bool)
	drv_stop := make(chan bool)
	//Initiate elevator IO (buttons are read in the event-handler)
	go elevio.PollFloorSensor(drv_floors)
	go elevio.PollObstructionSwitch(drv_obstr)
	go elevio.PollStopButton(drv_stop)
	//State-machine for elevator-behavior
	go StateMachineBehavior()

	//Ticker for frequent update of elevator-states
	var timer = time.NewTicker(300 * time.Millisecond)
	defer timer.Stop()

	for {
		select {
		case <-timer.C:
			elevator.G_Ch_elevator_update <- elevator.G_this_Elevator
			if elevator.G_this_Elevator.Behavior != types.BehaviorOpen {
				requestUpdates()
			}

		case a := <-drv_floors:
			CheckFloorCurrent(a)

		case a := <-drv_obstr:
			fmt.Printf("%+v\n", a)
			if a {
				fmt.Print("Obstruction\n")
				elevator.G_this_Elevator.SetElevatorBehaviour(types.BehaviorObst)
			} else {
				fmt.Print("Obstruction removed\n")
				elevator.G_this_Elevator.SetElevatorBehaviour(types.BehaviorIdle)
			}

		case a := <-drv_stop:
			fmt.Printf("Stopbutton: %+v\n", a)

			if a {
				elevio.SetStopLamp(true)
				elevio.SetMotorDirection(types.MD_Stop)
			} else {
				elevio.SetStopLamp(false)
			}
		case requests := <-ch_requests:
			fmt.Print(requests)
			mapNewRequests(requests)
			requestUpdates()
			CheckFloorCurrent(elevator.G_this_Elevator.Floor)
			lampChange()
		}
	}
}

func StateMachineBehavior() { //Hold the door (3 seconds)
	//Close the door
	elevator.G_door_open_counter = 0
	clearOrderFlag := true

	for {
		switch elevator.G_this_Elevator.Behavior {
		case types.BehaviorOpen:
			if clearOrderFlag {
				requests.ClearOrders(elevator.G_this_Elevator)
			}
			elevator.G_door_open_counter++
			if elevator.G_door_open_counter > elevator.G_ticks {
				elevio.SetDoorOpenLamp(false)
				elevator.G_this_Elevator.SetElevatorBehaviour(types.BehaviorIdle)
				elevator.G_door_open_counter = 0
				clearOrderFlag = false
			}
			time.Sleep(10 * time.Millisecond)
		case types.BehaviorIdle:
			//Usikker på om det er behov for noe her
			//requestUpdates()
			time.Sleep(10 * time.Millisecond)
		case types.BehaviorMoving:
			//litt usikker på hva som skal skje her egentlig, men lar den stå
			time.Sleep(10 * time.Millisecond)
		case types.BehaviorObst:
			fmt.Print("obstruction - state machine!")
		}
	}
}
