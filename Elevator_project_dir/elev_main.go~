package main

import (
	"./order_distribution"
	"./FSM"
	"./order_handler"
	"./network"
	"./logger"
	"./driver"
	"fmt"
	"time"
)

const NumFloors = driver.NumFloors
const NumButtons = driver.NumButtons

var OrderQueue := []Order


//FSM
/*
floorEvent <- chan int, //heis ved etasje til order_handler
buttonEvent <- chan OrderButton //knappetrykk til order_handler
*/

//Elev_driver
floor_event := make(chan int) //heis ved etasje til order_handler
button_event := make(chan OrderButton) //knappetrykk til order_handler

//Order_dist.
peers := make(chan int) //status of peer to peer fra network
newOrder := make(chan Order) //ekstern ordre fra order handler for kost funksjonen
assignedNewOrder := make(chan Order) //ekstern ordre til order handler, med heisID

//Order_handler
assignedNewOrder := make(chan Order) //ekstern ordre fra order_dist. , med heisID


elev_send_cost_value = make(chan Cost)
elev_send_new_order = make(chan Order)
elev_send_remove_order = make(chan Order)
elev_receive_cost_value = make(chan int)
elev_receive_new_order = make(chan Order)
elev_receive_remove_order = make (chan Order)



//Network
peers := make(chan int) //status of peer to peer til order_dist.




func main(){
	fmt.Printf("Elev_driver started\n")
	driver.ElevInit()

	//Initialisere
	go FSM.init()
	go order_handler.init()
	go network.UPD_init(elev_send_cost_value,
					elev_send_new_order,
					elev_send_remove_order,
					elev_receive_cost_value,
					elev_receive_new_order,
					elev_receive_remove_order)
	
	
	buttonEvent := make(chan driver.OrderButton)
	floorEvent := make(chan int)
	for {
		ButtonLights()
	}
	go driver.EventListener(buttonEvent, floorEvent) 
	
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
