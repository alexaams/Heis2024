package cost

import (
	"ProjectHeis/config_folder/config"
	"ProjectHeis/config_folder/types"
	"ProjectHeis/drivers/elevator"
	"ProjectHeis/network/peers"
	"encoding/json"
	"fmt"
	"os/exec"
	"runtime"
	"strconv"
)

// Struct members must be public in order to be accessible by json.Marshal/.Unmarshal
// This means they must start with a capital letter, so we need to use field renaming struct tags to make them camelCase

type HRAElevState struct {
	Behavior    string `json:"behaviour"`
	Floor       int    `json:"floor"`
	Direction   string `json:"direction"`
	CabRequests []bool `json:"cabRequests"`
}

type HRAInput struct {
	HallRequests types.OrdersHall        `json:"hallRequests"`
	States       map[string]HRAElevState `json:"states"`
}

func OrderEmpty(order types.OrdersHall) bool {
	for i := 0; i < globals.NumFloors; i++ {
		for j := 0; j < 2; j++ {
			if order[i][j] {
				return false
			}
		}
	}
	return true

}

func CostFunc(elevatorObject peers.PeersData) types.OrdersHall {
	if OrderEmpty(elevatorObject.GlobalOrderHall) {
		fmt.Println("No orders available in hall request")
		return elevatorObject.GlobalOrderHall
	}
	hraExecutable := ""
	switch runtime.GOOS {
	case "linux":
		hraExecutable = "hall_request_assigner"
	case "windows":
		hraExecutable = "hall_request_assigner.exe"
	default:
		panic("OS not supported")
	}

	peersActive := len(peers.G_PeersUpdate.Peers)
	statesElevators := make(map[string]HRAElevState, peersActive)
	idstring := strconv.Itoa(globals.ElevatorID)
	peers.G_Datamap[globals.ElevatorID] = elevatorObject

	//Mapping all elevators to the algorithm
	for i := 0; i < peersActive; i++ {
		id, _ := strconv.Atoi(peers.G_PeersUpdate.Peers[i])
		data := peers.G_Datamap[id]
		statesElevators[strconv.Itoa(id)] = elevatorToHRAState(data.Elevator)
	}

	input := HRAInput{
		HallRequests: elevatorObject.GlobalOrderHall, //Dette skal være en globalt gjeldende liste, så vi må få på plass funksjonalitet for å sikre at denne er oppdatert!
		States:       statesElevators,
	}

	jsonBytes, err := json.Marshal(input)
	if err != nil {
		fmt.Println("json.Marshal error: ", err)
	}

	ret, err := exec.Command(hraExecutable, "-i", string(jsonBytes)).CombinedOutput()
	if err != nil {
		fmt.Println("exec.Command error: ", err)
		fmt.Println(string(ret))
	}

	output := new(map[string][][2]bool)
	err = json.Unmarshal(ret, &output)
	if err != nil {
		fmt.Println("json.Unmarshal error: ", err)
	}

	fmt.Printf("output: \n")
	for k, v := range *output {
		fmt.Printf("%6v :  %+v\n", k, v)
	}
	ordersFixed := (*output)[idstring]
	return ordersFixed
}

func elevatorToHRAState(elev elevator.Elevator) HRAElevState {
	return HRAElevState{
		Behavior:    elev.ElevatorBehaviorToString(),
		Floor:       elev.Floor,
		Direction:   elev.ElevatorDirectionToString(),
		CabRequests: elev.Requests.CabFloor,
	}
}
