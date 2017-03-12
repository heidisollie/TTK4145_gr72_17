package order_distribution

import (
	"fmt"
	"../structs"
	//"strconv"
	//"../localState"
	"strings"

)


func cost_function(new_order structs.Order, state structs.Elev_state) int {
	cost_value := 0

	diff := new_order.Floor - state.Last_passed_floor


	//Turn reward
	if (state.Current_direction == 1 && diff > 0) || (state.Current_direction == -1 && diff < 0) {
		cost_value -= 50
	//Turn penalty
	} else if (state.Current_direction == 1 && diff < 0) || (state.Current_direction == -1 && diff > 0) {
		cost_value += 225

	//Distance
	} else if (state.Current_direction == 0) {
		if (diff == -1) || (diff == 1) {
			cost_value += 25		
		} else if (diff == -2) || (diff == 2) {
			cost_value += 50
		} else if (diff == -3) || (diff == 3) {
			cost_value += 75
		}
	} 

	return cost_value
}




func SplitIP(IP string) string {
	s := strings.Split(IP, ".")
	return s[3]
}



//Sends order to order_handler of with IP address of lowest cost. If tie, picks the one with lower IP
func action_select(assignedNewOrder chan<- structs.Order, 
	number_of_peers int, current_new_order structs.Order, 
	localIP string, 
	State_controller map[string]structs.Elev_state) {

	CostList := make([]structs.Cost, number_of_peers)
	var winner structs.Cost
	i := 0


	// This is dependent on the fact that disconnected elevators are not in State_Controller
	for index, state := range State_controller {
		current_new_order.IP = index //change the IP so it matches the cost_value
		//This is so when we assign the order, the IP will match the winner
		cost_value := cost_function(current_new_order, state)
		cost := structs.Cost{cost_value, current_new_order}
		CostList[i] = cost
		i++
	}

	//Sort CostList via insertion sort
	for i := 0; i<len(CostList); i++ {
		j := i
		for j>0 && CostList[j-1].Cost_value > CostList[j].Cost_value {
			temp := CostList[j-1]
			CostList [j-1] = CostList[j]
			CostList[j] = temp
			j -= 1		
		}
   	}
   	
   	winner = CostList[0]

   	//Check for tie, lowest IP wins
   	if CostList[0] == CostList[1] {
   		if SplitIP(CostList[0].Current_order.IP) > SplitIP(CostList[1].Current_order.IP) {
   			winner = CostList[1]
   		}
   	}




	// Debugging
	fmt.Printf("The not-lowest score is:  %d \n", CostList[1].Cost_value)
	fmt.Printf("The IP address of not-lowest cost value: %s \n", CostList[1].Current_order.IP)

	fmt.Printf("The lowest score is:  %d \n", CostList[0].Cost_value)
	fmt.Printf("The IP address of lowest cost value: %s \n", CostList[0].Current_order.IP)
	
	//We want to change the evaluated order to contain the IP adress of the winning elevator
	assignedNewOrder <- winner.Current_order

}



func Order_dist_init(localIP string,
	new_order <-chan structs.Order,
	assigned_new_order chan<- structs.Order,
	elev_receive_state <-chan structs.Elev_state,
	number_of_peers <-chan int) {

	var peers int = 0
	
	State_controller := map[string]structs.Elev_state{}

	for {
		select {

		case state_update := <- elev_receive_state:
			State_controller[state_update.IP] = state_update
		//Update number of peers if change
		case i := <- number_of_peers:
				peers = i
		case current_new_order := <-new_order:
			action_select(assigned_new_order, peers, current_new_order, localIP, State_controller)
		}
	}
}


