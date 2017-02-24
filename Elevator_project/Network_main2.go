package main

import (
	"./network/bcast"
	"./network/localip"
	"./network/peers"
	"flag"
	"fmt"
	"os"
)


const port_peer int = 20002
const port1 int = 37712
const port2 int = 37713
const port3 int = 37714


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

	// We make a channel for receiving updates on the id's of the peers that are
	//  alive on the network
	peerUpdateCh := make(chan peers.PeerUpdate)
	// We can disable/enable the transmitter after it has been started.
	// This could be used to signal that we are somehow "unavailable".
	peerTxEnable := make(chan bool)
	
	go peers.Transmitter(port_peer, id, peerTxEnable)
	go peers.Receiver(port_peer, peerUpdateCh)
	
	
	sendCostValue := make(chan string)
	sendNewOrder := make(chan string)
	sendRemoveOrder := make(chan string)
	receiveCostValue := make(chan string)
	receiveNewOrder := make(chan string)
	receiveRemoveOrder := make (chan string)
	
	go bcast.Transmitter(port1, sendCostValue)
	go bcast.Receiver(port1, receiveCostValue)
	go bcast.Transmitter(port2, sendNewOrder)
	go bcast.Receiver(port2, receiveNewOrder)
	go bcast.Transmitter(port3, sendRemoveOrder)
	go bcast.Receiver(port3, receiveRemoveOrder)
	
	
	
	
	for{
		select{
			case msg := <- receiveCostValue:
				fmt.Printf("Received cost value: %s\n", msg)
			
			case msg := <- receiveNewOrder:
				fmt.Printf("Received new order: %s\n", msg)
			
			case msg := <- receiveRemoveOrder:
				fmt.Printf("Received remove order: %s\n", msg)
		
		}
	}
}
