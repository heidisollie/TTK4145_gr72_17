package localElevator

import (
	def "../definitions"
)

var localElevator def.Elevator

func ReadLocalElevator() def.Elevator {
	return localElevator
}

func ChangeLocalElevator_flr(newFloor int) {
	localElevator.LastPassedFloor = newFloor
}

func ChangeLocalElevator_dir(newDir def.MotorDirection) {
	localElevator.CurrentDirection = newDir
}

func ChangeLocalElevator_stuck(stuck bool) {
	localElevator.Stuck = stuck
}

func ChangeLocalElevator_IP(IP string) {
	localElevator.IP = IP
}

func ChangeLocalElevator_Online(state bool) {
	localElevator.Online = state
}
