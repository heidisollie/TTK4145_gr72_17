package FSM

import (
	def "../definitions"
	"../driver"
	elev "../localElevator"
	"fmt"
	"sync"
	"time"
)

const stuckPeriod = 6 * time.Second
const statePeriod = 200 * time.Millisecond
const doorPeriod = 2 * time.Second

type stateFSM int

const (
	idle     stateFSM = 0
	moving   stateFSM = 1
	doorOpen stateFSM = 2
	stuck    stateFSM = -1
)

var mu sync.Mutex

func elevGoDown(mu sync.Mutex) stateFSM {
	driver.SetMotorDirection(def.DirnDown)
	mu.Lock()
	elev.ChangeLocalElevator_dir(def.DirnDown)
	mu.Unlock()
	return moving

}

func elevGoUp(mu sync.Mutex) stateFSM {
	driver.SetMotorDirection(def.DirnUp)
	mu.Lock()
	elev.ChangeLocalElevator_dir(def.DirnUp)
	mu.Unlock()
	return moving
}

func elevStop(mu sync.Mutex) {
	driver.SetMotorDirection(def.DirnStop)
	mu.Lock()
	elev.ChangeLocalElevator_dir(def.DirnStop)
	mu.Unlock()
}

func FSMInit(floorEvent <-chan int,
	newTargetFloor <-chan int,
	floorCompleted chan<- int,
	elevSendState chan<- def.Elevator) {

	targetFloor := -1
	state := idle

	stateTicker := time.NewTicker(statePeriod)
	stuckTimer := time.NewTimer(stuckPeriod)
	stuckTimer.Stop()
	doorTimer := time.NewTimer(doorPeriod)
	doorTimer.Stop()

	for {
		select {
		case <-stateTicker.C:
			mu.Lock()
			elevSendState <- elev.ReadLocalElevator()
			mu.Unlock()
		case <-stuckTimer.C:
			switch state {
			case idle:

			case moving:
				state = stuck
				mu.Lock()
				elev.ChangeLocalElevator_stuck(true)
				mu.Unlock()
				elevSendState <- elev.ReadLocalElevator()
				fmt.Printf("Elevator timed out\n")

			case doorOpen:
			case stuck:
			}
		case <-doorTimer.C:
			switch state {
			case idle:
			case moving:
			case doorOpen:
				driver.SetDoorOpenLamp(0)
				fmt.Printf("Door closed\n")
				floorCompleted <- elev.ReadLocalElevator().LastPassedFloor
				state = idle
			case stuck:
			}

		case floor := <-newTargetFloor:
			fmt.Printf("New target floor given\n")
			targetFloor = floor
			stuckTimer.Reset(stuckPeriod)
			switch state {
			case idle:
				if targetFloor < elev.ReadLocalElevator().LastPassedFloor {
					state = elevGoDown(mu)
				} else if targetFloor > elev.ReadLocalElevator().LastPassedFloor {
					state = elevGoUp(mu)
				} else {
					elevStop(mu)
					doorTimer.Reset(doorPeriod)
					driver.SetDoorOpenLamp(1)
					state = doorOpen
				}
			case moving:
				if targetFloor < elev.ReadLocalElevator().LastPassedFloor {
					elevGoDown(mu)
				} else if targetFloor > elev.ReadLocalElevator().LastPassedFloor {
					elevGoUp(mu)
				} else {
					elevStop(mu)
					doorTimer.Reset(doorPeriod)
					driver.SetDoorOpenLamp(1)
					state = doorOpen
				}
			case doorOpen:
			case stuck:
			}

		case floor := <-floorEvent:
			driver.SetFloorIndicator(floor)
			elev.ChangeLocalElevator_flr(floor)
			if (floor == 0) || (floor == def.NumFloors-1) {
				elevStop(mu)
			}

			switch state {
			case idle:

			case moving:

				if targetFloor < elev.ReadLocalElevator().LastPassedFloor {
					elevGoDown(mu)
					stuckTimer.Reset(stuckPeriod)
				} else if targetFloor > elev.ReadLocalElevator().LastPassedFloor {
					elevGoUp(mu)
					stuckTimer.Reset(stuckPeriod)
				} else {
					elevStop(mu)
					stuckTimer.Stop()
					doorTimer.Reset(doorPeriod)
					driver.SetDoorOpenLamp(1)
					state = doorOpen
				}
			case doorOpen:

			case stuck:
				elevStop(mu)
				elev.ChangeLocalElevator_stuck(false)
				state = idle

			}
		}
	}
}
