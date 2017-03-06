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

func main() {
	fmt.Printf("Elev_driver started\n")
	driver.ElevInit()

	const NumFloors = driver.NumFloors
	const NumButtons = driver.NumButtons

	var localIP string

	/*
		localIP, err := localip.LocalIP()
		if err != nil {
			fmt.Println(err)
			localIP = "DISCONNECTED"
		}*/

	localIP = "3456"

	var OrderQueue []structs.Order

	new_target_floor := make(chan int, 1)
	floor_event := make(chan int)                      //heis ved etasje til order_handler
	button_event := make(chan driver.OrderButton, 100) //knappetrykk til order_handler
	//peers := make(chan int)                        //status of peer to peer fra network
	new_order := make(chan structs.Order)          //ekstern ordre fra order handler for kost funksjonen
	assigned_new_order := make(chan structs.Order) //ekstern ordre fra order_dist. , med heisID
	elev_send_cost_value := make(chan structs.Cost)
	elev_send_new_order := make(chan structs.Order)
	elev_send_remove_order := make(chan structs.Order)
	elev_receive_cost_value := make(chan structs.Cost)
	elev_receive_new_order := make(chan structs.Order)
	elev_receive_remove_order := make(chan structs.Order)
	floor_completed := make(chan int)

	State := structs.Elev_state{Last_passed_floor: 0, Current_direction: driver.DirnStop, IP: localIP}

	driver.ElevInit()
	go driver.EventListener(button_event, floor_event)
	go FSM.FSM_init(State, floor_event, new_target_floor, floor_completed)

	go order_handler.Order_handler_init(State,
		OrderQueue,
		localIP,
		floor_completed,
		button_event,
		assigned_new_order,
		new_order,
		elev_send_new_order,
		elev_send_remove_order,
		elev_receive_new_order,
		elev_receive_remove_order,
		new_target_floor)

	go network.UDP_init(elev_receive_cost_value,
		elev_receive_new_order,
		elev_receive_remove_order,
		elev_send_cost_value,
		elev_send_new_order,
		elev_send_remove_order)

	go order_distribution.Order_dist_init(localIP,
		new_order,
		assigned_new_order,
		elev_receive_cost_value,
		elev_send_cost_value)

	select {}
}
