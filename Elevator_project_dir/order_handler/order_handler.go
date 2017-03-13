package order_handler

import (
	"../driver"
	//. "../network"
	"../structs"
	"fmt"
	//"time"
	"../localState"
)

var OrderNil = structs.Order{0,0,false,"0"}

//for all orders in queue, sends new floor if order is command or matches direction
func other_orders_in_dir(OrderQueue []structs.Order, 
	new_current_order chan<- structs.Order, 
	localIP string) {
	var floorSignal = driver.GetFloorSignal()
	if floorSignal != -1 {
		for _, order := range OrderQueue {
			if order.Floor == floorSignal && (int(order.Type) == int(localState.ReadLocalState().Current_direction)+1 || int(order.Type) == 1) {
				new_current_order <- order
			} else if get_new_order(OrderQueue, new_current_order, localIP) != OrderNil {
				new_current_order <- order	
			}
		}
	}
}

func is_duplicate(order structs.Order, OrderQueue []structs.Order) bool {
	for _, order_iter := range OrderQueue {
		if order == order_iter {
			return true
		}
	}
	return false
}


func printOrderQueue(OrderQueue []structs.Order) {
	fmt.Printf("---------------------------------\n")
	fmt.Printf("PRINTING FROM FUNCTION\n")
	fmt.Printf("Order Queue: \n")
	for i, order := range OrderQueue {
		fmt.Printf("Element: %d \n", i+1)
		fmt.Print("Button  Floor\n")
		fmt.Print(order.Type, order.Floor+1)
		fmt.Printf("\n")
		fmt.Printf("IP: %s\n", order.IP)
	}
	fmt.Printf("Length: ")
	fmt.Print(len(OrderQueue))
	fmt.Printf("\n")
	fmt.Printf("---------------------------------\n")
}
//Finds first order in queue that is ours
func get_new_order(OrderQueue []structs.Order, new_current_order chan<- structs.Order, localIP string) structs.Order {
	for _, order := range OrderQueue {
		if order.IP == localIP {
			return order	
		}
	}

	return OrderNil
}

func add_order(order structs.Order, OrderQueue []structs.Order, new_current_order chan<- structs.Order, localIP string) []structs.Order {
	if is_duplicate(order, OrderQueue) == false {
		fmt.Printf("Adding new order\n")
		OrderQueue = append(OrderQueue, order)
		printOrderQueue(OrderQueue)
		
		//If added order is only order, (because if not we already have a current_order)	
		if len(OrderQueue) == 1 {
			order := get_new_order(OrderQueue, new_current_order, localIP)
			if order != OrderNil{
				new_current_order <- order
			}
		}
		driver.SetButtonLamp(order.Type, order.Floor, 1)
	}
	return OrderQueue
}

func remove_order(order structs.Order, OrderQueue []structs.Order) []structs.Order {
	for index, order_iter := range OrderQueue {
		if order_iter == order {
			OrderQueue = removeElementSlice(OrderQueue, index)
			//driver.SetButtonLamp(OrderQueue[index].Type, OrderQueue[index].Floor, 0)
			return OrderQueue
		}
	}
	return nil
}

func removeElementSlice(slice []structs.Order, index int) []structs.Order {
	//Hvis index er siste elementet
	if len(slice) == 0 {	
	} else if index == len(slice)-1 {

		if index == 0 {

			slice = []structs.Order{}
		} else {
		slice = slice[:index]
		}
	} else {
		slice = append(slice[:index], slice[index+1:]...)
	}

return slice
}

func remove_all(floor int, elev_send_remove_order chan<- structs.Order, OrderQueue []structs.Order) []structs.Order {
	for _, order := range OrderQueue {
		if order.Floor == floor {
			fmt.Printf("Order handler: Found order in floor, removing\n")
			fmt.Printf("Order handler: Floor: %d %s\n", order.Floor, order.Type)
			OrderQueue = remove_order(order, OrderQueue)
			driver.SetButtonLamp(order.Type, order.Floor, 0)
			printOrderQueue(OrderQueue)
			if order.Internal == false {
			elev_send_remove_order <- order
			}
		}
	}
	return OrderQueue
}

func Order_handler_init(localIP string,
	floor_completed <-chan int,
	button_event <-chan driver.OrderButton,
	assignedNewOrder <-chan structs.Order,
	newOrder chan<- structs.Order,
	elev_send_new_order chan<- structs.Order,
	elev_send_remove_order chan<- structs.Order,
	elev_receive_new_order <-chan structs.Order,
	elev_receive_remove_order <-chan structs.Order,
	new_current_order chan<- structs.Order) {

	var OrderQueue []structs.Order


	for {
		select {

		//From FSM, current order done, must send new
		case floor := <-floor_completed:
			fmt.Printf("Order handler: Floor completed message received\n")
			OrderQueue = remove_all(floor, elev_send_remove_order, OrderQueue)
			//If more orders, fetch them
			if len(OrderQueue) != 0 {
				printOrderQueue(OrderQueue)
				fmt.Printf("Order handler: Retrieving new order\n")
				get_new_order(OrderQueue, new_current_order, localIP)
			}

		case order_button := <-button_event:
			fmt.Printf("Received new button event\n")
			if order_button.Type == driver.ButtonCallCommand {
				//fmt.Printf("Order handler: Button pressed is command button\n")
				new_order := structs.Order{Type: order_button.Type, Floor: order_button.Floor, Internal: true, IP: localIP}
				OrderQueue = add_order(new_order, OrderQueue, new_current_order, localIP)

			} else { // if external, send to order_distribution
				new_order := structs.Order{Type: order_button.Type, Floor: order_button.Floor, Internal: false, IP: localIP}
				elev_send_new_order <- new_order // for å sende til network
				//fmt.Printf("Order handler: Sending new order to order_dist\n")
				newOrder <- new_order

			}

		//Received new order from order_distribution
		case new_order := <-assignedNewOrder:
			fmt.Printf("Order handler: Received new order from ord_dist\n")
			OrderQueue = add_order(new_order, OrderQueue, new_current_order, localIP)
			driver.SetButtonLamp( new_order.Type, new_order.Floor, 1)

		//Removing external order
		case order := <-elev_receive_remove_order:
			OrderQueue = remove_order(order, OrderQueue)
			printOrderQueue(OrderQueue)
			driver.SetButtonLamp(order.Type, order.Floor, 0)	

		//Adding external order
		case new_order := <-elev_receive_new_order:
			fmt.Printf("Order handler: Adding external order to our queue\n")
			newOrder <- new_order

		//Checking for other orders
		default:
			other_orders_in_dir(OrderQueue, new_current_order, localIP)
		}

	}
}