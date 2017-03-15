package structs

import (
	"../driver"
)

type MotorDirection int
const (
	DirnDown MotorDirection = -1
	DirnStop MotorDirection = 0
	DirnUp   MotorDirection = 1
)

type ButtonType int
const (
	ButtonCallDown    ButtonType = 0
	ButtonCallCommand ButtonType = 1
	ButtonCallUp      ButtonType = 2
)

type OrderButton struct {
	Type  ButtonType
	Floor int
}


const (
	NumFloors  = int(C.N_FLOORS)
	NumButtons = int(C.N_BUTTONS)
	NumElev = 3
)


const Filename = "orderBackup"


type ElevState struct {
	LastPassedFloor   int
	CurrentDirection  MotorDirection
	Stuck             bool
	IP                string
}

type Order struct {
	Type   ButtonType
	Floor  int
	IP     string
}

type Cost struct {
	CostValue    int
	CurrentOrder Order
}

