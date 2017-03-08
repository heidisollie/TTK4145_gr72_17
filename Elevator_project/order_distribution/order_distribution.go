package order_distribution

import (
	"fmt"
	"../structs"
	//"strconv"
	"../localState"
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
	
	fmt.Printf("State: %d, %d \n", localState.ReadLocalState().Last_passed_floor, localState.ReadLocalState().Current_direction)
	return cost_value
}

//Sends order to order_handler of our cost is the lowest. If tie, picks the one with lower IP
func action_select(assignedNewOrder chan<- structs.Order, elev_receive_cost_value <-chan structs.Cost, cost structs.Cost) {
	//cost_value1 := <-elev_receive_cost_value
	//cost_value2 := <-elev_receive_cost_value
	cost_value1 := 10000
	cost_value2 := 10000
	//i, _ := strconv.Atoi(cost.IP)
	//a, _ := strconv.Atoi(cost_value1.IP)
	//b, _ := strconv.Atoi(cost_value2.IP)
	if cost.Cost_value < cost_value1 && cost.Cost_value < cost_value2 {
		assignedNewOrder <- cost.Current_order
		fmt.Printf("Sending new order to order_handler\n")
	} 
	/*else if cost.Cost_value == cost_value1.Cost_value || cost.Cost_value == cost_value2.Cost_value {
		if cost.Cost_value == cost_value1.Cost_value && cost.Cost_value < cost_value2.Cost_value {

			if i < a {
				assignedNewOrder <- cost.Current_order
			}
		} else if cost.Cost_value == cost_value1.Cost_value && cost.Cost_value < cost_value2.Cost_value {
			if i < b {
				assignedNewOrder <- cost.Current_order
			}	
		}
	}*/
}


func Order_dist_init(localIP string,
	new_order <-chan structs.Order,
	assigned_new_order chan<- structs.Order,
	elev_receive_cost_value <-chan structs.Cost,
	elev_send_cost_value chan<- structs.Cost) {
	for {
		select {
		case current_new_order := <-new_order:
			fmt.Printf("Received new order\n")
			//Kjører kost funksjon og avgjør om den får bestillingen
			current_cost := structs.Cost{Cost_value: cost_function(current_new_order), Current_order: current_new_order, IP: localIP}
			fmt.Printf("Cost value: %d\n", current_cost.Cost_value)

			//elev_send_cost_value <- current_cost

			//Ta inn kost verdi fra nettverkskanal og kjøre action_select
			action_select(assigned_new_order, elev_receive_cost_value, current_cost)
		}
	}
}





