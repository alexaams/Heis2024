package peers

import (
	"ProjectHeis/Scratch/config_folder/types"
	"ProjectHeis/Scratch/config_folder/globals"
	"ProjectHeis/drivers/elevator"
	"ProjectHeis/network/bcast"
	"ProjectHeis/network/conn"
	"ProjectHeis/network/localip"
	"fmt"
	"net"
	"sort"
	"strconv"
	"time"
)

// ___________Global variables___________
var G_Ch_PeersData_Tx = make(chan PeersData)
var G_Ch_PeersData_Rx = make(chan PeersData)
var G_PeersUpdate PeerUpdate

type PeerUpdate struct {
	Peers []string
	New   string
	Lost  []string
}

type PeersData struct {
	Elevator         elevator.Elevator
	Id               int
	SingleOrdersHall types.OrdersHall
	GlobalOrderHall  types.OrdersHall
}

const interval = 15 * time.Millisecond
const timeout = 500 * time.Millisecond

func Transmitter(port int, id string, transmitEnable <-chan bool) {

	conn := conn.DialBroadcastUDP(port)
	addr, _ := net.ResolveUDPAddr("udp4", fmt.Sprintf("255.255.255.255:%d", port))

	enable := true
	for {
		select {
		case enable = <-transmitEnable:
		case <-time.After(interval):
		}
		if enable {
			conn.WriteTo([]byte(id), addr)
		}
	}
}

func Receiver(port int, peerUpdateCh chan<- PeerUpdate) {
	//This code is specifically made for peers and checking if there are new peers/lost peers.
	//Should not be used for anything else.
	var buf [1024]byte
	var p PeerUpdate
	lastSeen := make(map[string]time.Time)

	conn := conn.DialBroadcastUDP(port)

	for {
		updated := false

		conn.SetReadDeadline(time.Now().Add(interval))
		n, _, _ := conn.ReadFrom(buf[0:])

		id := string(buf[:n])

		// Adding new connection
		p.New = ""
		if id != "" {
			if _, idExists := lastSeen[id]; !idExists {
				p.New = id
				updated = true
			}

			lastSeen[id] = time.Now()
		}

		// Removing dead connection
		p.Lost = make([]string, 0)
		for k, v := range lastSeen {
			if time.Now().Sub(v) > timeout {
				updated = true
				p.Lost = append(p.Lost, k)
				delete(lastSeen, k)
			}
		}

		// Sending update
		if updated {
			p.Peers = make([]string, 0, len(lastSeen))

			for k, _ := range lastSeen {
				p.Peers = append(p.Peers, k)
			}

			sort.Strings(p.Peers)
			sort.Strings(p.Lost)
			peerUpdateCh <- p
		}
	}
}

func InitPeers() PeersData {
	return PeersData{
		Elevator:         elevator.InitElevator(),
		Id:               localip.CreateID(),
		SingleOrdersHall: types.InitEmptyOrder(),
		GlobalOrderHall:  types.InitEmptyOrder(),
	}
}

func SendPeersData_init() {
	go bcast.Transmitter(16580, G_Ch_PeersData_Tx)
	go bcast.Receiver(16580, G_Ch_PeersData_Rx)
}

func PeersHeartBeat() {
	config.ElevatorID = localip.CreateID()

	fmt.Printf("Our ID is: %d\n", config.ElevatorID)

	peerUpdateCh := make(chan PeerUpdate)
	peerTxEnable := make(chan bool)

	go Transmitter(15659, strconv.Itoa(config.ElevatorID), peerTxEnable)
	go Receiver(15659, peerUpdateCh)

	fmt.Println("Heartbeat-sequency initiated")
	for {
		select {
		case p := <-peerUpdateCh:
			G_PeersUpdate = p
			p.PrintPeersUpdate()
			//Sende data videre til kostfunksjon
		}
	}
}

func (p PeerUpdate) PrintPeersUpdate() {
	fmt.Printf("Peer update:\n")
	fmt.Printf("  Peers:    %q\n", p.Peers)
	fmt.Printf("  New:      %q\n", p.New)
	fmt.Printf("  Lost:     %q\n", p.Lost)
}