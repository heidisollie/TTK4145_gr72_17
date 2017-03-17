package orderHandler

import (
	"../backup"
	def "../definitions"
	"../driver"
	elev "../localElevator"
	"fmt"
	"time"
)

func otherOrdersInDir(orderQueue []def.Order, newTargetFloor chan<- int) {
	var floorSignal = driver.GetFloorSignal()
	if len(orderQueue) != 0 && floorSignal != -1 {
		for _, order := range orderQueue {
			if order.Floor == floorSignal && (int(order.Type) == int(elev.ReadLocalElevator().CurrentDirection)+1 || int(order.Type) == 1) {
				newTargetFloor <- order.Floor
			}
		}
	}
}

func isDuplicate(order def.Order, orderQueue []def.Order) bool {
	for _, order_iter := range orderQueue {
		if order == order_iter {
			return true
		}
	}
	return false
}

//Sends first local order in queue to FSM
func sendNextOrder(orderQueue []def.Order, newTargetFloor chan<- int, localIP string) {
	for index, order := range orderQueue {
		if order.IP == localIP {
			newTargetFloor <- orderQueue[index].Floor
			break
		}
	}
}

func addOrder(order def.Order, orderQueue []def.Order, newTargetFloor chan<- int, localIP string) []def.Order {
	if isDuplicate(order, orderQueue) == false {
		orderQueue = append(orderQueue, order)
		fmt.Printf("Added order to queue\n")
		sendNextOrder(orderQueue, newTargetFloor, localIP)
		driver.SetButtonLamp(order.Type, order.Floor, 1)
	}
	return orderQueue
}

func removeOrder(order def.Order, orderQueue []def.Order) []def.Order {
	for index, order_iter := range orderQueue {
		if order_iter == order {
			orderQueue = sliceRemove(orderQueue, index)
			fmt.Printf("Removed order from queue \n")
			driver.SetButtonLamp(order_iter.Type, order_iter.Floor, 0)
			return orderQueue
		}
	}
	return nil
}

func sliceRemove(slice []def.Order, index int) []def.Order {
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

func removeOrdersAtFloor(floor int, elevSendRemoveOrder chan<- def.Order, orderQueue []def.Order) []def.Order {
	for _, order := range orderQueue {
		if order.Floor == floor {
			orderQueue = removeOrder(order, orderQueue)
			if order.Type != 1 {
				elevSendRemoveOrder <- order
			}
		}
	}
	return orderQueue
}

func OrderHandlerInit(localIP string,
	floorCompleted <-chan int,
	buttonEvent <-chan def.OrderButton,
	assignedNewOrder <-chan def.Order,
	processNewOrder chan<- def.Order,
	elevSendNewOrder chan<- def.Order,
	elevSendRemoveOrder chan<- def.Order,
	elevReceiveNewOrder <-chan def.Order,
	elevReceiveRemoveOrder <-chan def.Order,
	elevatorLost <-chan string,
	newTargetFloor chan<- int,
	floorEvent <-chan int) {

	var orderQueue []def.Order
	loggingPeriod := 200 * time.Millisecond
	loggingTicker := time.NewTicker(loggingPeriod)

	//Initilizing orderQueue from backup
	backup.ReadQueueFromFile(&orderQueue, def.Filename)
	for _, order := range orderQueue {
		driver.SetButtonLamp(order.Type, order.Floor, 1)
	}
	if len(orderQueue) != 0 {
		sendNextOrder(orderQueue, newTargetFloor, localIP)
	}

	for {
		select {
		case floor := <-floorCompleted:
			orderQueue = removeOrdersAtFloor(floor, elevSendRemoveOrder, orderQueue)
			if len(orderQueue) != 0 {
				sendNextOrder(orderQueue, newTargetFloor, localIP)
			}

		case orderButton := <-buttonEvent:
			fmt.Printf("Button was pressed\n")
			if orderButton.Type == def.ButtonCallCommand {
				fmt.Printf("1\n")
				newIntOrder := def.Order{Type: orderButton.Type, Floor: orderButton.Floor, IP: localIP}
				orderQueue = addOrder(newIntOrder, orderQueue, newTargetFloor, localIP)

			} else {
				fmt.Printf("2\n")
				if elev.ReadLocalElevator().Online == true {
					newExtOrder := def.Order{Type: orderButton.Type, Floor: orderButton.Floor, IP: localIP}
					elevSendNewOrder <- newExtOrder
					processNewOrder <- newExtOrder
				} else {
					fmt.Printf("Elevator has no nettwork\n")
				}

			}

		case order := <-assignedNewOrder:
			orderQueue = addOrder(order, orderQueue, newTargetFloor, localIP)
			driver.SetButtonLamp(order.Type, order.Floor, 1)

		case order := <-elevReceiveRemoveOrder:
			orderQueue = removeOrder(order, orderQueue)
			driver.SetButtonLamp(order.Type, order.Floor, 0)
		case order := <-elevReceiveNewOrder:
			processNewOrder <- order

			//Ha noe logikk her angående tap av nettverk?
		case IP := <-elevatorLost:
			for _, order := range orderQueue {
				if order.IP == IP {
					processNewOrder <- order
				}

			}

		case <-loggingTicker.C:
			backup.WriteQueueToFile(orderQueue, def.Filename)

		case <-floorEvent:
			otherOrdersInDir(orderQueue, newTargetFloor)
		}
	}
}
