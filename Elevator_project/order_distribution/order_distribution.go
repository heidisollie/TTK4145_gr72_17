package order_distribution

import (
	//"fmt"
	"../structs"
	"strconv"
)


//Needs access to State elev_State
func cost_function(new_order structs.Order, ) int {
	cost_value := 0
	diff := new_order.Floor - State.Last_passed_floor

	if (State.Current_dirrection == 1 && diff > 0) || (State.Current_direction == -1 && diff < 0) {
		cost_value -= 0.5
	} else if (State.Current_direction == 1 && diff < 0) || (State.Current_direction == -1 && diff > 0) {
		cost_value += 2.25
	} else if (State.Current_direction == 0) {
		if (diff == -1) || (diff == 1) {
			cost_value += 0.25		
		} else if (diff == -2) || (diff == 2) {
			cost_value += 0.5
		} else if (diff == -3) || (diff == 3) {
			cost_value += 0.75
		}
	}
	return cost_value
}

//Sends order to order_handler of our cost is the lowest. If tie, picks the one with lower IP
func action_select(assignedNewOrder chan<- structs.Order, elev_receive_cost_value <-chan structs.Cost, cost Cost) {
	cost_value1 := <-elev_receive_cost_value
	cost_value2 := <-elev_receive_cost_value
	if cost.Cost_value < cost_value1.Cost_value && cost.Cost_value < cost_value2.Cost_value {
		assignedNewOrder <- newOrder
	} else if cost.Cost_value == cost_value1.Cost_value || cost.Cost_value == cost_value2.Cost_value {
		if cost.Cost_value == cost_value1.Cost_value && cost.Cost_value < cost_value2.Cost_value {
			if strconv(cost.IP) < strconv(cost_value1.IP) {
				assignedNewOrder <- newOrder
			}
		} else if cost.Cost_value == cost_value1.Cost_value && cost.Cost_value < cost_value2.Cost_value {
			if strconv(cost.IP) < strconv(cost_value2.IP) {
				assignedNewOrder <- newOrder
			}	
		}
	}
}


func Order_dist_init(localIP string,
	new_order <-chan structs.Order,
	assigned_new_order chan<- structs.Order,
	elev_receive_cost_value <-chan structs.Cost,
	elev_send_cost_value chan<- structs.Cost) {
	for {
		select {
		case current_new_order := <-new_order:

			//Kjører kost funksjon og avgjør om den får bestillingen
			current_cost := structs.Cost{Cost_value: cost_function(current_new_order), Current_order: current_new_order, IP: localIP}
			elev_send_cost_value <- current_cost

			//Ta inn kost verdi fra nettverkskanal og kjøre action_select
			action_select(assigned_new_order, elev_receive_cost_value, current_cost)
		}
	}
}





