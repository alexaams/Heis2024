package main

import (
	"fmt"
	"net"
	"os/exec"
	"runtime"
	"strconv"
	"time"
)

const address = "localhost:8080"

var ournumber int = 0

func main() {

	go backup_sequence()
	go other_func()

	select {}
}

func other_func() {
	variable := 1
	for {
		variable++
		fmt.Printf("Backgroundvariable is now: %d\n", variable)
		time.Sleep(time.Second * 3)
	}
}
func backup_sequence() {
	//Creating a UDP-listener socket (not sure about the terminology)
	conn := listenToUDP(address)

	//Initiating backup-phase with the listener-socket, data and specified deadline
	backupPhase(&ournumber, 3, conn)

	conn.Close()

	//Start a new backup-process in a new terminal when no other is detected on UDP-broadcast
	openTerminal()

	//Creating a UDP-sender socket (not sure about the terminology)
	conn_send := sendToUDP(address)

	//Initiating primary-phase with the sender-socket, data and specified deadline
	primaryPhase(&ournumber, 3, conn_send)
}

func primaryPhase(number *int, deadlineSecs int, conn *net.UDPConn) {
	fmt.Println("---Initiating primary process---")
	for {
		_, err := conn.Write([]byte(fmt.Sprint(*number)))
		fmt.Println(*number)

		if err != nil {
			fmt.Println(err)
		}

		*number++
		time.Sleep(time.Second * 2)
	}
}

func backupPhase(number *int, deadlineSecs int, conn *net.UDPConn) {
	//Initiating backup-phase
	fmt.Println("---Initiating backup process---")
	buffer := make([]byte, 1024) //Creating buffer to receive data
	for {

		deadline := time.Now().Add(time.Duration(deadlineSecs) * time.Second)
		conn.SetReadDeadline(deadline)
		n, addr, err := conn.ReadFromUDP(buffer)

		//Listen for message until deadline is reached
		if err != nil {
			fmt.Println("--Primary phase--")
			break
		}

		number_temp, err := strconv.Atoi(string(buffer[:n]))
		if err != nil {
			fmt.Println("Feil ved konvertering av streng til heltall:", err)
			return
		}
		*number = number_temp
		fmt.Printf("Mottatt %d bytes fra %s: %d\n", n, addr, (*number))
	}
}

func listenToUDP(address string) *net.UDPConn {
	//Connect to UDP-socket
	udpAddr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		fmt.Println("Error resolving UDP address: ", err)
		return nil
	}

	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		fmt.Println("Error listening on UDP port:", err)
		return nil
	}
	return conn
}

func sendToUDP(address string) *net.UDPConn {
	udpAddr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		fmt.Println("Error resolving UDP address: ", err)
		return nil
	}

	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		fmt.Println("Error opening UDP address: ", err)
		return nil
	}
	return conn
}

func openTerminal() {
	//Declaring pointer to Cmd
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		//Executing terminal command which runs main.go in this folder
		cmd = exec.Command("cmd.exe", "/c", "start", "cmd.exe", "/k", "go", "run", "main.go")
	default:
		//''Catching'' error
		fmt.Println("Dette operativsystemet støttes ikke")
		return
	}

	//Start new terminal
	err := cmd.Start()
	if err != nil {
		fmt.Println("Feil ved åpning av terminal: ", err)
		return
	}
}
