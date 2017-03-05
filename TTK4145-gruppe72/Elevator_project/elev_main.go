package main

import (
	"./FSM"
	"./driver"
	"./network"
	"./order_distribution"
	"./order_handler"
	"./structs"
	"fmt"
	"time"
)

func main() {
	fmt.Printf("Elev_driver started\n")
	driver.ElevInit()

	const NumFloors = driver.NumFloors
	const NumButtons = driver.NumButtons

	localIP, err := localip.LocalIP()
	if err != nil {
		fmt.Println(err)
		localIP = "DISCONNECTED"
	}

	new_target_floor := make(chan int)
	order_motor_dir := make(chan int)
	floor_event := make(chan int)          //heis ved etasje til order_handler
	button_event := make(chan OrderButton) //knappetrykk til order_handler
	peers := make(chan int)                //status of peer to peer fra network
	new_order := make(chan Order)          //ekstern ordre fra order handler for kost funksjonen
	assigned_new_order := make(chan Order) //ekstern ordre fra order_dist. , med heisID

	//Network channels
	elev_send_cost_value = make(chan Cost)
	elev_send_new_order = make(chan Order)
	elev_send_remove_order = make(chan Order)
	elev_receive_cost_value = make(chan int)
	elev_receive_new_order = make(chan Order)
	elev_receive_remove_order = make(chan Order)
	floor_completed = make(chan int)

	//Network
	peers := make(chan int) //status of peer to peer til order_dist.

	//Initialisere
	go FSM.FSM_init(floor_event, new_target_floor, floor_completed)
	go order_handler.order_handler_init(floor_completed, button_event, assigned_new_order, new_order, new_target_floor)

	go network.UPD_init(elev_send_cost_value,
		elev_send_new_order,
		elev_send_remove_order,
		elev_receive_cost_value,
		elev_receive_new_order,
		elev_receive_remove_order)

	go order_distribution.order_dist_init(new_order,
		assigned_new_order,
		elev_receive_cost_value,
		elev_send_cost_value)

	go driver.ElevInit()
	go driver.EventListener(button_event, floor_event)
	fmt.Printf("Elev_driver ended\n")
}

/*
//FSM
setMotorDir := make(chan<- int)  //set motor direction fra order_handler

//Elev_driver
floorEvent := make(<-chan int) //heis ved etasje til order_handler
buttonEvent := make(<-chan OrderButton) //knappetrykk til order_handler

//Order_dist.
peers := make(chan<- int) //status of peer to peer fra network
sendCostValue := make(<-chan int) //kostverdi til nettverk for broadcasting
receiveCostValue := make(chan<- int) //kostverdig fra nettverk til action_select
newOrder := make(chan <- Order) //ekstern ordre fra order handler for kost funksjonen
assignedNewOrder := make(<- chan Order) //ekstern ordre til order handler, med heisID

//Order_handler
setMotorDir := make(<-chan int)  //Set motor direction til FSM, FSM calls function
floorEvent := make(chan<- int) //heis ved etasje fra elev_driver
buttonEvent := make(chan<- OrderButton) //knappetrykk fra elev_driver
newOrder := make(<- chan Order) //ekstern ordre til order_distribution for kost funksjonen
assignedNewOrder := make(chan <- Order) //ekstern ordre fra order_dist. , med heisID
sendRemoveOrder := make(<- chan Order) //ekstern ordre til network som skal fjernes
receiveRemoveOrder := make(chan <- Order) //ekstern order fra network som skal fjernes
sendNewOrder := make(<- chan Order) //ekstern ordre til network som skal legges til
receiveNewOrder := make(chan <- Order) //ekstern ordre fra network som skal legges til

//Network
peers := make(chan<- int) //status of peer to peer til order_dist.
sendCostValue := make(chan<- int) //kostverdi fra order_dist. for broadcasting
receiveCostValue := make(<-chan int) //kostverdi til order_dist fra broadcast
sendRemoveOrder := make(chan<- Order) //ekstern ordre fra orderhandler som skal broadcastes som remove order
receiveRemoveOrder := make(<-chan Order) //ekstern order til orderhandler  som skal fjernes
sendNewOrder := make(chan <- Order) //ekstern ordre fra orderhandler som skal broadcastes
receiveNewOrder := make(<- chan Order) //ekstern ordre til orderhandler som skal legges til


*/
