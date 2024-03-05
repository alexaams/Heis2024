package cost

import (
	"ProjectHeis/drivers/config"
	"ProjectHeis/drivers/elevator"
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
	HallRequests config.OrdersHall       `json:"hallRequests"`
	States       map[string]HRAElevState `json:"states"`
}

func OrderEmpty(order config.OrdersHall) bool {
	for i := 0; i < config.NumFloors; i++ {
		for j := 0; j < config.NumElevators; j++ {
			if order[i][j] {
				return false
			}
		}
	}
	return true

}

func CostFunc(elevatorObject config.PeersData, hallRequests config.OrdersHall, peers config.PeersConnection) config.OrdersHall {
	if OrderEmpty(elevatorObject.OrdersHall) {
		fmt.Println("No orders available in hall request")
		return elevatorObject.OrdersHall
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

	peersActive := len(peers.Peers)
	statesElevators := make(map[string]HRAElevState, peersActive)
	idstring := strconv.Itoa(elevatorObject.Id)

	for range statesElevators {
		id := elevatorObject.Id
		statesElevators[strconv.Itoa(id)] = elevatorToHRAState(elevatorObject.Elevator)
	}

	input := HRAInput{
		HallRequests: elevatorObject.OrdersHall,
		States:       statesElevators,
	}

	jsonBytes, err := json.Marshal(input)
	if err != nil {
		fmt.Println("json.Marshal error: ", err)
	}

	ret, err := exec.Command("../hall_request_assigner/"+hraExecutable, "-i", string(jsonBytes)).CombinedOutput()
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
	return config.OrdersHall(ordersFixed)
}

func elevatorToHRAState(elev elevator.Elevator) HRAElevState {
	return HRAElevState{
		Behavior:    elevator.ElevatorBehaviorToString(elev),
		Floor:       elev.Floor,
		Direction:   elevator.ElevatorDirectionToString(elev),
		CabRequests: elev.CabRequests[:],
	}

}
