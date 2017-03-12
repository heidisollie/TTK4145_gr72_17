package main

import (
	"./FSM"
	"./driver"
	//"./localip"
	"./network"
	"./order_distribution"
	"./order_handler"
	"./structs"
	"fmt"
)

//Hei

func main() {
	fmt.Printf("Elev_driver started\n")
	driver.ElevInit()

	const NumFloors = driver.NumFloors
	const NumButtons = driver.NumButtons

	var localIP string
	localIP = network.GetIP()


	newTargetFloor := make(chan int, 1)
	floorEvent := make(chan int)                      //heis ved etasje til order_handler
	buttonEvent := make(chan driver.OrderButton, 100) //knappetrykk til order_handler
	//peers := make(chan int)                        //status of peer to peer fra network
	processNewOrder := make(chan structs.Order, 100)          //ekstern ordre fra order handler for kost funksjonen
	assignedNewOrder := make(chan structs.Order, 100) //ekstern ordre fra order_dist. , med heisID


	elevSendState := make(chan structs.ElevState, 100)
	elevSendNewOrder := make(chan structs.Order, 100)
	elevSendRemoveOrder := make(chan structs.Order, 100)
	

	elevReceiveState := make(chan structs.ElevState, 100)
	elevReceiveNewOrder := make(chan structs.Order, 100)
	elevReceiveRemoveOrder := make(chan structs.Order, 100)

	floorCompleted := make(chan int, 100)
	numberOfPeers := make(chan int, 100)

	driver.ElevInit()
	go driver.EventListener(buttonEvent, floorEvent)
	go FSM.FSMInit(floorEvent, newTargetFloor, floorCompleted, elevSendState)

	go order_handler.OrderHandlerInit(localIP,
		floorCompleted,
		buttonEvent,
		assignedNewOrder,
		processNewOrder,
		elevSendNewOrder,
		elevSendRemoveOrder,
		elevReceiveNewOrder,
		elevReceiveRemoveOrder,
		newTargetFloor)

	go network.NetworkInit(localIP,
		elevSendState,
		elevSendNewOrder,
		elevSendRemoveOrder,
		elevReceiveState,
		elevReceiveNewOrder,
		elevReceiveRemoveOrder,
		numberOfPeers)

	go order_distribution.OrderDistInit(localIP,
		processNewOrder,
		assignedNewOrder,
		elevReceiveState,
		numberOfPeers)


	select {}
}
