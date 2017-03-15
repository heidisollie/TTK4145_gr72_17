package orderHandler

import (
	"../backup"
	"../driver"
	"../localElevator"
	def "../definitions"
	"fmt"
	"time"
)

func printOrderQueue(OrderQueue []def.Order) {
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

func printOrder(Order def.Order) {
	fmt.Printf("-------PRINTING ORDER------- \n")
	fmt.Printf("Button: %d, Floor: %d \n", Order.Type, Order.Floor+1)
	fmt.Printf("-----------------------------\n")
}

func otherOrdersInDir(OrderQueue []def.Order, newTargetFloor chan<- int) {
	var floorSignal = driver.GetFloorSignal()
	if len(OrderQueue) != 0 {
		if floorSignal != -1 {
			for _, order := range OrderQueue {
				if order.Floor == floorSignal && (int(order.Type) == int(localElevator.ReadLocalElevator().CurrentDirection)+1 || int(order.Type) == 1) {
					newTargetFloor <- order.Floor
				}
			}
		}
	}
}

func isDuplicate(order def.Order, OrderQueue []def.Order) bool {
	for _, order_iter := range OrderQueue {
		if order == order_iter {
			return true
		}
	}
	return false
}

func getNewOrder(OrderQueue []def.Order, newTargetFloor chan<- int, localIP string) {
	for index, order := range OrderQueue {
		if order.IP == localIP {
			fmt.Printf("Found new order \n")
			newTargetFloor <- OrderQueue[index].Floor
			fmt.Printf("Sent new order \n")
			break
		}
	}
}

func addOrder(order def.Order, OrderQueue []def.Order, newTargetFloor chan<- int, localIP string) []def.Order {
	if isDuplicate(order, OrderQueue) == false {
		OrderQueue = append(OrderQueue, order)
		getNewOrder(OrderQueue, newTargetFloor, localIP)
		driver.SetButtonLamp(order.Type, order.Floor, 1)
	}
	return OrderQueue
}

func removeOrder(order def.Order, OrderQueue []def.Order) []def.Order {
	for index, order_iter := range OrderQueue {
		if order_iter == order {
			OrderQueue = removeElementSlice(OrderQueue, index)
			driver.SetButtonLamp(order_iter.Type, order_iter.Floor, 0)
			return OrderQueue
		}
	}
	return nil
}

func removeElementSlice(slice []def.Order, index int) []def.Order {
	if len(slice) == 0 {

	} else if len(slice) == 1 {
		slice = []def.Order{}
	} else {
		if index == len(slice)-1 {
			slice = slice[:index]
		} else if index == 0 {
			slice = slice[index+1:]
		} else {
			slice = append(slice[:index], slice[index+1:]...)
		}
	}
	return slice
}

func removeAll(floor int, elevSendRemoveOrder chan<- def.Order, OrderQueue []def.Order) []def.Order {
	for _, order := range OrderQueue {
		if order.Floor == floor {
			OrderQueue = removeOrder(order, OrderQueue)
			if order.Type != 1 {
				elevSendRemoveOrder <- order
			}
		}
	}
	return OrderQueue
}

func OrderHandlerInit(localIP string,
	floorCompleted <-chan int,
	buttonEvent <-chan driver.OrderButton,
	assignedNewOrder <-chan def.Order,
	processNewOrder chan<- def.Order,
	elevSendNewOrder chan<- def.Order,
	elevSendRemoveOrder chan<- def.Order,
	elevReceiveNewOrder <-chan def.Order,
	elevReceiveRemoveOrder <-chan def.Order,
	elevLost <-chan string,
	newTargetFloor chan<- int,
	floorEvent <-chan int) {


	loggingPeriod := 10 * time.Millisecond
	loggingTicker := time.NewTicker(loggingPeriod)


	var OrderQueue []def.Order
	backup.ReadQueueFromFile(&OrderQueue, def.Filename)

	for _, order := range OrderQueue {
		driver.SetButtonLamp(order.Type, order.Floor, 1)
	}

	if len(OrderQueue) != 0 {
		getNewOrder(OrderQueue, newTargetFloor, localIP)
	}

	for {
		select {
		case floor := <-floorCompleted:
			OrderQueue = removeAll(floor, elevSendRemoveOrder, OrderQueue)
			if len(OrderQueue) != 0 {
				getNewOrder(OrderQueue, newTargetFloor, localIP)
			}
		case orderButton := <-buttonEvent:
			if orderButton.Type == driver.ButtonCallCommand {
				newIntOrder := def.Order{Type: orderButton.Type, Floor: orderButton.Floor, IP: localIP}
				OrderQueue = addOrder(newIntOrder, OrderQueue, newTargetFloor, localIP)

			} else {
				newExtOrder := def.Order{Type: orderButton.Type, Floor: orderButton.Floor, IP: localIP}
				elevSendNewOrder <- newExtOrder
				processNewOrder <- newExtOrder
			}
		case newOrder := <-assignedNewOrder:
			fmt.Printf("Order handler: Received new order from ord dist\n")
			OrderQueue = addOrder(newOrder, OrderQueue, newTargetFloor, localIP)

		case order := <-elevReceiveRemoveOrder:
			OrderQueue = removeOrder(order, OrderQueue)

		case newOrder := <-elevReceiveNewOrder:
			processNewOrder <- newOrder

		case IP := <-elevLost:
			for _, order := range OrderQueue {
				if order.IP == IP {
					processNewOrder <- order
				}
			}

		case <-loggingTicker.C:
			backup.WriteQueueToFile(OrderQueue, def.Filename)

		case <-floorEvent:
			otherOrdersInDir(OrderQueue, newTargetFloor)

		}
	}
}
