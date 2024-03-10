package cost2

import (
	"ProjectHeis/drivers/elevator"
)

type Button int
type OnClearedRequestFunc func(btn Button, floor int)

func RequestClearAtCurrentFloor(OldElev elevator.Elevator) {}
