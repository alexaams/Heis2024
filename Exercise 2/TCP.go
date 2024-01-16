package main

import (
	"fmt"
	"net"
	"time"
)

func reciever(conn net.Conn) {

	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	fmt.Println("Recieved: ", string(buffer[0:n]))
	if err != nil {
		fmt.Println("Feil ved lesing")
	}
}

func send(conn net.Conn) {

	msg_init := "Connect to: 10.100.23.15:34933\x00"
	buffer := []byte(msg_init)
	_, err := conn.Write(buffer)
	if err != nil {
		fmt.Println(msg_init, err)
	}
	time.Sleep(time.Second * 1)
	for {
		msg := "test gruppe 5\x00"
		buffer := []byte(msg)
		_, err := conn.Write(buffer)
		if err != nil {
			fmt.Println(err)
		}
		time.Sleep(time.Second * 1)
	}
}

func main() {
	conn, err := net.Dial("tcp", "10.100.23.129:34933")
	if err != nil {
		fmt.Println(err)
	}

	go send(conn)
	go reciever(conn)

	for {
		time.Sleep(time.Second * 10)
	}
}
