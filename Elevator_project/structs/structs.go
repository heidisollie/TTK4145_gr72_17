package structs

import (
	"log"
	"time"
	"../network"
	"encoding/json"
	
)


type Elev_state struct{
	last_passed_floor  int
	current_direction driver.MotorDirection
	id int
}

type Order struct{
	ButtonType int
	Floor int
	Internal bool
}

type Cost struct {
	cost_value int
	Order Order
	IP string
}

