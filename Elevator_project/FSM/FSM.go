package FSM

import (
	"../driver"
	"../localState"
	"../structs"
	"fmt"
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
				fmt.Printf("FSM: Elevator timed out. State = stuck\n")
				driver.SetMotorDirection(driver.DirnStop)
				mu.Lock()
				localState.ChangeLocalState_dir(driver.DirnStop)
				localState.ChangeLocalState_stuck(true)
				mu.Unlock()
				State = stuck
				// Can end program here
			case stuck:
			}

		case floor := <-newTargetFloor:
			/*if floor == driver.GetFloorSignal() {
				driver.SetMotorDirection(driver.DirnStop)
				mu.Lock()
				localState.ChangeLocalState_dir(driver.DirnStop)
				mu.Unlock()
				fmt.Printf("We are at this floor \n")
				fmt.Printf("Door 1\n")
				driver.OpenCloseDoor()
				floorCompleted <- floor
			}*/
			if floor == targetFloor {
				driver.OpenCloseDoor()
				floorCompleted <- floor
				fmt.Printf("New target floor is current target floor \n")
				break
			}

			targetFloor = floor
			fmt.Printf("FSM: New target floor is %d\n", targetFloor+1)

			switch State {
			case idle:
				if targetFloor == -1 {
					fmt.Printf("FSM: No target floor\n")
				} else {
					fmt.Printf("Resetting stuck timer: 1 \n")
					stuckTimer.Reset(stuckPeriod)
					if targetFloor < localState.ReadLocalState().LastPassedFloor {
						fmt.Printf("Going down\n")
						driver.SetMotorDirection(driver.DirnDown)
						mu.Lock()
						localState.ChangeLocalState_dir(driver.DirnDown)
						mu.Unlock()
						State = moving
					} else if targetFloor > localState.ReadLocalState().LastPassedFloor {
						fmt.Printf("Going up\n")
						driver.SetMotorDirection(driver.DirnUp)
						mu.Lock()
						localState.ChangeLocalState_dir(driver.DirnUp)
						mu.Unlock()
						State = moving
					} else {
						fmt.Printf("Staying put \n")
						driver.SetMotorDirection(driver.DirnStop)
						State = idle
						mu.Lock()
						localState.ChangeLocalState_dir(driver.DirnStop)
						mu.Unlock()
						floorCompleted <- localState.ReadLocalState().LastPassedFloor
						fmt.Printf("Door 2\n")
						driver.OpenCloseDoor()
						fmt.Print("FSM: [IDLE] Reached target floor\n")

					}
				}
			case moving:
				if targetFloor == localState.ReadLocalState().LastPassedFloor {
					driver.SetMotorDirection(driver.DirnStop)
					stuckTimer.Stop()
					mu.Lock()
					localState.ChangeLocalState_dir(driver.DirnStop)
					mu.Unlock()
					State = idle
					floorCompleted <- localState.ReadLocalState().LastPassedFloor
					fmt.Printf("Door 3\n")
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
				fmt.Printf("Reaching bounds, stopping \n")
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

				if targetFloor == -1 {
					break
				} else if targetFloor < localState.ReadLocalState().LastPassedFloor {
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
