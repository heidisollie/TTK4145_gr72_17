package order_distribution

import (
	//"fmt"
	. "../structs"
)

func cost_function(new_order Order) int {
	cost_value := 0

	return cost_value
}

func action_select(assignedNewOrder chan<- Order, elev_receive_cost_value <-chan Cost, newOrder Order) {

	cost_value := <-elev_receive_cost_value

	//if won:
	assignedNewOrder <- newOrder
}

func order_dist_init(new_order <-chan Order,
	assigned_new_order chan<- Order,
	elev_receive_cost_value <-chan Cost,
	elev_send_cost_value chan<- Cost) {
	for {
		select {
		case current_new_order := <-new_order:
			//Kjører kost funksjon og avgjør om den får bestillingen
			current_cost := Cost{cost_value: cost_function(current_new_order), Order: current_new_order, IP: localIP}
			elev_send_cost_value <- current_cost
			//Ta inn kost verdi fra nettverkskanal og kjøre action_select
			action_select(assigned_new_order, elev_receive_cost_value, current_new_order)
		}
	}
}
