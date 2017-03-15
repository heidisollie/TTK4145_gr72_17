package FSM

import (
	"../driver"
	"../localState"
	"../structs"
	"fmt"
	//"os"
	"sync"
	"time"
)

const stuckPeriod = 6 * time.Second
const statePeriod = 200 * time.Millisecond
const doorPeriod = 2 * time.Second

type FSMState int

const (
	idle     FSMState = 0
	moving   FSMState = 1
	doorOpen FSMState = 2
	stuck    FSMState = -1
)

var mu sync.Mutex

func elevGoDown(mu sync.Mutex) FSMState {
	driver.SetMotorDirection(driver.DirnDown)
	mu.Lock()
	localState.ChangeLocalState_dir(driver.DirnDown)
	mu.Unlock()
	return moving

}

func elevGoUp(mu sync.Mutex) FSMState {
	driver.SetMotorDirection(driver.DirnUp)
	mu.Lock()
	localState.ChangeLocalState_dir(driver.DirnUp)
	mu.Unlock()
	return moving
}

func elevStop(mu sync.Mutex) {
	driver.SetMotorDirection(driver.DirnStop)
	mu.Lock()
	localState.ChangeLocalState_dir(driver.DirnStop)
	mu.Unlock()

}

func FSMInit(floorEvent <-chan int, newTargetFloor <-chan int, floorCompleted chan<- int, elevSendState chan<- structs.ElevState) {

	targetFloor := -1
	State := idle

	stateTicker := time.NewTicker(statePeriod)
	stuckTimer := time.NewTimer(stuckPeriod)
	stuckTimer.Stop()
	doorTimer := time.NewTimer(doorPeriod)
	doorTimer.Stop()

	for {
		select {
		case <-stateTicker.C:
			mu.Lock()
			////fmt.Printf("My local states: flr: %d, dir: %d, IP: %s \n", localState.ReadLocalState().LastPassedFloor, localState.ReadLocalState().CurrentDirection, localState.ReadLocalState().IP)
			elevSendState <- localState.ReadLocalState()
			mu.Unlock()
		case <-stuckTimer.C:
			switch State {
			case idle:

			case moving:
				////fmt.Printf("Elevator timed out. Exiting program \n")
				//elevStop(mu)
				State = stuck
				mu.Lock()
				localState.ChangeLocalState_stuck(true)
				mu.Unlock()
				elevSendState <- localState.ReadLocalState()
				//driver.ClearAllButtonLamps
				fmt.Printf("Elevator timed out\n")

			case stuck:
			}
		case <-doorTimer.C:
			switch State {
			case idle:
				////fmt.Printf("Skal aldri printes 1\n")
			case doorOpen:
				driver.SetDoorOpenLamp(0)
				floorCompleted <- localState.ReadLocalState().LastPassedFloor
				State = idle
			case moving:
				//fmt.Printf("Skal aldri printes 2\n")
			case stuck:
			}

		case floor := <-newTargetFloor:

			//fmt.Printf("Floor Signal %d, floor: %d \n", driver.GetFloorSignal()+1, floor+1)

			targetFloor = floor
			fmt.Printf("FSM: New target floor is %d\n", targetFloor+1)
			stuckTimer.Reset(stuckPeriod)
			switch State {
			case idle:
				//if targetFloor != -1 {
				if targetFloor < localState.ReadLocalState().LastPassedFloor {
					State = elevGoDown(mu)
				} else if targetFloor > localState.ReadLocalState().LastPassedFloor {
					State = elevGoUp(mu)
				} else {
					elevStop(mu)
					doorTimer.Reset(doorPeriod)
					driver.SetDoorOpenLamp(1)
					State = doorOpen
				}
				//}
			case moving:
				if targetFloor < localState.ReadLocalState().LastPassedFloor {
					elevGoDown(mu)

				} else if targetFloor > localState.ReadLocalState().LastPassedFloor {
					elevGoUp(mu)

				}
			case doorOpen:
				//fmt.Printf("Door open\n")
			case stuck:
				//fmt.Printf("FSM: [STUCK] \n")
			}

		case floor := <-floorEvent:
			fmt.Printf("We're at a new floor: %d \n", floor+1)
			driver.SetFloorIndicator(floor)
			localState.ChangeLocalState_flr(floor)
			if (floor == 0) || (floor == driver.NumFloors-1) {
				elevStop(mu)
			}

			switch State {
			case idle:
				fmt.Printf("FSM: CASE [floor event]: STATE [idle], floor: %d\n", floor+1)

			case moving:
				fmt.Printf("FSM: CASE [floor event]: STATE [moving], floor: %d\n", floor+1)

				if targetFloor < localState.ReadLocalState().LastPassedFloor {
					elevGoDown(mu)
					stuckTimer.Reset(stuckPeriod)
				} else if targetFloor > localState.ReadLocalState().LastPassedFloor {
					elevGoUp(mu)
					stuckTimer.Reset(stuckPeriod)
				} else {
					//fmt.Printf("FSM state MOVING and floor == targetFloor\n")
					elevStop(mu)
					stuckTimer.Stop()
					doorTimer.Reset(doorPeriod)
					driver.SetDoorOpenLamp(1)
					State = doorOpen
				}
			case doorOpen:
				//fmt.Printf("Door open\n")
			case stuck:
				elevStop(mu)
				localState.ChangeLocalState_stuck(false)
				State = idle

			}
		}
	}
}
