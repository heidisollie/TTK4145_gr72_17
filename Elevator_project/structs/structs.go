package structs

import (
	"../driver"
)

type ElevState struct {
	LastPassedFloor  int
	CurrentDirection driver.MotorDirection
	Stuck            bool
	IP               string
}

type Order struct {
	Type  driver.ButtonType
	Floor int
	IP    string
}

type Cost struct {
	CostValue    int
	CurrentOrder Order
}

const NUMELEV = 3

/*
func Dummy_func() {

}*/
