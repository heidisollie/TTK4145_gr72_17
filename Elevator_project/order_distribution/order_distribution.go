package order_distribution

import (
	"fmt"
	"../structs"
	//"strconv"
	"../localState"
	"strings"
	 "math/rand"

)


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
	} else if (state.CurrentDirection == 0) {
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


func sortListSort(list []structs.Cost) []structs.Cost {
	fmt.Printf("Maybe this is the problem \n")
    if len(list) < 2 {
        return list
    }

    left, right := 0, len(list)-1

    // Pick a pivot
    pivotIndex := rand.Int() % len(list)
    fmt.Printf("Hello %d \n", pivotIndex)
    // Move the pivot to the right
    list[pivotIndex], list[right] = list[right], list[pivotIndex]

    // Pile elements smaller than the pivot on the left
    fmt.Printf("Must be here somewhere\n")
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



//Sends order to order handler of with IP address of lowest cost. If tie, picks the one with lower IP
func actionSelect(assignedNewOrder chan<- structs.Order, 
	numberOfPeers int, 
	currentNewOrder structs.Order, 
	localIP string, 
	stateController map[string]structs.ElevState) {

	CostList := make([]structs.Cost, numberOfPeers)
	var winner structs.Cost
	i := 0


	// This is dependent on the fact that disconnected elevators are not in State_Controller
	//which is the case when starting up with one elevator, but we have not implemented something the removes disconnected elevators
	for index, state := range stateController {
		fmt.Printf("Iterating over stateController \n")
		currentNewOrder.IP = index //change the IP so it matches the costValue
		//This is so when we assign the order, the IP will match the winner
		costValue := costFunction(currentNewOrder, state)
		cost := structs.Cost{costValue, currentNewOrder}
		CostList[i] = cost
		i += 1
	}

/*
	//Sort CostList via insertion sort
	for i := 0; i<len(CostList); i++ {
		fmt.Printf("Sorting list \n")
		j := i
		for j>0 && CostList[j-1].CostValue > CostList[j].CostValue {
			temp := CostList[j-1]
			CostList [j-1] = CostList[j]
			CostList[j] = temp
			j -= 1		
		}
   	}
 */
   	CostList = sortListSort(CostList)	
   	winner = CostList[0]
   	//Check for tie, lowest IP wins
   	if len(CostList) > 1 {
	   	if CostList[0] == CostList[1] {
	   		if SplitIP(CostList[0].CurrentOrder.IP) > SplitIP(CostList[1].CurrentOrder.IP) {
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



func OrderDistInit(localIP string,
	newOrder <-chan structs.Order,
	assignedNewOrder chan<- structs.Order,
	elevReceiveState <-chan structs.ElevState,
	numberOfPeers <-chan int) {

	var peers int = 0
	
	stateController := map[string]structs.ElevState{}

	stateController[localIP] = structs.ElevState{localState.ReadLocalState().LastPassedFloor, 0, localIP}


	for {
		select {

		//Received updated state from other elevator
		case stateUpdate := <- elevReceiveState:
			fmt.Printf("Received updated state from %s \n", stateUpdate.IP)
			fmt.Print("State: dir: %d, flr: %d \n", stateUpdate.CurrentDirection, stateUpdate.LastPassedFloor+1)
			stateController[stateUpdate.IP] = stateUpdate

		//Update number of peers if change
		case i := <- numberOfPeers:
			fmt.Printf("Received updated number of peers: %d \n", i)
			peers = i

		//Received new order from order handler
		case order := <-newOrder:
			fmt.Printf("Received new external order - button: %d, floor: %d\n", order.Type, order.Floor + 1)
			actionSelect(assignedNewOrder, peers, order, localIP, stateController)
		}
	}
}


