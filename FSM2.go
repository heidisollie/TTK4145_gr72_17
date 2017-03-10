package FSM2

import (
	"../driver"
	//"../structs"
	"fmt"
	"time"
	"../localState"
)

//------------------------------------------------

//--------------------------------------------------


const stuck_period = 5 * time.Second

type FSM_state int

const (
	idle = 0
	//door_open = 1
	moving = 1
	stuck  = -1
)

func stuckTimerOut(state FSM_state){
		switch ReadState_dir() {	
		case driver.DirnStop:

		case driver.DirnUp:
			fmt.Printf("FSM: Elevator timed out. State = stuck\n")
			ChangeState_dir(driver.DirnStop)

		case driver.DirnDown: 
			fmt.Printf("FSM: Elevator timed out. State = stuck\n")
			ChangeState_dir(driver.DirnStop)
			

		//case stuck:
	}	
}

func ReadState_floor() int{
	return localState.ReadLocalState().Last_passed_floor
}

func ReadState_dir() driver.MotorDirection{
	return localState.ReadLocalState().Current_direction
}
func ChangeState_dir(dir driver.MotorDirection) {
	localState.ChangeLocalState_dir(dir)
}


func newTargetFloor(floor int, target_floor int, floor_completed chan<- int, ) int {
		target_floor = floor
		fmt.Printf("FSM: CASE [target floor]: New target floor is %d\n", target_floor+1)
		switch ReadState_dir() {
			case driver.DirnStop:
					if target_floor < ReadState_floor(){
						ChangeState_dir(driver.DirnDown)
					} else if target_floor > ReadState_floor() {
						ChangeState_dir(driver.DirnUp)
					} else {
						ChangeState_dir(driver.DirnStop)
						floor_completed <- ReadState_floor()
						driver.OpenCloseDoor()
						fmt.Print("FSM: [IDLE] Reached target floor\n")
					}
	
			case driver.DirnDown:

			case driver.DirnUp:/*
				if target_floor == ReadState_floor() {
					ChangeState_dir(driver.DirnStop)
					state = idle
					floor_completed <- ReadState_floor()
					driver.OpenCloseDoor()
					fmt.Print("FSM: [MOVING] Reached target floor\n")
				}
				*/
		}	
		return target_floor
}

func newFloorEvent(floor int, target_floor int, floor_completed chan<- int) (bool) {

		driver.SetFloorIndicator(floor)
		localState.ChangeLocalState_flr(floor)
		var timer bool = false

		switch ReadState_dir(){
			case driver.DirnStop:
				fmt.Printf("FSM: CASE [floor event]: STATE [idle], floor: %d\n", floor+1)
			case driver.DirnDown:
				fmt.Printf("FSM: CASE [floor event]: STATE [moving], floor: %d\n", floor+1)
				

				if target_floor == -1 {
					break
				} else if target_floor < ReadState_floor() {
					ChangeState_dir(driver.DirnDown)
					timer = true
				} else if target_floor > ReadState_floor() {
					ChangeState_dir(driver.DirnUp)
					timer = true
				} else {
					ChangeState_dir(driver.DirnStop)
					floor_completed <- ReadState_floor()
					driver.OpenCloseDoor()
					timer = true
				}
		case driver.DirnUp:
				fmt.Printf("FSM: CASE [floor event]: STATE [moving], floor: %d\n", floor+1)
				

				if target_floor == -1 {
					break
				} else if target_floor < ReadState_floor() {
					ChangeState_dir(driver.DirnDown)
					timer = true
				} else if target_floor > ReadState_floor() {
					ChangeState_dir(driver.DirnUp)
					timer = true
				} else {
					ChangeState_dir(driver.DirnStop)
					floor_completed <- ReadState_floor()
					driver.OpenCloseDoor()
				}	
		}
		return timer

}

func FSM_init(floor_event <-chan int, new_target_floor <-chan int, floor_completed chan<- int) {

	target_floor := -1
	var state FSM_state = idle
	var timer bool
	stuck_timer := time.NewTimer(stuck_period)
	stuck_timer.Stop()

	for {
		select {
			case <-stuck_timer.C:
				stuckTimerOut(state)
	
			case floor := <-new_target_floor:
			
			if floor == target_floor {
				driver.OpenCloseDoor()
				floor_completed <- ReadState_floor()
			}

			stuck_timer.Reset(stuck_period)
			target_floor = newTargetFloor(floor, target_floor, floor_completed)
			case floor := <-floor_event:
				if (floor == 0) || (floor == driver.NumFloors-1) {
					ChangeState_dir(driver.DirnStop)
				}
			
			timer = newFloorEvent(floor, target_floor, floor_completed)
				if (timer == true) {
					stuck_timer.Reset(stuck_period)
				}


		}
	}
}