package FSM

import (
	"../driver"
	"../structs"
	"fmt"
	"time"
	"../localState"
	"sync"
)

//------------------------------------------------

//--------------------------------------------------


const stuckPeriod = 5 * time.Second

type FSMState int

const (
	idle = 0
	moving = 1
	stuck  = -1
)


var mu sync.Mutex
	

func FSMInit(floorEvent <-chan int, newTargetFloor <-chan int, floorCompleted chan<- int, elevSendState chan<- structs.ElevState) {

	targetFloor := -1
	State := idle

	stuckTimer := time.NewTimer(stuckPeriod)
	stuckTimer.Stop()



	for {
		select {
		case <-stuckTimer.C:
			switch State {
			case idle:

			case moving:
				
				fmt.Printf("FSM: Elevator timed out. State = stuck\n")
				driver.SetMotorDirection(driver.DirnStop)
				mu.Lock()
				localState.ChangeLocalState_dir(driver.DirnStop)
				elevSendState <- localState.ReadLocalState()
				mu.Unlock()
				State = stuck
			case stuck:
			}

		case floor := <-newTargetFloor:
			if floor == targetFloor {
				break
			}
			targetFloor = floor
			fmt.Printf("FSM: New target floor is %d\n", targetFloor+1)
			
			switch State {
			case idle:
				if targetFloor == -1 {
					//fmt.Printf("FSM: No target floor\n")
				} else {
					stuckTimer.Reset(stuckPeriod)
					if targetFloor < localState.ReadLocalState().LastPassedFloor{
						State = moving
						driver.SetMotorDirection(driver.DirnDown)
						mu.Lock()
						localState.ChangeLocalState_dir(driver.DirnDown)
						elevSendState <- localState.ReadLocalState()
						mu.Unlock()
					} else if targetFloor > localState.ReadLocalState().LastPassedFloor {
						mu.Lock()
						State = moving
						mu.Unlock()
						driver.SetMotorDirection(driver.DirnUp)
						localState.ChangeLocalState_dir(driver.DirnUp)
						elevSendState <- localState.ReadLocalState()
					} else {
						driver.SetMotorDirection(driver.DirnStop)
						State = idle
						mu.Lock()
						localState.ChangeLocalState_dir(driver.DirnStop)
						elevSendState <- localState.ReadLocalState()
						mu.Unlock()
						floorCompleted <- localState.ReadLocalState().LastPassedFloor
						driver.OpenCloseDoor()
						fmt.Print("FSM: [IDLE] Reached target floor\n")

					}
				}
			case moving:
				if targetFloor == localState.ReadLocalState().LastPassedFloor {
					driver.SetMotorDirection(driver.DirnStop)
					mu.Lock()
					localState.ChangeLocalState_dir(driver.DirnStop)
					elevSendState <- localState.ReadLocalState()
					mu.Unlock()
					State = idle
					floorCompleted <- localState.ReadLocalState().LastPassedFloor
					driver.OpenCloseDoor()
					fmt.Print("FSM: [MOVING] Reached target floor\n")

				}
			case stuck:
				fmt.Printf("FSM: [STUCK] \n")
			}

		case floor := <-floorEvent:
			driver.SetFloorIndicator(floor)

			if (floor == 0) || (floor == driver.NumFloors-1) {
				driver.SetMotorDirection(driver.DirnStop)
				mu.Lock()
				localState.ChangeLocalState_dir(driver.DirnStop)
				elevSendState <- localState.ReadLocalState()
				mu.Unlock()
			}

			localState.ChangeLocalState_flr(floor)

			switch State {
			case idle:
				fmt.Printf("FSM: CASE [floor event]: STATE [idle], floor: %d\n", floor+1)

			case moving:
				fmt.Printf("FSM: CASE [floor event]: STATE [moving], floor: %d\n", floor+1)
				stuckTimer.Reset(stuckPeriod)

				if targetFloor == -1 {
					break
				} else if targetFloor < localState.ReadLocalState().LastPassedFloor {
					driver.SetMotorDirection(driver.DirnDown)
					mu.Lock()
					localState.ChangeLocalState_dir(driver.DirnDown)
					elevSendState <- localState.ReadLocalState()
					mu.Unlock()
				} else if targetFloor > localState.ReadLocalState().LastPassedFloor {
					driver.SetMotorDirection(driver.DirnUp)
					mu.Lock()
					localState.ChangeLocalState_dir(driver.DirnUp)
					elevSendState <- localState.ReadLocalState()
					mu.Unlock()
				} else {
					driver.SetMotorDirection(driver.DirnStop)
					mu.Lock()
					localState.ChangeLocalState_dir(driver.DirnStop)
					elevSendState <- localState.ReadLocalState()
					mu.Unlock()
					stuckTimer.Stop()
					State = idle
					floorCompleted <- localState.ReadLocalState().LastPassedFloor
					driver.OpenCloseDoor()
				}
			case stuck:
				fmt.Printf("FSM: STATE [stuck]; cannot reach floor %d\n", floor+1)

			}
		}
	}
}