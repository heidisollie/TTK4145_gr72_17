package structs

import (
	"../driver"
)

type Elev_state struct {
	last_passed_floor int
	current_direction driver.MotorDirection
	id                int
}

type Order struct {
	Type     driver.ButtonType
	Floor    int
	Internal bool
	IP       string
}

type Cost struct {
	cost_value int
	Order      Order
	IP         string
}
