package localip

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

var localIP string

func LocalIP() (string, error) {
	if localIP == "" {
		conn, err := net.Dial("udp4", "8.8.8.8:53")
		if err != nil {
			return "", err
		}
		defer conn.Close()
		localIP = strings.Split(conn.LocalAddr().String(), ":")[0]
	}
	return localIP, nil
}

// Creating ID with local ip and PID
func CreateID() int {
	idStr := ""

	if idStr == "" {
		localIP, err := LocalIP()
		if err != nil {
			fmt.Println(err)
			localIP = "DISCONNECTED"
		}
		idStr = localIP
		temp_arr := strings.Split(idStr, ".")
		idStr = temp_arr[3]

	}
	idInt, err := strconv.Atoi(idStr)
	if err != nil {
		fmt.Println("Error Converting stringID to Int: ", err)
	}
	return idInt
}
