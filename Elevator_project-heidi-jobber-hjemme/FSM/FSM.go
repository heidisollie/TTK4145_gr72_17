package FSM

import (
	"../driver"
	"../structs"
	"fmt"
	"time"
	"../localState"
	"sync"
)



//Want to change so we focus on choosing direction not target_floor


 // -------QUESTIONS-------------------------
//Do we need mutex for reading local state?

const stuck_period = 5 * time.Second
type FSM_state int

const (
	idle = 0
	moving = 1
	stuck  = -1
)


var mu sync.Mutex
	

func FSM_init(floor_event <-chan int, 
	new_target_order <-chan structs.Order, 
	floor_completed chan<- int, 
	elev_send_state chan<- structs.Elev_state) {


	FSM_state := idle
	stuck_timer := time.NewTimer(stuck_period)
	stuck_timer.Stop()



	for {
		select {
		case <-stuck_timer.C:
			switch FSM_state {
			case idle:
				//If elevator in idle and should move but does not, must be deaclared stuck.
			case moving:
				fmt.Printf("FSM: Elevator timed out. State = stuck\n")
				driver.SetMotorDirection(driver.DirnStop)
				mu.Lock()
				localState.ChangeLocalState_dir(driver.DirnStop)
				elev_send_state <- localState.ReadLocalState()
				mu.Unlock()
				FSM_state = stuck
			case stuck:
				//This should never be the case 
			}

		case order := <- new_target_order:
			if order.Floor == localState.ReadLocalState().Last_passed_floor {
				driver.OpenCloseDoor()
				floor_completed <- localState.ReadLocalState().Last_passed_floor
				break
			}
			fmt.Printf("FSM: New target order is: button: %d, floor: %d \n", order.Type, order.Floor+1)
			switch FSM_state {
			case idle:
					stuck_timer.Reset(stuck_period)
					if order.Floor < localState.ReadLocalState().Last_passed_floor{
						FSM_state = moving
						driver.SetMotorDirection(driver.DirnDown)
						mu.Lock()
						localState.ChangeLocalState_dir(driver.DirnDown)
						elev_send_state <- localState.ReadLocalState()
						mu.Unlock()
					} else if order.Floor > localState.ReadLocalState().Last_passed_floor {
						FSM_state = moving
						driver.SetMotorDirection(driver.DirnUp)
						mu.Lock()
						localState.ChangeLocalState_dir(driver.DirnUp)
						elev_send_state <- localState.ReadLocalState()
						mu.Unlock()
					}
			case moving:
				if order.Floor == localState.ReadLocalState().Last_passed_floor {
					driver.SetMotorDirection(driver.DirnStop)
					mu.Lock()
					localState.ChangeLocalState_dir(driver.DirnStop)
					elev_send_state <- localState.ReadLocalState()
					mu.Unlock()
					FSM_state = idle
					floor_completed <- localState.ReadLocalState().Last_passed_floor
					driver.OpenCloseDoor()
					fmt.Print("FSM: [MOVING] Reached target floor\n")

				}
			case stuck:
				fmt.Printf("FSM: [STUCK] Cannot take new orders \n")
			}

		case floor := <-floor_event:
			//driver.SetFloorIndicator(floor) - moved this to driver

			if (floor == 0) || (floor == driver.NumFloors-1) {
				driver.SetMotorDirection(driver.DirnStop)
				mu.Lock()
				localState.ChangeLocalState_dir(driver.DirnStop)
				elev_send_state <- localState.ReadLocalState()
				mu.Unlock()
			}

			mu.Lock()
			localState.ChangeLocalState_flr(floor)
			elev_send_state <- localState.ReadLocalState()
			mu.Unlock()

			switch FSM_state {
			case idle:
				fmt.Printf("FSM: CASE [floor event]: STATE [idle], floor: %d\n", floor+1)

			case moving:
				fmt.Printf("FSM: CASE [floor event]: STATE [moving], floor: %d\n", floor+1)
				stuck_timer.Reset(stuck_period)
				mu.Lock() //Read current floor under mutex
				state_flr := localState.ReadLocalState().Last_passed_floor
				mu.Unlock()
				if target_floor == -1 {
					break
				} else if target_floor < state_flr {
					driver.SetMotorDirection(driver.DirnDown)
					mu.Lock()
					elev_send_state <- localState.ReadLocalState()
					mu.Unlock()
				} else if target_floor > state_flr {
					driver.SetMotorDirection(driver.DirnUp)
					mu.Lock()
					localState.ChangeLocalState_dir(driver.DirnUp)
					elev_send_state <- localState.ReadLocalState()
					mu.Unlock()
				} else {
					driver.SetMotorDirection(driver.DirnStop)
					mu.Lock()
					localState.ChangeLocalState_dir(driver.DirnStop)
					elev_send_state <- localState.ReadLocalState()
					mu.Unlock()
					stuck_timer.Stop()
					FSM_state = idle
					mu.Lock()
					floor_completed <- localState.ReadLocalState().Last_passed_floor
					mu.Unlock()
					driver.OpenCloseDoor()
				}

			case stuck:
				fmt.Printf("FSM: STATE [stuck]; cannot reach floor %d\n", floor+1)

			} */
		}
	}
}

//We need two cases. First we need a goal. And in addition we need to check for other orders in the same direction that we can pick up.


//Currently we only set direction, and expect other_orders_in_dir to make it stop. But this only happens if it matches the criteria
//for "other orders in direction" which is not the case when it is the only order. 