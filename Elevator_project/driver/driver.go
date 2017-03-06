package driver

/*
 #cgo CFLAGS: -std=c11
 #cgo LDFLAGS: -lcomedi -lm
 #include "elev.h"
 #include "io.h"
*/
import "C"
import "fmt"
import "time"

type MotorDirection int

const (
	DirnDown MotorDirection = -1
	DirnStop MotorDirection = 0
	DirnUp   MotorDirection = 1
)

type ButtonType int

const door_period = 3 * time.Second

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
)

func ElevInit() {
	C.elev_init()
	ClearAllButtonLamps()
	SetStopLamp(0)
	SetDoorOpenLamp(0)
	SetFloorIndicator(0)

	SetMotorDirection(DirnUp)
	for GetFloorSignal() == -1 {
	}
	SetMotorDirection(DirnStop)

}

func EventListener(button_event chan OrderButton, floor_event chan int) {
	buttonWasActive := [NumFloors][NumButtons]int{
		{0, 0, 0},
		{0, 0, 0},
		{0, 0, 0},
		{0, 0, 0}}

	var buttonSignal, floorSignal int
	lastPassedFloor := -1
	for {
		floorSignal = GetFloorSignal()
		if floorSignal != lastPassedFloor && floorSignal != -1 {
			floor_event <- floorSignal
			lastPassedFloor = floorSignal
			fmt.Println(floorSignal)
			for button := ButtonCallDown; int(button) < NumButtons; button++ {
				SetButtonLamp(button, floorSignal, 0)
				buttonWasActive[floorSignal][button] = GetButtonSignal(button, floorSignal)
			}
		}
		for floor := 0; floor < NumFloors; floor++ {
			for button := ButtonCallDown; int(button) < NumButtons; button++ {
				if (floor == 0) && (button == ButtonCallDown) {
					continue
				}
				if (floor == NumFloors-1) && (button == ButtonCallUp) {
					continue
				}

				buttonSignal = GetButtonSignal(button, floor)
				if buttonSignal == 1 && (buttonWasActive[floor][button] == 0) {
					button_event <- OrderButton{Type: button, Floor: floor}
					buttonWasActive[floor][button] = GetButtonSignal(button, floor)
					SetButtonLamp(button, floor, 1)
					for i := 0; i < NumFloors; i++ {
						fmt.Print(buttonWasActive[i])
						fmt.Printf("\n")
					}
					fmt.Printf("\n")
				}
			}
		}
	}
}
func OpenCloseDoor() {
	door_timer := time.NewTimer(door_period)
	SetDoorOpenLamp(1)
	<-door_timer.C
	SetDoorOpenLamp(0)
}

func ClearAllButtonLamps() {
	for floor := 0; floor < NumFloors; floor++ {
		if floor < NumFloors-1 {
			SetButtonLamp(ButtonCallUp, floor, 0)
		}
		if floor > 0 {
			SetButtonLamp(ButtonCallDown, floor, 0)
		}
		SetButtonLamp(ButtonCallCommand, floor, 0)
	}
}

func SetMotorDirection(dirn MotorDirection) {
	C.elev_set_motor_direction(C.elev_motor_direction_t(dirn))
}

func SetButtonLamp(button ButtonType, floor, value int) {
	C.elev_set_button_lamp(C.elev_button_type_t(button), C.int(floor), C.int(value))
}

func SetFloorIndicator(floor int) {
	C.elev_set_floor_indicator(C.int(floor))
}

func SetDoorOpenLamp(value int) {
	C.elev_set_door_open_lamp(C.int(value))
}

func SetStopLamp(value int) {
	C.elev_set_stop_lamp(C.int(value))
}

func GetButtonSignal(button ButtonType, floor int) int {
	return int(C.elev_get_button_signal(C.elev_button_type_t(button), C.int(floor)))
}

func GetFloorSignal() int {
	return int(C.elev_get_floor_sensor_signal())
}

func GetStopSignal() int {
	return int(C.elev_get_stop_signal())
}

func GetObstructionSignal() int {
	return int(C.elev_get_obstruction_signal())
}
