package network

import (
	. "../structs"
	"./bcast"
	"./localip"
	"./peers"
	"flag"
	"fmt"
	"os"
	"time"
	//"../localState"
)

const (
	port_peer          = 37899
	get_order_port     = 37776
	cost_value_port    = 37714
	remove_order_port  = 37715
	backup_port        = 37716
	broadcast_interval = 1 * time.Second
)

type UDPmessage_cost struct {
	Address string
	Data    Cost
}

type UDPmessage_order struct {
	Address string
	Data    Order
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


func UDPPeerBroadcast(id string, number_of_peers chan<- int) {
	peerUpdateCh := make(chan peers.PeerUpdate)
	peerTxEnable := make(chan bool)

	go peers.Transmitter(port_peer, id, peerTxEnable)
	go peers.Receiver(port_peer, peerUpdateCh)

	for {
		select {
		case p := <-peerUpdateCh:
			fmt.Printf("Peer update:\n")
			fmt.Printf("  Peers:    %q\n", p.Peers)
			fmt.Printf("  New:      %q\n", p.New)
			fmt.Printf("  Lost:     %q\n", p.Lost)
			number_of_peers <- len(p.Peers)
		}
	}
}


func TransmitMsg(localIP string,
	elev_send_cost_value <-chan Cost,
	elev_send_new_order <-chan Order,
	elev_send_remove_order <-chan Order ){

	net_send_cost_value := make(chan UDPmessage_cost, 100)
	net_send_new_order := make(chan UDPmessage_order, 100)
	net_send_remove_order := make(chan UDPmessage_order, 100)


	go bcast.Transmitter(get_order_port, net_send_new_order)
	go bcast.Transmitter(remove_order_port, net_send_remove_order)
	go bcast.Transmitter(cost_value_port, net_send_cost_value)

	for {
		select {
	
			case msg := <-elev_send_new_order:
				message := UDPmessage_order{Address: localIP, Data: msg}
				net_send_new_order <- message
				fmt.Printf("BROADCASTING NEW ORDER\n")
	
			case msg := <-elev_send_remove_order:
				message := UDPmessage_order{Address: localIP, Data: msg}
				net_send_remove_order <- message
				fmt.Printf("BROADCASTING REMOVE ORDER\n")
	
			case msg := <-elev_send_cost_value:
				fmt.Printf("BROADCASTING COST VALUE\n")
				msg.Current_order.IP = localIP
				message := UDPmessage_cost{Address: localIP, Data: msg}
				net_send_cost_value <- message
	
		}
	}
}
func ReceiveMsg(localIP string,
	elev_receive_cost_value chan<- Cost,
	elev_receive_new_order chan<- Order,
	elev_receive_remove_order chan<- Order) {


	net_receive_cost_value := make(chan UDPmessage_cost, 100)
	net_receive_new_order := make(chan UDPmessage_order, 100)
	net_receive_remove_order := make(chan UDPmessage_order, 100)



	go bcast.Receiver(get_order_port, net_receive_new_order)
	go bcast.Receiver(remove_order_port, net_receive_remove_order)
	go bcast.Receiver(cost_value_port, net_receive_cost_value)

	for {
		select {
			case msg := <-net_receive_new_order:
				if msg.Address != localIP {
					fmt.Printf("RECEIVING NEW ORDER\n")
					msg.Data.IP = msg.Address
					fmt.Printf("ADDRESS1: %s\n", msg.Data.IP)
					elev_receive_new_order <- msg.Data
				}
			case msg := <-net_receive_remove_order:
				if msg.Address != localIP {
					fmt.Printf("RECEIVING REMOVE ORDER\n")
					msg.Data.IP = msg.Address
					fmt.Printf("ADDRESS2: %s\n", msg.Data.IP)
					elev_receive_remove_order <- msg.Data
				}
			case msg := <-net_receive_cost_value:
				if msg.Address != localIP {
					fmt.Printf("RECEIVING COST VALUE\n")
					msg.Data.Current_order.IP = msg.Address
					fmt.Printf("ADDRESS: %s\n", msg.Data.Current_order.IP)
					elev_receive_cost_value <- msg.Data	
				}		
		}
	}

}

func UDP_init(
	localIP string,
	elev_receive_cost_value chan<- Cost,
	elev_receive_new_order chan<- Order,
	elev_receive_remove_order chan<- Order,
	elev_send_cost_value <-chan Cost,
	elev_send_new_order <-chan Order,
	elev_send_remove_order <-chan Order,
	number_of_peers chan<- int) {

	fmt.Printf("Initializing network\n")

	go UDPPeerBroadcast(localIP, number_of_peers)

	go TransmitMsg(localIP, elev_send_cost_value, elev_send_new_order, elev_send_remove_order)

	go ReceiveMsg(localIP, elev_receive_cost_value, elev_receive_new_order, elev_receive_remove_order)

}






/*
func UDP_init(
	elev_receive_cost_value chan<- Cost,
	elev_receive_new_order chan<- Order,
	elev_receive_remove_order chan<- Order,
	elev_send_cost_value <-chan Cost,
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
	net_send_cost_value := make(chan<- UDPmessage_cost)
	net_send_new_order := make(chan<- UDPmessage_order)
	net_send_remove_order := make(chan<- UDPmessage_order)

	net_receive_cost_value := make(<-chan UDPmessage_cost)
	net_receive_new_order := make(<-chan UDPmessage_order)
	net_receive_remove_order := make(<-chan UDPmessage_order)

	peerUpdateCh := make(chan peers.PeerUpdate)
	peerTxEnable := make(chan bool)

	//binding channels and ports
	go peers.Transmitter(port_peer, id, peerTxEnable)
	go peers.Receiver(port_peer, peerUpdateCh)

	go bcast.Transmitter(get_order_port, net_send_new_order)
	go bcast.Transmitter(remove_order_port, net_send_remove_order)
	go bcast.Transmitter(cost_value_port, net_send_cost_value)

	go bcast.Receiver(get_order_port, net_receive_new_order)
	go bcast.Receiver(remove_order_port, net_receive_remove_order)
	go bcast.Receiver(cost_value_port, net_receive_cost_value)

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

		case msg := <-elev_send_cost_value:
			fmt.Printf("Broadcasting cost value\n")
			for {
				message := UDPmessage_cost{Address: localIP, Data: msg}
				net_send_cost_value <- message
				time.Sleep(broadcast_interval)
			}

		//cases where NW receives data from the network and passes it to the right channel
		case msg := <-net_receive_new_order:
			fmt.Printf("Received new order from NW\n")
			elev_receive_new_order <- msg.Data

		case msg := <-net_receive_remove_order:
			fmt.Printf("Received remove order from NW\n")
			elev_receive_remove_order <- msg.Data

		case msg := <-net_receive_cost_value:
			fmt.Printf("Received cost value from NW\n")
			elev_receive_cost_value <- msg.Data

		}
	}
}
*/
