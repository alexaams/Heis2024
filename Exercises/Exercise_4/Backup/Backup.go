package backup

import (
	"fmt"
	"os"
	"time"
)

func terminate() {
	time.Sleep(3 * time.Second)
	defer fmt.Println("Program terminated")
	os.Exit(3)
}

func main(){
	terminate()
}
