package FSM

import (
	"../driver"
	"../localState"
	"../structs"
	"fmt"
	"os"
	"sync"
	"time"
)

const stuckPeriod = 6 * time.Second
const statePeriod = 100 * time.Millisecond

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
				fmt.Printf("Elevator timed out. Exiting program \n")
				driver.SetMotorDirection(driver.DirnStop)
				State = stuck
				mu.Lock()
				localState.ChangeLocalState_dir(driver.DirnStop)
				localState.ChangeLocalState_stuck(true)
				mu.Unlock()
				elevSendState <- localState.ReadLocalState()
				os.Exit(2)

			case stuck:
			}

		case floor := <-newTargetFloor:

			fmt.Printf("Floor Signal %d, floor: %d \n", driver.GetFloorSignal()+1, floor+1)
			/*
				if floor == driver.GetFloorSignal() {
					driver.SetMotorDirection(driver.DirnStop)
					State = idle
					mu.Lock()
					localState.ChangeLocalState_dir(driver.DirnStop)
					mu.Unlock()
					driver.OpenCloseDoor()

					floorCompleted <- floor
					break
				}
			*/
			/*
				if floor == targetFloor {
					driver.OpenCloseDoor()
					floorCompleted <- floor
					fmt.Printf("New target floor is current target floor \n")
					break
				}
			*/

			targetFloor = floor
			fmt.Printf("FSM: New target floor is %d\n", targetFloor+1)
			stuckTimer.Reset(stuckPeriod)
			switch State {
			case idle:
				if targetFloor == -1 {
					//fmt.Printf("FSM: No target floor\n")
				} else {

					if targetFloor < localState.ReadLocalState().LastPassedFloor {
						driver.SetMotorDirection(driver.DirnDown)
						State = moving
						mu.Lock()
						localState.ChangeLocalState_dir(driver.DirnDown)
						mu.Unlock()

					} else if targetFloor > localState.ReadLocalState().LastPassedFloor {
						driver.SetMotorDirection(driver.DirnUp)
						State = moving
						mu.Lock()
						localState.ChangeLocalState_dir(driver.DirnUp)
						mu.Unlock()
					} else {
						driver.SetMotorDirection(driver.DirnStop)
						mu.Lock()
						localState.ChangeLocalState_dir(driver.DirnStop)
						mu.Unlock()
						driver.OpenCloseDoor()
						floorCompleted <- floor
					}
				}
			case moving:
				if targetFloor < localState.ReadLocalState().LastPassedFloor {
					driver.SetMotorDirection(driver.DirnDown)
					mu.Lock()
					localState.ChangeLocalState_dir(driver.DirnDown)
					mu.Unlock()

				} else if targetFloor > localState.ReadLocalState().LastPassedFloor {
					driver.SetMotorDirection(driver.DirnUp)
					mu.Lock()
					localState.ChangeLocalState_dir(driver.DirnUp)
					mu.Unlock()

				}
			case stuck:
				fmt.Printf("FSM: [STUCK] \n")
			}

		case floor := <-floorEvent:
			fmt.Printf("We're at a new floor\n")
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

				if targetFloor < localState.ReadLocalState().LastPassedFloor {
					driver.SetMotorDirection(driver.DirnDown)
					stuckTimer.Reset(stuckPeriod)
					mu.Lock()
					localState.ChangeLocalState_dir(driver.DirnDown)
					mu.Unlock()
				} else if targetFloor > localState.ReadLocalState().LastPassedFloor {
					driver.SetMotorDirection(driver.DirnUp)
					stuckTimer.Reset(stuckPeriod)
					mu.Lock()
					localState.ChangeLocalState_dir(driver.DirnUp)
					mu.Unlock()
				} else {
					fmt.Printf("FSM state MOVING and floor == targetFloor\n")
					State = idle
					driver.SetMotorDirection(driver.DirnStop)
					mu.Lock()
					localState.ChangeLocalState_dir(driver.DirnStop)
					mu.Unlock()
					stuckTimer.Stop()
					State = idle
					floorCompleted <- localState.ReadLocalState().LastPassedFloor
					fmt.Printf("Door 4\n")
					driver.OpenCloseDoor()
				}

			case stuck:
				fmt.Printf("FSM: STATE [stuck]; cannot reach floor %d\n", floor+1)

			}
		}
	}
}
