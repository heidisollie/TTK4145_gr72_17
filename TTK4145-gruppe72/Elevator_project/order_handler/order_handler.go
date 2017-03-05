package order_handler

import (
	"../driver"
	. "../structs"
	."../network"
)

//var OrderQueue []Order

localIP, err := localip.LocalIP()
if err != nil {
	fmt.Println(err)
	localIP = "DISCONNECTED"
}



type Queue struct {
	list []Order
}

var OrderQueue Queue

func other_orders_in_dir(OrderQueue Queue, new_target_floor chan int) {
	Floor := Elev_state.last_passed_floor
	Direction := Elev_state.current_direction
	for _, order := range OrderQueue.list {
		if (order.Floor == Floor && order.Type == Direction) || order.Type == 0 {
			new_target_floor <- Floor
		}
	}
}

func is_duplicate(order Order, OrderQueue Queue) bool {
	for _, order_iter := range OrderQueue.list {
		if order == order_iter {
			return true
		}

	}
	return false
}

func get_new_order(OrderQueue Queue, new_target_floor chan<- int) {
	for index, order := range OrderQueue.list {
		if order.IP == localIP {
			new_target_floor <- OrderQueue.list[0].Floor
			break
		}
	}
}

func add_order(order Order, OrderQueue Queue) {
	if !is_duplicate(order, OrderQueue) == false {
		OrderQueue.list = append(OrderQueue.list, order)
		if len(OrderQueue.list) == 1 {
			get_new_order(OrderQueue)
		}
		driver.SetButtonLamp(Order.Type, Order.Floor, 1)
	}
}

func remove_order(order Order, OrderQueue Queue) {
	for index, order_iter := range OrderQueue.list {
		if order_iter == order {
			OrderQueue.list = append(OrderQueue.list[:index], OrderQueue.list[index+1]...)
			driver.SetButtonLamp(OrderQueue.list[index].Type, OrderQueue.list[index].Floor, 0)
			break
		}
	}
}

func remove_all(floor int, elev_send_remove_order chan Order) {
	for _, order := range OrderQueue.list {
		if order.Floor == floor {
			remove_order(order, OrderQueue)
			elev_send_remove_order <- order
		}
	}
}

/*
func merge_orders(){
	//Ved nettverksbrud, merger ordre med våre
	//Kanskje ikke nødvendig da alle ekstern bestillingen allerede er i køen vår
}*/

func order_handler_init(
	floor_completed <-chan int,
	button_event <-chan driver.OrderButton,
	assignedNewOrder <-chan Order,
	newOrder chan<- Order,
	elev_send_new_order chan<- Order,
	elev_send_remove_order chan<- Order,
	elev_receive_new_order <-chan Order,
	elev_receive_remove_order <-chan Order,
	new_target_floor chan<- int) {

	for {
		other_orders_in_dir(OrderQueue, new_target_floor)
		select {
		case floor := <-floor_completed:
			remove_all(floor, elev_send_remove_order)
			if len(OrderQueue.list) != 0 {
				get_new_order(OrderQueue, new_target_floor)
			}
		case order_button := <-button_event:
			if order_button.Type == ButtonCallCommand {
				new_order = Order{Type: order_button.Type, Floor: order_button.Floor, Internal: true, IP: localIP}
				add_order(new_order, OrderQueue)
				elev_send_new_order <- order
			} else {
				new_order = Order{Type: order_button.Type, Floor: order_button.Floor, Internal: false, IP: localIP}
				newOrder <- new_order
			}
		case new_order := <-assignedNewOrder:
			add_order(new_order, OrderQueue)
			elev_send_new_order <- order

		case remove_order := <-elev_receive_remove_order:
			remove_order(order, OrderQueue)

		case new_order := <-elev_receive_new_order:
			add_order(new_order, OrderQueue)
		}
	}
}
