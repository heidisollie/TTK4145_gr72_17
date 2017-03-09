package FSM

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




func FSM_init(floor_event <-chan int, new_target_floor <-chan int, floor_completed chan<- int) {

	target_floor := -1
	FSM_state := idle

	stuck_timer := time.NewTimer(stuck_period)
	stuck_timer.Stop()

	for {
		select {
		case <-stuck_timer.C:
			switch FSM_state {
			case idle:

			case moving:
				fmt.Printf("FSM: Elevator timed out. State = stuck\n")
				driver.SetMotorDirection(driver.DirnStop)
				localState.ChangeLocalState_dir(driver.DirnStop)
				FSM_state = stuck
			case stuck:
			}

		case floor := <-new_target_floor:
			if floor == target_floor {
				break
			}
			target_floor = floor
			fmt.Printf("FSM: New target floor is %d\n", target_floor+1)
			if FSM_state == 0{
				fmt.Printf("FSM: State: idle\n")
			} else if FSM_state == 1 {
				fmt.Printf("FSM: State moving\n")
			} else {
				fmt.Printf("FSM: State stuck\n")
			}
			
			switch FSM_state {
			case idle:
				if target_floor == -1 {
					fmt.Printf("FSM: No target floor\n")
				} else {
					stuck_timer.Reset(stuck_period)
					if target_floor < localState.ReadLocalState().Last_passed_floor{
						FSM_state = moving
						driver.SetMotorDirection(driver.DirnDown)
						localState.ChangeLocalState_dir(driver.DirnDown)
					} else if target_floor > localState.ReadLocalState().Last_passed_floor {
						FSM_state = moving
						driver.SetMotorDirection(driver.DirnUp)
						localState.ChangeLocalState_dir(driver.DirnUp)
					} else {
						driver.SetMotorDirection(driver.DirnStop)
						FSM_state = idle
						floor_completed <- localState.ReadLocalState().Last_passed_floor
						driver.OpenCloseDoor()
						fmt.Print("FSM: [IDLE] Reached target floor\n")
						fmt.Printf("FSM: [case: new_target_floor, state = idle] Opening doors\n")

					}
				}
			case moving:
				if target_floor == localState.ReadLocalState().Last_passed_floor {
					driver.SetMotorDirection(driver.DirnStop)
					localState.ChangeLocalState_dir(driver.DirnStop)
					FSM_state = idle

					floor_completed <- localState.ReadLocalState().Last_passed_floor
					driver.OpenCloseDoor()
					fmt.Print("FSM: [MOVING] Reached target floor\n")
					fmt.Printf("FSM: [case: new_target_floor, state = moving] Opening doors\n")

				}
			case stuck:
				fmt.Printf("FSM: STUUCk\n")
			}

		case floor := <-floor_event:
			driver.SetFloorIndicator(floor)

			if (floor == 0) || (floor == driver.NumFloors-1) {
				driver.SetMotorDirection(driver.DirnStop)
				localState.ChangeLocalState_dir(driver.DirnStop)
			}

			localState.ChangeLocalState_flr(floor)

			switch FSM_state {
			case idle:
				fmt.Printf("FSM: CASE [floor event]: STATE [idle], floor: %d\n", floor+1)

			case moving:
				fmt.Printf("FSM: CASE [floor event]: STATE [moving], floor: %d\n", floor+1)
				stuck_timer.Reset(stuck_period)

				if target_floor == -1 {
					break
				} else if target_floor < localState.ReadLocalState().Last_passed_floor {
					FSM_state = moving
					driver.SetMotorDirection(driver.DirnDown)
					localState.ChangeLocalState_dir(driver.DirnDown)
				} else if target_floor > localState.ReadLocalState().Last_passed_floor {
					FSM_state = moving
					driver.SetMotorDirection(driver.DirnUp)
					localState.ChangeLocalState_dir(driver.DirnUp)
				} else {
					driver.SetMotorDirection(driver.DirnStop)
					localState.ChangeLocalState_dir(driver.DirnStop)
					stuck_timer.Stop()
					FSM_state = idle
					floor_completed <- localState.ReadLocalState().Last_passed_floor
					driver.OpenCloseDoor()
					fmt.Printf("FSM: [case: floor_event] Opening doors\n")
				}

			case stuck:
				fmt.Printf("FSM: STATE [stuck]; cannot reach floor %d\n", floor+1)

			}
		}
	}
}
