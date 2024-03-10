package globals

// ------------------ CONST ------------------
const (
	NumElevators   int = 3
	NumFloors      int = 4
	NumButtonTypes int = 3

	DoorOpenDuration float64 = 4.0 // [s] open door duration
	ElevatorInitID   int     = -1

	BackupFile string = "SystemBackup.txt"
	BackupDir  string = "BackupFiles"
)

// ------------------ VARIABLES ------------------

var (
	ElevatorPort int = 12345
)

// ------------------ CHANNELS ------------------
