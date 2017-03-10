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


	new_target_floor := make(chan int, 1)
	floor_event := make(chan int)                      //heis ved etasje til order_handler
	button_event := make(chan driver.OrderButton, 100) //knappetrykk til order_handler
	//peers := make(chan int)                        //status of peer to peer fra network
	new_order := make(chan structs.Order, 100)          //ekstern ordre fra order handler for kost funksjonen
	assigned_new_order := make(chan structs.Order, 100) //ekstern ordre fra order_dist. , med heisID
	elev_send_cost_value := make(chan structs.Cost, 100)
	elev_send_new_order := make(chan structs.Order, 100)
	elev_send_remove_order := make(chan structs.Order, 100)
	elev_receive_cost_value := make(chan structs.Cost, 100)
	elev_receive_new_order := make(chan structs.Order, 100)
	elev_receive_remove_order := make(chan structs.Order, 100)
	floor_completed := make(chan int, 100)
	number_of_peers := make(chan int, 100)
	//other_others_in_dir := make(chan bool, 100)

	driver.ElevInit()
	go driver.EventListener(button_event, floor_event)
	go FSM.FSM_init(floor_event, new_target_floor, floor_completed)

	go order_handler.Order_handler_init(localIP,
		floor_completed,
		button_event,
		assigned_new_order,
		new_order,
		elev_send_new_order,
		elev_send_remove_order,
		elev_receive_new_order,
		elev_receive_remove_order,
		new_target_floor)

	go network.UDP_init(localIP,
		elev_receive_cost_value,
		elev_receive_new_order,
		elev_receive_remove_order,
		elev_send_cost_value,
		elev_send_new_order,
		elev_send_remove_order,
		number_of_peers)

	go order_distribution.Order_dist_init(localIP,
		new_order,
		assigned_new_order,
		elev_receive_cost_value,
		elev_send_cost_value,
		number_of_peers)

	select {}
}
