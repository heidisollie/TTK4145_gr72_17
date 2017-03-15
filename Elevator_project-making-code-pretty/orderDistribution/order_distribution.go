package orderDistribution

import (
	"../localState"
	"../network/peers"
	def "../definitions"
	"fmt"
	"math/rand"
	"strings"
)

func costFunction(newOrder def.Order, state def.ElevState) int {
	costValue := 0

	diff := newOrder.Floor - state.LastPassedFloor

	//Turn reward
	if (state.CurrentDirection == 1 && diff > 0) || (state.CurrentDirection == -1 && diff < 0) {
		costValue -= 50
	//Turn penalty
	} else if (state.CurrentDirection == 1 && diff < 0) || (state.CurrentDirection == -1 && diff > 0) {
		costValue += 225

	//Distance penalty
	} else if state.CurrentDirection == 0 {
		if (diff == -1) || (diff == 1) {
			costValue += 25
		} else if (diff == -2) || (diff == 2) {
			costValue += 50
		} else if (diff == -3) || (diff == 3) {
			costValue += 75
		}
	}
	return costValue
}

func SplitIP(IP string) string {
	s := strings.Split(IP, ".")
	return s[3]
}

func sortListSort(list []def.Cost) []def.Cost {
	if len(list) < 2 {
		return list
	}
	left, right := 0, len(list)-1
	pivotIndex := rand.Int() % len(list)
	list[pivotIndex], list[right] = list[right], list[pivotIndex]
	for i := range list {
		if list[i].CostValue < list[right].CostValue {
			list[i], list[left] = list[left], list[i]
			left++
		}
	}
	list[left], list[right] = list[right], list[left]
	sortListSort(list[:left])
	sortListSort(list[left+1:])
	return list
}

func actionSelect(assignedNewOrder chan<- def.Order,
	numberOfPeers int,
	currentNewOrder def.Order,
	localIP string,
	stateController map[string]def.ElevState) {
	print("The new order is; %d ")
	CostList := make([]def.Cost, numberOfPeers)
	var winner def.Cost
	i := 0

	for index, state := range stateController {
		fmt.Printf("index: %s \n", index)
		currentNewOrder.IP = index

		costValue := costFunction(currentNewOrder, state)
		cost := def.Cost{costValue, currentNewOrder}
		CostList[i] = cost
		i += 1

	}

	CostList = sortListSort(CostList)
	for i := 0; i < len(CostList); i++ {
		fmt.Printf("Cost value of index %d \n", i)
		fmt.Printf("%d \n", CostList[i].CostValue)
	}
	winner = CostList[0]

	if len(CostList) > 1 {
		if CostList[0] == CostList[1] {
			if SplitIP(CostList[0].CurrentOrder.IP) > SplitIP(CostList[1].CurrentOrder.IP) {
				winner = CostList[1]
			}
		}
	}
	assignedNewOrder <- winner.CurrentOrder
}


func OrderDistInit(localIP string,
	processNewOrder <-chan def.Order,
	assignedNewOrder chan<- def.Order,
	elevReceiveState <-chan def.ElevState,
	elevLost chan<- string,
	Peers <-chan peers.PeerUpdate) {

	var peers int = 0

	stateController := map[string]def.ElevState{}
	stateController[localIP] = def.ElevState{localElevatr.ReadLocalElevator().LastPassedFloor, 0, false, localIP}

	for {
		select {
		case stateUpdate := <-elevReceiveState:
			if stateUpdate.Stuck != true {
				stateController[stateUpdate.IP] = stateUpdate
			} else {
				delete(stateController, stateUpdate.IP)
				elevLost <- stateUpdate.IP
			}
		case p := <-Peers:
			peers = len(p.Peers)
			if len(p.Lost) != 0 {
				delete(stateController, p.Lost[0])
				elevLost <- p.Lost[0]
			}
		case order := <-processNewOrder:
			actionSelect(assignedNewOrder, peers, order, localIP, stateController)
		}
	}
}
