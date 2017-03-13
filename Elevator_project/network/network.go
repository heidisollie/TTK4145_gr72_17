package network

import (
	"../structs"
	"./bcast"
	"./localip"
	"./peers"
	"flag"
	"fmt"
	"os"
	"time"
)

const (
	portPeer          = 37899
	newOrderPort      = 37776
	removeOrderPort   = 37715
	statePort         = 37714
	backupPort        = 37716
	broadcastInterval = 1 * time.Second
)

type UDPMessageState struct {
	Address string
	Data    structs.ElevState
}

type UDPMessageOrder struct {
	Address string
	Data    structs.Order
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
	id = fmt.Sprintf("peer-%s-%d", localIP, os.Getpid())

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
	elevSendState <-chan structs.ElevState,
	elevSendNewOrder <-chan structs.Order,
	elevSendRemoveOrder <-chan structs.Order) {

	netSendState := make(chan UDPMessageState, 100)
	netSendNewOrder := make(chan UDPMessageOrder, 100)
	netSendRemoveOrder := make(chan UDPMessageOrder, 100)

	go bcast.Transmitter(statePort, netSendState)
	go bcast.Transmitter(newOrderPort, netSendNewOrder)
	go bcast.Transmitter(removeOrderPort, netSendRemoveOrder)

	for {
		select {

		case msg := <-elevSendState:
			msg.IP = localIP
			message := UDPMessageState{Address: localIP, Data: msg}
			for {
				netSendState <- message
			}
			fmt.Printf("BROADCASTING STATE\n")
		case msg := <-elevSendNewOrder:
			message := UDPMessageOrder{Address: localIP, Data: msg}
			for {
				netSendNewOrder <- message
			}
			fmt.Printf("BROADCASTING NEW ORDER\n")

		case msg := <-elevSendRemoveOrder:
			message := UDPMessageOrder{Address: localIP, Data: msg}
			for {
				netSendRemoveOrder <- message
			}
			fmt.Printf("BROADCASTING REMOVE ORDER\n")

		}
	}
}
func ReceiveMessage(localIP string,
	elevReceiveState chan<- structs.ElevState,
	elevReceiveNewOrder chan<- structs.Order,
	elevReceiveRemoveOrder chan<- structs.Order) {

	netReceiveState := make(chan UDPMessageState, 100)
	netReceiveNewOrder := make(chan UDPMessageOrder, 100)
	netReceiveRemoveOrder := make(chan UDPMessageOrder, 100)

	go bcast.Receiver(statePort, netReceiveState)
	go bcast.Receiver(newOrderPort, netReceiveNewOrder)
	go bcast.Receiver(removeOrderPort, netReceiveRemoveOrder)

	for {
		select {

		case msg := <-netReceiveState:
			if msg.Address != localIP {
				fmt.Printf("RECEIVING STATE\n")
				msg.Data.IP = msg.Address
				fmt.Printf("ADDRESS: %s \n", msg.Data.IP)
				elevReceiveState <- msg.Data
			}
		case msg := <-netReceiveNewOrder:
			if msg.Address != localIP {
				fmt.Printf("RECEIVING NEW ORDER\n")
				msg.Data.IP = msg.Address
				fmt.Printf("ADDRESS1: %s \n", msg.Data.IP)
				elevReceiveNewOrder <- msg.Data
			}
		case msg := <-netReceiveRemoveOrder:
			if msg.Address != localIP {
				fmt.Printf("RECEIVING REMOVE ORDER\n")
				msg.Data.IP = msg.Address
				fmt.Printf("ADDRESS2: %s \n", msg.Data.IP)
				elevReceiveRemoveOrder <- msg.Data
			}
		}
	}
}

func NetworkInit(
	localIP string,
	elevSendState <-chan structs.ElevState,
	elevSendNewOrder <-chan structs.Order,
	elevSendRemoveOrder <-chan structs.Order,
	elevReceiveState chan<- structs.ElevState,
	elevReceiveNewOrder chan<- structs.Order,
	elevReceiveRemoveOrder chan<- structs.Order,
	Peers chan<- peers.PeerUpdate) {

	fmt.Printf("Initializing network\n")

	go UDPPeerBroadcast(localIP, Peers)
	go TransmitMessage(localIP, elevSendState, elevSendNewOrder, elevSendRemoveOrder)
	go ReceiveMessage(localIP, elevReceiveState, elevReceiveNewOrder, elevReceiveRemoveOrder)

}

/*
func UDP_init(
	elev_receive_state chan<- Cost,
	elev_receive_new_order chan<- Order,
	elev_receive_remove_order chan<- Order,
	elev_send_state <-chan Cost,
	elev_send_new_order <-chan Order,
	elev_send_remove_order <-chan Order) {

	fmt.Printf("Initializing network\n")

	var id string

		flag.StringVar(&id, "id", "", "id of this peer")
		flag.Parse()

			var localIP string
			localIP, err := localip.LocalIP()
			if err != nil {
				fmt.Println(err)
				localIP = "DISCONNECTED"
			}
			id = fmt.Sprintf("peer-%s-%d", localIP, os.Getpid())

	var localIP string
	localIP = "3456"
	id = localIP
	//channels for network
	net_send_state := make(chan<- UDPmessage_cost)
	net_send_new_order := make(chan<- UDPmessage_order)
	net_send_remove_order := make(chan<- UDPmessage_order)

	net_receive_state := make(<-chan UDPmessage_cost)
	net_receive_new_order := make(<-chan UDPmessage_order)
	net_receive_remove_order := make(<-chan UDPmessage_order)

	peerUpdateCh := make(chan peers.PeerUpdate)
	peerTxEnable := make(chan bool)

	//binding channels and ports
	go peers.Transmitter(port_peer, id, peerTxEnable)
	go peers.Receiver(port_peer, peerUpdateCh)

	go bcast.Transmitter(get_order_port, net_send_new_order)
	go bcast.Transmitter(remove_order_port, net_send_remove_order)
	go bcast.Transmitter(state_port, net_send_state)

	go bcast.Receiver(get_order_port, net_receive_new_order)
	go bcast.Receiver(remove_order_port, net_receive_remove_order)
	go bcast.Receiver(state_port, net_receive_state)

	//send_ticker := time.NewTicker(broadcast_interval) // bruke dette?

	for {
		select {

		//cases where NW recieves message from elevatar and broadcastes it on the network
		case msg := <-elev_send_new_order:
			fmt.Printf("Broadcasting new order from elev\n")
			for {
				message := UDPmessage_order{Address: localIP, Data: msg}
				net_send_new_order <- message
				time.Sleep(broadcast_interval)
			}

		case msg := <-elev_send_remove_order:
			fmt.Printf("Broadcasting remove order from elev\n")
			for {
				message := UDPmessage_order{Address: localIP, Data: msg}
				net_send_remove_order <- message
				time.Sleep(broadcast_interval)
			}

		case msg := <-elev_send_state:
			fmt.Printf("Broadcasting cost value\n")
			for {
				message := UDPmessage_cost{Address: localIP, Data: msg}
				net_send_state <- message
				time.Sleep(broadcast_interval)
			}

		//cases where NW receives data from the network and passes it to the right channel
		case msg := <-net_receive_new_order:
			fmt.Printf("Received new order from NW\n")
			elev_receive_new_order <- msg.Data

		case msg := <-net_receive_remove_order:
			fmt.Printf("Received remove order from NW\n")
			elev_receive_remove_order <- msg.Data

		case msg := <-net_receive_state:
			fmt.Printf("Received cost value from NW\n")
			elev_receive_state <- msg.Data

		}
	}
}
*/
