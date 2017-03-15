package driver

/*
 #cgo CFLAGS: -std=c11
 #cgo LDFLAGS: -lcomedi -lm
 #include "elev.h"
 #include "io.h"
*/
import (
 	"C"
	"fmt"
	"time"
	"../definitions"
)

const door_period = 3 * time.Second


func ElevInit() {
	C.elev_init()
	ClearAllButtonLamps()
	SetStopLamp(0)
	SetDoorOpenLamp(0)
	SetFloorIndicator(0)

	SetMotorDirection(def.DirnUp)
	for GetFloorSignal() == -1 {
	}
	SetMotorDirection(def.DirnStop)

}

func DriverInit(buttonEvent chan def.OrderButton, floorEvent chan int) {
	buttonWasActive := [def.NumFloors][def.NumButtons]int{
		{0, 0, 0},
		{0, 0, 0},
		{0, 0, 0},
		{0, 0, 0}}

	var buttonSignal, floorSignal int
	lastPassedFloor := -1
	for {
		floorSignal = GetFloorSignal()
		if floorSignal != lastPassedFloor && floorSignal != -1 {
			floorEvent <- floorSignal
			lastPassedFloor = floorSignal
			for button := def.ButtonCallDown; int(button) < NumButtons; button++ {
				buttonWasActive[floorSignal][button] = GetButtonSignal(button, floorSignal)
			}
		}
		for floor := 0; floor < def.NumFloors; floor++ {
			for button := def.ButtonCallDown; int(button) < NumButtons; button++ {
				if (floor == 0) && (button == def.ButtonCallDown) {
					continue
				}
				if (floor == def.NumFloors-1) && (button == def.ButtonCallUp) {
					continue
				}

				buttonSignal = GetButtonSignal(button, floor)
				if buttonSignal == 1 && (buttonWasActive[floor][button] == 0) {
					buttonEvent <- def.OrderButton{Type: button, Floor: floor}
					buttonWasActive[floor][button] = GetButtonSignal(button, floor)
					for i := 0; i < def.NumFloors; i++ {
						
					}
				}
			}
		}
	}
}

func OpenCloseDoor() {
	door_timer := time.NewTimer(def.door_period)
	SetDoorOpenLamp(1)
	<-door_timer.C
	SetDoorOpenLamp(0)
}

func ClearAllButtonLamps() {
	for floor := 0; floor < def.NumFloors; floor++ {
		if floor < NumFloors-1 {
			SetButtonLamp(def.ButtonCallUp, floor, 0)
		}
		if floor > 0 {
			SetButtonLamp(def.ButtonCallDown, floor, 0)
		}
		SetButtonLamp(def.ButtonCallCommand, floor, 0)
	}
}

func SetMotorDirection(dirn def.MotorDirection) {
	C.elev_set_motor_direction(C.elev_motor_direction_t(dirn))
}

func SetButtonLamp(button def.ButtonType, floor, value int) {
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

func GetButtonSignal(button def.ButtonType, floor int) int {
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
