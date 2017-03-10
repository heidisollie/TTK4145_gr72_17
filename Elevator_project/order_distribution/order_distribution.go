package order_distribution

import (
	"fmt"
	"../structs"
	//"strconv"
	"../localState"
	"strings"
	"sort"
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



//Sends order to order_handler of our cost is the lowest. If tie, picks the one with lower IP
func action_select(assignedNewOrder chan<- structs.Order, number_of_peers int, current_new_order structs.Order, localIP string, State_controller [string]structs.Elev_state) {
	
	const n = number_of_peers

	cost_list := [n]structs.Cost
	i := 0
	for index, state := range State_controller {
		cost_value = cost_function(current_new_order, state)
		current_new_order.IP = index
		cost := Cost{cost_value, current_new_order}
		cost_list[i] = cost_function(current_new_order, state)
		i++
	}

	smallestcost = 9999999
	for index, cost := range cost_list {
		if cost_list[index] < smallestcost {
				smallestcost = cost_list[index]
		}

	} 

	sort.Ints(cost_list)


	// IP ADDRESS OF COST AND ORDER ARE DIFFERENT
	fmt.Printf("The not-lowest score is:  %d \n", cost_list[1].Cost_value)
	fmt.Printf("The IP address of not-lowest cost value: %s \n", cost_list[1].Current_order.IP)

	fmt.Printf("The lowest score is:  %d \n", lowestScore.Cost_value)
	fmt.Printf("The IP address of lowest cost value: %s \n", lowestScore.Current_order.IP)
	
	//We want to change the evaluated order to contain the IP adress of the winning elevator

	assignedNewOrder <- lowestScore.Current_order

}



func Order_dist_init(localIP string,
	new_order <-chan structs.Order,
	assigned_new_order chan<- structs.Order,
	elev_receive_state <-chan structs.Cost,
	number_of_peers <-chan int) {

	var peers int
	var State_controller := make([string]structs.Elev_state)


	for {
		select {

		case state_update := elev_receive_state:
			State_controller[state_update.IP] = state_update

		case i := <- number_of_peers:
				peers = i
		case current_new_order := <-new_order:

			action_select(assigned_new_order, peers, current_new_order, localIP, State_controller)
		}
	}
}





