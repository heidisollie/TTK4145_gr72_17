package localState

import(
	"../structs"
	"../driver"
)



var localElevator structs.ElevState



func ReadLocalState() structs.ElevState{
	return localElevator
}

func ChangeLocalState_flr(newFloor int){
	localElevator.LastPassedFloor = newFloor
}

func ChangeLocalState_dir(newDir driver.MotorDirection){
	driver.SetMotorDirection(newDir)
	localElevator.CurrentDirection = newDir
}