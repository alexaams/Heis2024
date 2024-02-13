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

var number int

func main() {

	//Connect to UDP-socket
	udpAddr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		fmt.Println("Error resolving UDP address:", err)
		return
	}

	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		fmt.Println("Error listening on UDP port:", err)
		return
	}
	//Initiating backup-phase
	fmt.Println("---Initiating backup process---")
	buffer := make([]byte, 1024) //Creating buffer to receive data

	for {

		deadline := time.Now().Add(3 * time.Second)
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
		number = number_temp
		fmt.Printf("Mottatt %d bytes fra %s: %d\n", n, addr, (number))
	}
	conn.Close()
	//Start a new backup-process in a new terminal when no other is detected on UDP-broadcast
	openTerminal()
	//Switch over to primary process and start sending a number
	conn_send, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		fmt.Println("Feil ved å åpne UDP-port: ", err)
	}
	fmt.Println("---Initiating primary process---")
	for {
		_, err = conn_send.Write([]byte(fmt.Sprint(number)))
		fmt.Println(number)

		if err != nil {
			fmt.Println(err)
		}

		number++
		time.Sleep(time.Second * 1)
	}
	//Terminate terminal.

}

func openTerminal() {
	// Bestem kommandoen for å åpne en ny terminal og kjøre go run main.go
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		// For Windows bruker vi "cmd.exe" for å åpne en ny terminal og kjøre go run main.go
		cmd = exec.Command("cmd.exe", "/c", "start", "cmd.exe", "/k", "go", "run", "main.go")
	default:
		fmt.Println("Dette operativsystemet støttes ikke")
		return
	}

	// Start den nye terminalen
	err := cmd.Start()
	if err != nil {
		fmt.Println("Feil ved åpning av terminal: ", err)
		return
	}
}
