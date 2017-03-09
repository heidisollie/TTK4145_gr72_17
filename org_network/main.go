package main

import (
	"./network/bcast"
	"./network/localip"
	"./network/peers"
	"flag"
	"fmt"
	"os"
	"time"
)


const TPort_peer = 15657
const RPort_peer = 15657

const TPort = 13034
const RPort = 13034

type UDPMSG_cost struct {
	Adress string
	Cost Cost
	Iter int
	
}

type Order struct {
	Type     ButtonType
	Floor    int
	Internal bool
	IP       string
}

type Cost struct {
	Cost_value    int
	Current_order Order
	IP            string
}


type ButtonType int


const (
	ButtonCallDown    ButtonType = 0
	ButtonCallCommand ButtonType = 1
	ButtonCallUp      ButtonType = 2
)

func main() {
	// Our id can be anything. Here we pass it on the command line, using
	//  `go run main.go -id=our_id`
	var id string
	flag.StringVar(&id, "id", "", "id of this peer")
	flag.Parse()

	// ... or alternatively, we can use the local IP address.
	// (But since we can run multiple programs on the same PC, we also append the
	//  process ID)
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

	go peers.Transmitter(TPort_peer, id, peerTxEnable)
	go peers.Receiver(RPort_peer, peerUpdateCh)


	helloTx := make(chan UDPMSG_cost)
	helloRx := make(chan UDPMSG_cost)

	go bcast.Transmitter(TPort, helloTx)
	go bcast.Receiver(RPort, helloRx)


	go func() {	
		order := Order{ButtonCallUp, 3, true, id}
		Cost := Cost{2, order, id}
		Cost_msg := UDPMSG_cost{id, Cost, 0}
		for {
			Cost_msg.iter++
			helloTx <- Cost_msg
			fmt.Printf("[SEND] Cost value: %d\n", Cost_msg.Cost.Cost_value)
			fmt.Printf("[SEND] Floor: %d\n", Cost_msg.Cost.Current_order.Floor)
			fmt.Printf("[SEND] Button: %d\n", Cost_msg.Cost.Current_order.Type)
			time.Sleep(2 * time.Second)
		}
	}()

	fmt.Println("Started")
	for {
		select {
		case p := <-peerUpdateCh:
			fmt.Printf("Peer update:\n")
			fmt.Printf("  Peers:    %q\n", p.Peers)
			fmt.Printf("  New:      %q\n", p.New)
			fmt.Printf("  Lost:     %q\n", p.Lost)

		case c:= <-helloRx:
			fmt.Printf("[RECEIVE] Cost value: %d\n", c.Cost.Cost_value)
			fmt.Printf("[RECEIVE] Floor: %d\n", c.Cost.Current_order.Floor)
			fmt.Printf("[RECEIVE] Button type: %d\n", c.Cost.Current_order.Type)
		}
	}
}
\