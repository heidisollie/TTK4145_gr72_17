package FSM

import (
	"../driver"
	"../localElevator"
	def "../def"
	"fmt"
	"sync"
	"time"
	"os"
)

const stuckPeriod = 6 * time.Second
const statePeriod = 100 * time.Millisecond

type stateFSM int
const (
	idle   		stateFSM = 0
	moving 		stateFSM = 1
)

var mu sync.Mutex

func FSMInit(floorEvent <-chan int, newTargetFloor <-chan int, floorCompleted chan<- int, elevSendState chan<- def.ElevState) {

	targetFloor := -1
	state stateFSM := idle
	stateTicker := time.NewTicker(statePeriod)
	stuckTimer := time.NewTimer(stuckPeriod)
	stuckTimer.Stop()

	for {
		select {
		case <-stateTicker.C:
			mu.Lock()
			elevSendState <- localElevator.ReadLocalElevator()
			mu.Unlock()
		case <-stuckTimer.C:
			switch state {
			case idle:

			case moving:
				driver.SetMotorDirection(def.DirnStop)
				mu.Lock()
				localElevator.ChangeLocalElevator_dir(def.DirnStop)
				localElevator.ChangeLocalElevator_stuck(true)
				mu.Unlock()
				os.Exit(2)

			}

		case floor := <-newTargetFloor:
			if floor == drivet.GetFloorSignal() {
				driver.SetMotorDirection(def.DirnStop)
				mu.Lock()
				localElevator.ChangeLocalElevator_dir(def.DirnStop)
				mu.Unlock()
				floorCompleted <- localElevator.ReadLocalElevator().LastPassedFloor
				driver.OpenCloseDoor()
			}
					
			if floor == targetFloor {
				driver.OpenCloseDoor()
				floorCompleted <- floor
				break
			}
			if floor != -1 {			
				targetFloor = floor
			}

			switch state {
			case idle:
					stuckTimer.Reset(stuckPeriod)
					if targetFloor < localElevator.ReadLocalElevator().LastPassedFloor {
						driver.SetMotorDirection(def.DirnDown)
						mu.Lock()
						localElevator.ChangeLocalElevator_dir(def.DirnDown)
						mu.Unlock()
						state = moving
					} else if targetFloor > localElevator.ReadLocalElevator().LastPassedFloor {
						driver.SetMotorDirection(def.DirnUp)
						mu.Lock()
						localElevator.ChangeLocalElevator_dir(def.DirnUp)
						mu.Unlock()
						state = moving
					}
			case moving:
					if targetFloor < localElevator.ReadLocalElevator().LastPassedFloor {
						driver.SetMotorDirection(def.DirnDown)
						mu.Lock()
						localElevator.ChangeLocalElevator_dir(def.DirnDown)
						mu.Unlock()
						state = moving
					} else if targetFloor > localElevator.ReadLocalElevator().LastPassedFloor {
						driver.SetMotorDirection(def.DirnUp)
						mu.Lock()
						localElevator.ChangeLocalElevator_dir(def.DirnUp)
						mu.Unlock()
						state = moving
					}
			}

		case floor := <-floorEvent:
			driver.SetFloorIndicator(floor)
			if (floor == 0) || (floor == def.NumFloors-1) {
				driver.SetMotorDirection(def.DirnStop)
				mu.Lock()
				localElevator.ChangeLocalElevator_dir(def.DirnStop)
				mu.Unlock()
			}

			localElevator.ChangeLocalElevator_flr(floor)

			switch state {
			case idle:

			case moving:

				if targetFloor == -1 {
					break
				} else if targetFloor < floor {
					driver.SetMotorDirection(def.DirnDown)
					stuckTimer.Reset(stuckPeriod)
					mu.Lock()
					localElevator.ChangeLocalElevator_dir(def.DirnDown)
					mu.Unlock()
				} else if targetFloor > floor {
					driver.SetMotorDirection(def.DirnUp)
					stuckTimer.Reset(stuckPeriod)
					mu.Lock()
					localElevator.ChangeLocalElevator_dir(def.DirnUp)
					mu.Unlock()
				} else {
					driver.SetMotorDirection(def.DirnStop)
					mu.Lock()
					localElevator.ChangeLocalElevator_dir(def.DirnStop)
					mu.Unlock()
					stuckTimer.Stop()
					state = idle
					floorCompleted <- localElevator.ReadLocalElevator().LastPassedFloor
					driver.OpenCloseDoor()
				}
			}
		}
	}
}
