package network

import (
	"./bcast"
	"./localip"
	"./peers"
	"flag"
	"fmt"
	"os"
	"time"
)


const (
	port_peer = "20002"
	get_order_port = "37712"
	cost_value_port = "37714"
	remove_order_port = "37715"
	backup_port = "37716"
	broadcast_interval  = 100 * time.Millisecond
) 

type UDPmessage struct {
		Address IP
		Data []byte
	}



	



func UDP_init(
	elev_receive_cost_value chan <-Cost,
	elev_receive_new_order chan <-Order,
	elev_receive_remove_order chan <- Order,

	elev_send_cost_value <-chan Cost,
	elev_send_new_order <- chan Order,
	elev_send_remove_order <- chan Order
	){
	
	fmt.Printf("initialazing network\n")
	
	//channels for network 
	net_send_cost_value := make(chan Cost)
	net_send_new_order := make(chan Order)
	net_send_remove_order := make(chan Order)
	net_receive_cost_value := make(chan int)
	net_receive_new_order := make(chan Order)
	net_receive_remove_order := make (chan Order)
	
	
	net_send_cost_value <- chan Cost
	net_send_new_order <- chan Order
	net_send_remove_order <- chan Order

	net_receive_cost_ value chan <- Cost
	net_receive_new_order chan <- Order
	net_receive_remove_order chan <- Order
	
	
	//binding channels and ports
	go peers.Transmitter(port_peer, id, peerTxEnable)
	go peers.Receiver(port_peer, peerUpdateCh)

	go bcast.Transmitter(get_order_port, net_send_new_order)
	go bcast.Transmitter(remove_order_port, net_send_remove_order)
	go bcast.Transmitter(cost_value_port, net_send_cost_value)

	go bcast.Receiver(get_order_port, net_receive_new_order)
	go bcast.Receiver(remove_order_port, net_receive_remove_order)
	go bcast.Receiver(cost_value_port, net_receive_cost_value)
	
	send_ticker := time.NewTicker(broadcast_interval) // bruke dette?
	
}




func main() {

	var id string
	flag.StringVar(&id, "id", "", "id of this peer")
	flag.Parse()

	if id == "" {
		localIP, err := localip.LocalIP()
		if err != nil {
			fmt.Println(err)
			localIP = "DISCONNECTED"
		}
		id = fmt.Sprintf("peer-%s-%d", localIP, os.Getpid())
	}
	
	peerUpdateCh := make(chan peers.PeerUpdate)
	peerTxEnable := make(chan bool)

	for{
		select{

			//cases where NW recieves message from elevatar and broadcastes it on the network
			case msg := <- elev_send_new_order:
				fmt.Printf("Broadcasting new order\n")
				for {
					message := UDPmessage{Address: localIP, Data: msg}
					net_send_new_order <- message
					
					time.Sleep(broadcast_interval)
				}
				
				
			case msg := <- elev_send_remove_order:
				fmt.Printf("Broadcasting remove order\n")
				for {
					message := UDPmessage{Address: localIP, Data: msg}
					net_send_remove_order <- message
					time.Sleep(broadcast_interval)
				}
				
			case msg := <- elev_send_cost_value:
				fmt.Printf("Broadcasting cost value\n")
				for {
					message := UDPmessage{Address: localIP, Data: msg}
					net_send_cost_value <- message
					time.Sleep(broadcast_interval)
				}
				

			//cases where NW receives data from the network and passes it to the right channel
			case msg := <- net_receive_new_order:
				fmt.Printf("Received new order from NW\n")
				elev_send_new_order <- msg

			
			case msg := <- net_receive_remove_order:
				fmt.Printf("Received remove order from NW\n")
				elev_receive_remove_order <- msg


			case msg := <- net_receive_cost_value:
				fmt.Printf("Rceived cost value from NW\n")
				elev_receive_cost_value <- msg
			
		}
	}
}
