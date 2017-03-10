package localState

import(
	"../structs"
	"../driver"
)






var localElevator structs.Elev_state



func ReadLocalState() structs.Elev_state{
	return localElevator
}

func ChangeLocalState_flr(new_floor int){
	localElevator.Last_passed_floor = new_floor
}

func ChangeLocalState_dir(new_dir driver.MotorDirection){
	driver.SetMotorDirection(new_dir)
	localElevator.Current_direction = new_dir
}