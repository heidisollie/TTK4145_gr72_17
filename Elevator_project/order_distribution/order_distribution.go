package order_distribution

import (
	"fmt"
	"../structs"
	//"strconv"
	"../localState"
	"strings"
)


func cost_function(new_order structs.Order) int {
	cost_value := 0
	diff := new_order.Floor - localState.ReadLocalState().Last_passed_floor

	//Turn reward
	if (localState.ReadLocalState().Current_direction == 1 && diff > 0) || (localState.ReadLocalState().Current_direction == -1 && diff < 0) {
		cost_value -= 50
	//Turn penalty
	} else if (localState.ReadLocalState().Current_direction == 1 && diff < 0) || (localState.ReadLocalState().Current_direction == -1 && diff > 0) {
		cost_value += 225
	//Distance
	} else if (localState.ReadLocalState().Current_direction == 0) {
		if (diff == -1) || (diff == 1) {
			cost_value += 25		
		} else if (diff == -2) || (diff == 2) {
			cost_value += 50
		} else if (diff == -3) || (diff == 3) {
			cost_value += 75
		}
	}
	fmt.Print(localState.ReadLocalState().Last_passed_floor)
	
	fmt.Printf("Order dist: State: %d, %d \n", localState.ReadLocalState().Last_passed_floor, localState.ReadLocalState().Current_direction)

	return cost_value
}

func findSmallestScore(list [structs.NUMELEV]structs.Cost) structs.Cost {
	for i:=0; i<len(list)-1; i++ {
		if list[i+1].Cost_value < list[i].Cost_value {
			temp := list[i]
			list[i] = list[i+1]
			list [i+1] = temp
		}
		if list[0] == list[1] {
			if (SplitIP(list[0].IP) < SplitIP(list[1].IP)){
				return list[0]
			} else {
				return list[1]
			}
		}
	} 
	return list[0]
}




func SplitIP(IP string) string {
	s := strings.Split(IP, ".")
	return s[3]
}



//Sends order to order_handler of our cost is the lowest. If tie, picks the one with lower IP
func action_select(assignedNewOrder chan<- structs.Order, elev_receive_cost_value <-chan structs.Cost, number_of_peers int, cost structs.Cost, localIP string) {
	filler_cost := structs.Cost{99999, cost.Current_order, cost.IP}
	var cost_list [structs.NUMELEV]structs.Cost
	for i := 0; i < structs.NUMELEV; i++ {
		cost_list[i] = filler_cost
	}

	switch (number_of_peers){

	case 1:
		cost_list[0] = cost
	case 2:
		new_cost := <-elev_receive_cost_value
		cost_list[0] = cost
		cost_list[1] = new_cost
	case 3:
		new_cost := <-elev_receive_cost_value
		new_cost2 := <-elev_receive_cost_value
		cost_list[0] = cost
		cost_list[1] = new_cost
		cost_list[2] = new_cost2
	}


	lowestScore := findSmallestScore(cost_list)

	fmt.Printf("ORDER DIST: Cost of lowest score: %d ----------------------------------------------- \n", lowestScore.Cost_value)
	if lowestScore.IP == localIP {
		fmt.Printf("Order_dist: Assigning new order")
		assignedNewOrder <- cost.Current_order
	}
}


func Order_dist_init(localIP string,
	new_order <-chan structs.Order,
	assigned_new_order chan<- structs.Order,
	elev_receive_cost_value <-chan structs.Cost,
	elev_send_cost_value chan<- structs.Cost,
	number_of_peers <-chan int) {

	var peers int

	for {
		select {
		case i := <- number_of_peers:
				peers = i
		case current_new_order := <-new_order:
			fmt.Printf("Order dist: Received new order\n")
			//Kjører kost funksjon og avgjør om den får bestillingen
			current_cost := structs.Cost{Cost_value: cost_function(current_new_order), Current_order: current_new_order, IP: localIP}
			fmt.Printf("Order dist: Cost value: %d\n", current_cost.Cost_value)

			//elev_send_cost_value <- current_cost

			//Ta inn kost verdi fra nettverkskanal og kjøre action_select
			action_select(assigned_new_order, elev_receive_cost_value, peers, current_cost, localIP)

		}
	}
}





