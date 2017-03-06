package order_distribution

import (
	//"fmt"
	"../structs"
)

func cost_function(new_order structs.Order) int {
	cost_value := 0
	//Movement penalty

	//Turn penalty

	//Order_dir_penalty

	return cost_value
}

func action_select(assignedNewOrder chan<- structs.Order, elev_receive_cost_value <-chan structs.Cost, newOrder structs.Order) {

	cost_value := <-elev_receive_cost_value
	if cost_value.Cost_value < 10 {

	}
	//if won:
	assignedNewOrder <- newOrder
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
			action_select(assigned_new_order, elev_receive_cost_value, current_new_order)
		}
	}
}
