package structs

import (
	"../driver"
)



type Elev_state struct {
	Last_passed_floor int
	Current_direction driver.MotorDirection
	IP                string
}

type Order struct {
	Type     driver.ButtonType
	Floor    int
	Internal bool
	IP       string
}

type Cost struct {
	Cost_value    int
	Current_order Order
	IP            string
}

	

func Dummy_func() {

}
