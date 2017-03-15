package order_handler

import (
	//	"../backup"
	"../backup"
	"../driver"
	"../localState"
	"../structs"
	"fmt"
	"time"
)

//for all orders in queue, sends new floor if order is command or matches direction
func otherOrdersInDir(OrderQueue []structs.Order, newTargetFloor chan<- int) {
	var floorSignal = driver.GetFloorSignal()
	if len(OrderQueue) != 0 {
		if floorSignal != -1 {
			for _, order := range OrderQueue {
				if order.Floor == floorSignal && (int(order.Type) == int(localState.ReadLocalState().CurrentDirection)+1 || int(order.Type) == 1) {
					fmt.Printf("Found extra order \n")
					newTargetFloor <- order.Floor
				}
			}
		}
	}
}

func isDuplicate(order structs.Order, OrderQueue []structs.Order) bool {
	for _, order_iter := range OrderQueue {
		if order == order_iter {
			return true
		}
	}
	return false
}

func printOrderQueue(OrderQueue []structs.Order) {
	fmt.Printf("-------PRINTING QUEUE------- \n")
	fmt.Printf("Order Queue: \n")
	for i, order := range OrderQueue {
		fmt.Printf("Element: %d \n", i+1)
		fmt.Print("Button  Floor\n")
		fmt.Printf("%d        %d \n", order.Type, order.Floor+1)
		fmt.Printf("\n")
		fmt.Printf("IP: %s\n", order.IP)
	}
	fmt.Printf("Length: ")
	fmt.Print(len(OrderQueue))
	fmt.Printf("\n")
	fmt.Printf("---------------------------------\n")
}

func printOrder(Order structs.Order) {
	//fmt.Printf("-------PRINTING ORDER------- \n")
	//fmt.Printf("Button: %d, Floor: %d \n", Order.Type, Order.Floor+1)
	//fmt.Printf("-----------------------------\n")
}

func getNewOrder(OrderQueue []structs.Order, newTargetFloor chan<- int, localIP string) {
	for index, order := range OrderQueue {
		if order.IP == localIP {
			//fmt.Printf("Found new order \n")
			newTargetFloor <- OrderQueue[index].Floor
			//fmt.Printf("Sent new order \n")
			break
		}
	}
}

func addOrder(order structs.Order, OrderQueue []structs.Order, newTargetFloor chan<- int, localIP string) []structs.Order {
	if isDuplicate(order, OrderQueue) == false {
		//fmt.Printf("Adding new order\n")
		OrderQueue = append(OrderQueue, order)
		printOrderQueue(OrderQueue)
		getNewOrder(OrderQueue, newTargetFloor, localIP)
		driver.SetButtonLamp(order.Type, order.Floor, 1)
	}
	return OrderQueue
}

func removeOrder(order structs.Order, OrderQueue []structs.Order) []structs.Order {
	for index, order_iter := range OrderQueue {
		if order_iter == order {
			OrderQueue = removeElementSlice(OrderQueue, index)
			driver.SetButtonLamp(order_iter.Type, order_iter.Floor, 0)
			return OrderQueue
		}
	}
	return nil
}

func removeElementSlice(slice []structs.Order, index int) []structs.Order {
	if len(slice) == 0 {

	} else if len(slice) == 1 {
		slice = []structs.Order{}
	} else {
		if index == len(slice)-1 { //if removed order is last element
			slice = slice[:index] //includes everything up to current index (since it is last)
		} else if index == 0 {
			slice = slice[index+1:] //includes everything except for current index (since it is first)
		} else {
			slice = append(slice[:index], slice[index+1:]...)
		}
	}

	return slice
}

func removeAll(floor int, elevSendRemoveOrder chan<- structs.Order, OrderQueue []structs.Order) []structs.Order {
	for _, order := range OrderQueue {
		if order.Floor == floor {
			OrderQueue = removeOrder(order, OrderQueue)

			if order.Type != 1 {
				elevSendRemoveOrder <- order
			}
		}
	}
	//fmt.Printf("After remove all \n")
	printOrderQueue(OrderQueue)
	return OrderQueue
}

