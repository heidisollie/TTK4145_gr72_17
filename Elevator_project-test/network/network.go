package network

import (
	def "../definitions"
	"./bcast"
	"./localip"
	"./peers"
	"flag"
	"fmt"
	"reflect"
)

const (
	portPeer        = 37899
	newOrderPort    = 37776
	removeOrderPort = 37715
	statePort       = 37714
	backupPort      = 37716
)

type UDPMessageState struct {
	Address string
	Data    def.Elevator
}

type UDPMessageOrder struct {
	Address string
	Data    def.Order
}

func GetIP() string {
	var id string
	flag.StringVar(&id, "id", "", "id of this peer")
	flag.Parse()
	localIP, err := localip.LocalIP()
	if err != nil {
		fmt.Println(err)
		localIP = "DISCONNECTED"
	}
	id = fmt.Sprintf(localIP)

	return id

}

func UDPPeerBroadcast(id string, numberOfPeers chan<- peers.PeerUpdate) {
	peerUpdateCh := make(chan peers.PeerUpdate)
	peerTxEnable := make(chan bool)

	go peers.Transmitter(portPeer, id, peerTxEnable)
	go peers.Receiver(portPeer, peerUpdateCh)

	for {
		select {
		case p := <-peerUpdateCh:
			fmt.Printf("Peer update:\n")
			fmt.Printf("  Peers:    %q\n", p.Peers)
			fmt.Printf("  New:      %q\n", p.New)
			fmt.Printf("  Lost:     %q\n", p.Lost)
			numberOfPeers <- p
		}
	}
}

func TransmitMessage(localIP string,
	elevSendState <-chan def.Elevator,
	elevSendNewOrder <-chan def.Order,
	elevSendRemoveOrder <-chan def.Order) {

	netSendState := make(chan UDPMessageState)
	netSendNewOrder := make(chan UDPMessageOrder)
	netSendRemoveOrder := make(chan UDPMessageOrder)

	go bcast.Transmitter(statePort, netSendState)
	go bcast.Transmitter(newOrderPort, netSendNewOrder)
	go bcast.Transmitter(removeOrderPort, netSendRemoveOrder)

	for {
		select {

		case msg := <-elevSendState:
			msg.IP = localIP
			message := UDPMessageState{Address: localIP, Data: msg}
			for i := 0; i < 3; i++ {
				netSendState <- message
			}
		case msg := <-elevSendNewOrder:
			message := UDPMessageOrder{Address: localIP, Data: msg}
			for i := 0; i < 3; i++ {
				netSendNewOrder <- message
			}

		case msg := <-elevSendRemoveOrder:
			message := UDPMessageOrder{Address: localIP, Data: msg}
			for i := 0; i < 3; i++ {
				netSendRemoveOrder <- message
			}

		}
	}
}

func ReceiveMessage(localIP string,
	elevReceiveState chan<- def.Elevator,
	elevReceiveNewOrder chan<- def.Order,
	elevReceiveRemoveOrder chan<- def.Order) {

	netReceiveState := make(chan UDPMessageState)
	netReceiveNewOrder := make(chan UDPMessageOrder)
	netReceiveRemoveOrder := make(chan UDPMessageOrder)

	go bcast.Receiver(statePort, netReceiveState)
	go bcast.Receiver(newOrderPort, netReceiveNewOrder)
	go bcast.Receiver(removeOrderPort, netReceiveRemoveOrder)

	for {
		select {

		case msg := <-netReceiveState:
			if msg.Address != localIP {
				msg.Data.IP = msg.Address
				elevReceiveState <- msg.Data
			}
		case msg := <-netReceiveNewOrder:
			if msg.Address != localIP {
				msg.Data.IP = msg.Address
				elevReceiveNewOrder <- msg.Data
			}
		case msg := <-netReceiveRemoveOrder:
			if msg.Address != localIP {
				msg.Data.IP = msg.Address
				elevReceiveRemoveOrder <- msg.Data
			}
		}
	}
}

func Repeater(ch_in interface{}, chs_out ...interface{}) {
	for {
		v, _ := reflect.ValueOf(ch_in).Recv()
		for _, c := range chs_out {
			reflect.ValueOf(c).Send(v)
		}
	}
}

func NetworkInit(
	localIP string,
	elevSendState <-chan def.Elevator,
	elevSendNewOrder <-chan def.Order,
	elevSendRemoveOrder <-chan def.Order,
	elevReceiveState chan<- def.Elevator,
	elevReceiveNewOrder chan<- def.Order,
	elevReceiveRemoveOrder chan<- def.Order,
	Peers chan<- peers.PeerUpdate) {

	fmt.Printf("Initializing network\n")

	go UDPPeerBroadcast(localIP, Peers)
	go TransmitMessage(localIP, elevSendState, elevSendNewOrder, elevSendRemoveOrder)
	go ReceiveMessage(localIP, elevReceiveState, elevReceiveNewOrder, elevReceiveRemoveOrder)
}
