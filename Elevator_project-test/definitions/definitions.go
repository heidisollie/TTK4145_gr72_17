package definitions

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
	NumFloors  = 4
	NumButtons = 3
	NumElev    = 3
)

const Filename = "orderBackup"

type Elevator struct {
	LastPassedFloor  int
	CurrentDirection MotorDirection
	Stuck            bool
	IP               string
	Online           bool
}

type Order struct {
	Type  ButtonType
	Floor int
	IP    string
}

type Cost struct {
	CostValue    int
	CurrentOrder Order
}
