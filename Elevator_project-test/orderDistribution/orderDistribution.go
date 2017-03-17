package orderDistribution

import (
	def "../definitions"
	elev "../localElevator"
	"../network/peers"
	"fmt"
	"math/rand"
	"strings"
)

func costFunction(newOrder def.Order, state def.Elevator) int {
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
	return costValue
}

func splitIP(IP string) string {
	s := strings.Split(IP, ".")
	return s[3]
}

func mergeSort(list []def.Cost) []def.Cost {
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
	mergeSort(list[:left])
	mergeSort(list[left+1:])
	return list
}

func orderAuction(assignedNewOrder chan<- def.Order,
	numberOfPeers int,
	newOrder def.Order,
	localIP string,
	stateController map[string]*def.Elevator) {

	CostList := make([]def.Cost, numberOfPeers)
	var winner def.Cost
	i := 0
	for index, state := range stateController {
		tempOrder := def.Order{newOrder.Type, newOrder.Floor, index}
		costValue := costFunction(newOrder, *state)
		cost := def.Cost{costValue, tempOrder}
		CostList[i] = cost
		i += 1
	}

	CostList = mergeSort(CostList)
	winner = CostList[0]

	//tie breaker
	if len(CostList) > 1 {
		if CostList[0].CostValue == CostList[1].CostValue {
			if splitIP(CostList[0].CurrentOrder.IP) > splitIP(CostList[1].CurrentOrder.IP) {
				winner = CostList[1]
			}
		}
	}
	assignedNewOrder <- winner.CurrentOrder
}

func OrderDistInit(localIP string,
	processNewOrder <-chan def.Order,
	assignedNewOrder chan<- def.Order,
	elevReceiveState <-chan def.Elevator,
	elevatorLost chan<- string,
	Peers <-chan peers.PeerUpdate) {

	var peers int = 0
	//Inkludere en bool i Elevator som er Online yes/no.
	stateController := make(map[string]*def.Elevator)
	stateController[localIP] = &def.Elevator{elev.ReadLocalElevator().LastPassedFloor, 0, false, localIP, true}

	for {
		select {
		case stateUpdate := <-elevReceiveState:
			if stateUpdate.Stuck != true {
				stateController[stateUpdate.IP] = &stateUpdate

			} else {
				delete(stateController, stateUpdate.IP)
				elevatorLost <- stateUpdate.IP
			}

			//Her kommer informasjonen om nye og lost peers
		case p := <-Peers:
			peers = len(p.Peers)
			if len(p.Lost) != 0 {
				if p.Lost[0] == localIP {
					fmt.Printf("----Lost network----\n")
					elev.ChangeLocalElevator_Online(false)
				} else {
					delete(stateController, p.Lost[0])
					elevatorLost <- p.Lost[0]
				}
			} else if len(p.New) != 0 && p.New == localIP {
				fmt.Printf("We're back online\n")
				elev.ChangeLocalElevator_Online(true)
			}

			//Over blir feil, pga vi sletter oss selv fra en tom liste!

			//Logikk som setter localElevator til false om ikke på nett og som setter den til true om den kommer på

		case order := <-processNewOrder:
			orderAuction(assignedNewOrder, peers, order, localIP, stateController)

		}
	}
}