func OrderHandlerInit(localIP string,
	floorCompleted <-chan int,
	buttonEvent <-chan driver.OrderButton,
	assignedNewOrder <-chan structs.Order,
	processNewOrder chan<- structs.Order,
	elevSendNewOrder chan<- structs.Order,
	elevSendRemoveOrder chan<- structs.Order,
	elevReceiveNewOrder <-chan structs.Order,
	elevReceiveRemoveOrder <-chan structs.Order,
	elevLost <-chan string,
	newTargetFloor chan<- int,
	floorEvent <-chan int) {

	var OrderQueue []structs.Order

	backup.ReadQueueFromFile(&OrderQueue, structs.Filename)

	for _, order := range OrderQueue {
		driver.SetButtonLamp(order.Type, order.Floor, 1)
	}
	fmt.Printf("Initial queue \n")
	printOrderQueue(OrderQueue)

	if len(OrderQueue) != 0 {
		getNewOrder(OrderQueue, newTargetFloor, localIP)
	}

	loggingPeriod := 100 * time.Millisecond
	loggingTicker := time.NewTicker(loggingPeriod)

	for {
		select {
		case floor := <-floorCompleted:
			//fmt.Printf("Order handler: Floor completed message received\n")
			OrderQueue = removeAll(floor, elevSendRemoveOrder, OrderQueue)
			////fmt.Printf("Order handler: Removed from order queue\n")
			if len(OrderQueue) != 0 {
				printOrderQueue(OrderQueue)
				//fmt.Printf("Order handler: Retrieving new order\n")
				getNewOrder(OrderQueue, newTargetFloor, localIP)
			}

		case orderButton := <-buttonEvent:
			if orderButton.Type == driver.ButtonCallCommand {
				////fmt.Printf("Order handler: Button pressed is command button\n")
				newIntOrder := structs.Order{Type: orderButton.Type, Floor: orderButton.Floor, IP: localIP}
				OrderQueue = addOrder(newIntOrder, OrderQueue, newTargetFloor, localIP)

			} else { // if external, send to order distribution
				newExtOrder := structs.Order{Type: orderButton.Type, Floor: orderButton.Floor, IP: localIP}
				elevSendNewOrder <- newExtOrder // for å sende til network
				////fmt.Printf("Order handler: Sending new order to order_dist\n")
				processNewOrder <- newExtOrder
			}

		case newOrder := <-assignedNewOrder:
			//fmt.Printf("Order handler: Received new order from ord dist\n")
			OrderQueue = addOrder(newOrder, OrderQueue, newTargetFloor, localIP)
			driver.SetButtonLamp(newOrder.Type, newOrder.Floor, 1)
			////fmt.Printf("Order handler: Added new order in Order Queue\n")

		case order := <-elevReceiveRemoveOrder:
			OrderQueue = removeOrder(order, OrderQueue)
			printOrderQueue(OrderQueue)
			driver.SetButtonLamp(order.Type, order.Floor, 0)
		case newOrder := <-elevReceiveNewOrder:

			//fmt.Printf("Order handler: Receiveing external order, sendt to order dist for processing\n")
			processNewOrder <- newOrder

		// Redistributing orders for lost elevator
		// Note: if elevator comes back on the network, nothing happens.
		// The elevators are just ineffective for a short period of time
		case IP := <-elevLost:
			for _, order := range OrderQueue {
				if order.IP == IP {
					processNewOrder <- order
				}

			}

		case <-loggingTicker.C:
			backup.WriteQueueToFile(OrderQueue, structs.Filename)

		case <-floorEvent:
			otherOrdersInDir(OrderQueue, newTargetFloor)
		}
	}
}
