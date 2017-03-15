package localElevator

import (
	def "../definitions"
)

var localElevator def.ElevState

func ReadLocalElevator() def.ElevState {
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
