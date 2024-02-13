package main

import (
	"fmt"
	"net"
	"os/exec"
	"strconv"
	"time"
)

const address = "127.0.0.1:27106"

var message int

//var last_time_stamp time.Time

func main() {
	//last_time_stamp = time.Now()

	fmt.Println("test")
	//channel_timer := make(chan int)

	backup()

	primary()
	//go primary()

	//go watch_dog(last_time_stamp, 2)
	//for {
	//select {
	// case a := <-channel_timer:
	// 	last_time_stamp = time.Now()
	// 	println(a)
	//}

	//}

}

// func watch_dog(time_stamp time.Time, time_limit int) {
// 	if time.Now().Second()-time_stamp.Second() > time_limit {
// 		fmt.Println("Time-limit exceeded")
// 	}
// }

func backup() {

	fmt.Println("-- Backup phase --")
	udpAddr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		fmt.Println("Feil ved å løse opp adresse: ", err)
		return
	}

	conn, _ := net.ListenUDP("udp", udpAddr)

	defer openTerminal()
	defer conn.Close()

	buffer := make([]byte, 1024)

	for {
		deadline := time.Now().Add(2 * time.Second)
		conn.SetReadDeadline(deadline)

		n, addr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println("--Primary phase--")
			break
		}
		message, _ = strconv.Atoi(string(buffer[:n]))

		fmt.Printf("Mottatt %d bytes fra %s: %s\n", n, addr, fmt.Sprint(message))
		//channel <- 1
	}

}

func primary() {

	udpAddr_send, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		fmt.Println("Feil ved å løse opp adresse: ", err)
		return
	}

	conn, err := net.DialUDP("udp", nil, udpAddr_send)
	if err != nil {
		fmt.Println("Feil ved å åpne UDP-port: ", err)
	}
	defer conn.Close()
	for {
		message++
		_, err = conn.Write([]byte(strconv.Itoa(message)))
		fmt.Println(message)
		if err != nil {
			fmt.Println(err)
		}
		time.Sleep(time.Second)

	}

}

func openTerminal() {
	cmd := exec.Command("gnome-terminal", "--", "go", "run", "main.go")
	err := cmd.Start()
	if err != nil {
		fmt.Println("Feil ved åpning av terminal: ", err)
		return
	}

	// Wait until the terminal window is closed before continuing
	err = cmd.Wait()
	if err != nil {
		fmt.Println("Feil ved venting på terminal: ", err)
	}

}

//func terminate() {
//	fmt.Println("Program terminated")
//	os.Exit(3)
//}
