package main

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"time"
)

const address = "localhost:8080"

var last_time_stamp time.Time

func main() {
	//last_time_stamp = time.Now()

	fmt.Println("test")
	//channel_timer := make(chan int)

	backup()
	openTerminal()
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

func watch_dog(time_stamp time.Time, time_limit int) {
	if time.Now().Second()-time_stamp.Second() > time_limit {
		fmt.Println("Time-limit exceeded")
	}
}

func backup() {
	fmt.Println("-- Backup phase --")
	udpAddr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		fmt.Println("Feil ved å løse opp adresse: ", err)
		return
	}

	conn, _ := net.ListenUDP("udp", udpAddr)

	buffer := make([]byte, 1024)
	deadline := time.Now().Add(3 * time.Second)

	for {
		conn.SetReadDeadline(deadline)

		n, addr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println("--Primary phase--")
			break
		}

		fmt.Printf("Mottatt %d bytes fra %s: %s\n", n, addr, string(buffer))
		//channel <- 1
	}
	conn.Close()

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
	var k int = 0
	for {
		_, err = conn.Write([]byte(fmt.Sprint(k)))
		fmt.Println(k)
		if err != nil {
			fmt.Println(err)
		}
		k++
		if k > 6 {
			break
		}
		time.Sleep(time.Second)
	}
	terminate()
}

func openTerminal() {
	cmd := exec.Command("gnome-terminal", "--", "go", "run", "main.go")
	err := cmd.Start()
	if err != nil {
		fmt.Println("Feil ved åpning av terminal: ", err)
		return
	}

	// Vent til terminalvinduet er ferdig før du fortsetter
	err = cmd.Wait()
	if err != nil {
		fmt.Println("Feil ved venting på terminal: ", err)
	}
}

func terminate() {
	fmt.Println("Program terminated")
	os.Exit(3)
}
