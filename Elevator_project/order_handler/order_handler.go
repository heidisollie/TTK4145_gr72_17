package order_handler

import (
	"fmt"
	."../structs"

)


var OrderQueue := []Order



func check_order(OrderQueue){
	Floor = ElevState.LastPassedFloor
	Direction = ElevState.CurrentDirection
	for _, order := range OrderQueue {
		if (order.Floor == Floor && order.ButtonType == Direction) || order.ButtonType == 0 {
			FSM.FSM_state = idle //Stop and open door
		}
	}
}

func is_duplicate(Order, OrderQueue) bool{
	for _, order := range OrderQueue {
		if Order == order {
			return true
		}
	
	}
	return false
}

func add_order(Order, OrderQueue){
	OrderQueue = append(OrderQueue, Order)

}

func remove_order(Order, OrderQueue, index){
	OrderQueue = append(OrderQueue[:index], OrderQueue[index+1]...)
}

func to_remove(ElevState){
	Floor = ElevState.LastPassedFloor
	for index, order := range OrderQueue {
		if order == Order {
			remove_order(Order, OrderQueue, index)
		}
	}
}

func get_order(){
	//Hva skal her	
}


func init(){
	
	





}


