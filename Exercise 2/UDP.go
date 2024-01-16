package main

import (
	//"bufio"
	"fmt"
	"net"
	//"os"
	"time"
)


func main(){

	address := "0.0.0.0:30000"

	udpAddr, err :=  net.ResolveUDPAddr("udp", address)
	if err != nil {
		fmt.Println("Feil ved å løse opp adresse: ", err)
		return
	}

	conn, err := net.ListenUDP("udp",udpAddr)

	if err!=nil {
		fmt.Println("Feil ved å åpne UDP-port: ", err)
	}

	defer conn.Close()

	fmt.Printf("Lytter på UDP-port %s\n",address)

	buffer := make([]byte,256)
	

	for {
		n, addr, err := conn.ReadFromUDP(buffer)
		if err!=nil{
			fmt.Println("Feil ved å lese data: ",err)
			continue
		}
		time.Sleep(1000*time.Millisecond)

		fmt.Printf("Mottatt %d bytes fra %s: %s\n", n, addr, string(buffer))
	}

	// ip udp server: 10.100.23.129 

	
}