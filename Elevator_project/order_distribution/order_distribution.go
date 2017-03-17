package order_distribution

import (
	"../localState"
	"../network/peers"
	"../structs"
	"fmt"
	"math/rand"
	"strings"
)

//Om heisen er stuck må kosten bli "uendelig"
//Sette last_passed_floor til 4? Er det mulig
//Det er "jallafix" (Y)

func costFunction(newOrder structs.Order, state structs.ElevState) int {
	costValue := 0

	diff := newOrder.Floor - state.LastPassedFloor

	//Turn reward
	if (state.CurrentDirection == 1 && diff > 0) || (state.CurrentDirection == -1 && diff < 0) {
		costValue -= 50
		//Turn penalty
	} else if (state.CurrentDirection == 1 && diff < 0) || (state.CurrentDirection == -1 && diff > 0) {
		costValue += 225

		//Distance
	} else if state.CurrentDirection == 0 {
		if (diff == -1) || (diff == 1) {
			costValue += 25
		} else if (diff == -2) || (diff == 2) {
			costValue += 50
		} else if (diff == -3) || (diff == 3) {
			costValue += 75
		}
	}
	/*
		if state.Stuck == true {
			costValue += 99999
		}*/

	return costValue
}

func SplitIP(IP string) string {
	s := strings.Split(IP, ".")
	return s[3]
}

func sortListSort(list []structs.Cost) []structs.Cost {

	if len(list) < 2 {
		return list
	}

	left, right := 0, len(list)-1

	// Pick a pivot
	pivotIndex := rand.Int() % len(list)

	// Move the pivot to the right
	list[pivotIndex], list[right] = list[right], list[pivotIndex]

	// Pile elements smaller than the pivot on the left

	for i := range list {
		if list[i].CostValue < list[right].CostValue {
			list[i], list[left] = list[left], list[i]
			left++
		}
	}

	// Place the pivot after the last smaller element
	list[left], list[right] = list[right], list[left]

	// Go down the rabbit hole
	sortListSort(list[:left])
	sortListSort(list[left+1:])

	return list
}

func printCostList(costList []structs.Cost) {
	fmt.Printf("------Printing cost list------\n")
	for i := 0; i < len(costList); i++ {
		fmt.Printf("Cost value of index %d \n", i)
		fmt.Printf("%d - IP: %s \n", costList[i].CostValue, costList[i].CurrentOrder.IP)
	}
	fmt.Printf("--------------------------------\n")
}

//Sends order to order handler of with IP address of lowest cost. If tie, picks the one with lower IP
func actionSelect(assignedNewOrder chan<- structs.Order,
	numberOfPeers int,
	newOrder structs.Order,
	localIP string,
	stateController map[string]*structs.ElevState) {

	fmt.Printf("The new order is: button: %d, floor: %d \n", newOrder.Type, newOrder.Floor+1)
	CostList := make([]structs.Cost, numberOfPeers)
	var winner structs.Cost
	i := 0

	for index, state := range stateController {
		fmt.Printf("Index: %s \n", index)
		fmt.Printf("State: flr: %d, dir: %d, IP: %s \n", state.LastPassedFloor, state.CurrentDirection, state.IP)
		tempOrder := structs.Order{newOrder.Type, newOrder.Floor, index}
		costValue := costFunction(newOrder, *state)
		cost := structs.Cost{costValue, tempOrder}
		CostList[i] = cost
		fmt.Printf("IP of order in cost value: %s \n", cost.CurrentOrder.IP)
		i += 1

	}

	CostList = sortListSort(CostList)
	printCostList(CostList)

	winner = CostList[0]
	//Check for tie, lowest IP wins
	if len(CostList) > 1 {
		if CostList[0].CostValue == CostList[1].CostValue {
			//fmt.Printf("last ip: %s, last ip: %s\n", SplitIP(CostList[0].CurrentOrder.IP), SplitIP(CostList[1].CurrentOrder.IP))
			if SplitIP(CostList[0].CurrentOrder.IP) > SplitIP(CostList[1].CurrentOrder.IP) {
				fmt.Printf("Switch winner because %s > %s \n", SplitIP(CostList[0].CurrentOrder.IP), SplitIP(CostList[1].CurrentOrder.IP))
				winner = CostList[1]
			}
		}
	}

	if len(CostList) > 1 {
		// Debugging
		fmt.Printf("The not-lowest score is:  %d \n", CostList[1].CostValue)
		fmt.Printf("The IP address of not-lowest cost value: %s \n", CostList[1].CurrentOrder.IP)
	}

	fmt.Printf("The lowest score is:  %d \n", CostList[0].CostValue)
	fmt.Printf("The IP address of lowest cost value: %s \n", CostList[0].CurrentOrder.IP)

	//We want to change the evaluated order to contain the IP adress of the winning elevator
	assignedNewOrder <- winner.CurrentOrder

}

func redistributionOfOrders() {

}

func OrderDistInit(localIP string,
	processNewOrder <-chan structs.Order,
	assignedNewOrder chan<- structs.Order,
	elevReceiveState <-chan structs.ElevState,
	elevLost chan<- string,
	Peers <-chan peers.PeerUpdate) {

	var peers int = 0
	stateController := make(map[string]*structs.ElevState)
	stateController[localIP] = &structs.ElevState{localState.ReadLocalState().LastPassedFloor, 0, false, localIP, true}

	for {
		select {

		//Received updated state from other elevator
		case stateUpdate := <-elevReceiveState:
			if stateUpdate.Stuck != true {
				stateController[stateUpdate.IP] = &stateUpdate
				//fmt.Printf("Received states are: flr: %d, dir: %d, IP: %s \n", stateUpdate.LastPassedFloor, stateUpdate.CurrentDirection, stateUpdate.IP)
			} else {
				//If elevator stuck
				delete(stateController, stateUpdate.IP)
				elevLost <- stateUpdate.IP
			}

		//Update number of peers if change and remove lost elevators som stateController
		case p := <-Peers:
			peers = len(p.Peers)
			if len(p.Lost) != 0 {
				if p.Lost[0] == localIP {
					fmt.Printf("----Lost network----\n")
					localState.ChangeLocalState_Online(false)
				} else {
					delete(stateController, p.Lost[0])
					elevLost <- p.Lost[0]
				}
			} else if len(p.New) != 0 && p.New == localIP {
				fmt.Printf("We're back online\n")
				localState.ChangeLocalState_Online(true)
			}
		//Received new order from order handler
		case order := <-processNewOrder:
			fmt.Printf("Received new external order - button: %d, floor: %d\n", order.Type, order.Floor+1)
			actionSelect(assignedNewOrder, peers, order, localIP, stateController)

		}
	}
}