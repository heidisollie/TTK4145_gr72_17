package FSM

import (
	"../driver"
	"../localState"
	"../structs"
	"fmt"
	"sync"
	"time"
)

const stuckPeriod = 5 * time.Second
const statePeriod = 10 * time.Millisecond

type FSMState int

const (
	idle   = 0
	moving = 1
	stuck  = -1
)

var mu sync.Mutex

func FSMInit(floorEvent <-chan int, newTargetFloor <-chan int, floorCompleted chan<- int, elevSendState chan<- structs.ElevState) {

	targetFloor := -1
	State := idle

	stateTicker := time.NewTicker(statePeriod)
	stuckTimer := time.NewTimer(stuckPeriod)
	stuckTimer.Stop()

	for {
		select {
		case <-stateTicker.C:
			mu.Lock()
			elevSendState <- localState.ReadLocalState()
			mu.Unlock()
		case <-stuckTimer.C:
			switch State {
			case idle:

			case moving:

				fmt.Printf("FSM: Elevator timed out. State = stuck\n")
				driver.SetMotorDirection(driver.DirnStop)
				mu.Lock()
				localState.ChangeLocalState_dir(driver.DirnStop)
				mu.Unlock()
				State = stuck
			case stuck:
			}

		case floor := <-newTargetFloor:
			fmt.Printf("Received new target floor\n")
			if floor == targetFloor {
				fmt.Printf("Floor is target floor\n")
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
					if targetFloor < localState.ReadLocalState().LastPassedFloor {
						State = moving
						driver.SetMotorDirection(driver.DirnDown)
						mu.Lock()
						localState.ChangeLocalState_dir(driver.DirnDown)
						mu.Unlock()
					} else if targetFloor > localState.ReadLocalState().LastPassedFloor {

						State = moving
						driver.SetMotorDirection(driver.DirnUp)
						mu.Lock()
						localState.ChangeLocalState_dir(driver.DirnUp)
						mu.Unlock()
					} else {
						driver.SetMotorDirection(driver.DirnStop)
						State = idle
						mu.Lock()
						localState.ChangeLocalState_dir(driver.DirnStop)
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
			fmt.Printf("Received new floor event\n")
			driver.SetFloorIndicator(floor)

			if (floor == 0) || (floor == driver.NumFloors-1) {
				driver.SetMotorDirection(driver.DirnStop)
				mu.Lock()
				localState.ChangeLocalState_dir(driver.DirnStop)
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
					mu.Unlock()
				} else if targetFloor > localState.ReadLocalState().LastPassedFloor {
					driver.SetMotorDirection(driver.DirnUp)
					mu.Lock()
					localState.ChangeLocalState_dir(driver.DirnUp)
					mu.Unlock()
				} else {
					driver.SetMotorDirection(driver.DirnStop)
					mu.Lock()
					localState.ChangeLocalState_dir(driver.DirnStop)
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
