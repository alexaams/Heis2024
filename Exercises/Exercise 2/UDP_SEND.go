package main

import (
	//"bufio"
	"fmt"
	"net"

	//"os"
	"time"
)

const send_address string = "10.100.23.129:20005"
const listen_address string = "0.0.0.0:20005"

func main() {

	udpAddr_send, err := net.ResolveUDPAddr("udp", send_address)
	if err != nil {
		fmt.Println("Feil ved å løse opp adresse: ", err)
		return
	}

	conn_send, err := net.DialUDP("udp", nil, udpAddr_send)
	if err != nil {
		fmt.Println("Feil ved å åpne UDP-port: ", err)
	}

	udpAddr, err := net.ResolveUDPAddr("udp", listen_address)
	if err != nil {
		fmt.Println("Feil ved å løse opp adresse: ", err)
		return
	}

	conn, err := net.ListenUDP("udp", udpAddr)

	if err != nil {
		fmt.Println("Feil ved å åpne UDP-port: ", err)
	}

	defer conn.Close()
	defer conn_send.Close()

	buffer := make([]byte, 256)

	for {
		_, err = conn_send.Write([]byte("Halloi fra Alex"))
		fmt.Println("Send")
		if err != nil {
			fmt.Println(err)
		}
		n, addr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println("Feil ved å lese data: ", err)
			continue
		}

		fmt.Printf("Mottatt %d bytes fra %s: %s\n", n, addr, string(buffer))
		time.Sleep(1000 * time.Millisecond)

	}

	// ip udp server: 10.100.23.129

}
