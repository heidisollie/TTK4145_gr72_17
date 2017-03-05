package FSM

import (
	"../driver"
	. "../structs"
	"fmt"
	"time"
)

const door_period = 3 * time.Second

type FSM_state int

const (
	idle      = 0
	door_open = 1
	moving    = 2
	stuck     = -1
)

State := Elev_state{last_passed_floor: 0, current_direction = DirnStop, id = localIP}



func FSM_init(floor_event <-chan int, new_target_floor <-chan int, floor_completed chan<- int) {

	target_floor := -1
	FSM_state := idle

	door_timer := time.NewTimer(door_period)
	door_timer.Stop()

	stuck_timer := time.NewTimer(5 * time.Second)
	stuck_timer.Stop()

	for {
		select {
		case <-stuck_timer.C: //Går inn her om timeren går ut
			switch FSM_state {
			case idle:
			//skal ikke gjøre noe som helst
			case door_open:
			//do nothing
			case moving:
				fmt.Printf("Elevator timed out. State = stuck\n")
				FSM_state = stuck
			case stuck:
			}
		case <-door_timer.C:
			switch FSM_state {
			case idle:
			case door_open:

				fmt.Printf("STATE [door_open]\n")
				stuck_timer.Stop()
				driver.SetDoorOpenLamp(0)
				floor_completed <- Elev_state.last_passed_floor
				target_floor = -1
				FSM_state = idle

			case moving:
			case stuck:
			}

		case floor := <-new_target_floor:

			target_floor = floor
			fmt.Printf("New target floor is %d\n", target_floor+1)
			switch FSM_state {
			case idle:
				if target_floor == -1 {
					break
				} else {
					stuck_timer.Reset(5 * time.Second)
					if target_floor < Elev_state.last_passed_floor {
						FSM_state = moving
						driver.SetMotorDirection(driver.DirnDown)
						Elev_state.current_direction = driver.DirnDown
					} else if target_floor > Elev_state.last_passed_floor {
						FSM_state = moving
						driver.SetMotorDirection(driver.DirnUp)
						Elev_state.current_direction = driver.DirnUp
					} else {
						driver.SetMotorDirection(driver.DirnStop)
						door_timer.Reset(door_period)
						fmt.Print("Target floor is current floor\n")
						driver.SetDoorOpenLamp(1)
						FSM_state = door_open
					}
				}
			//disse skal ikke gjøre noe som helst
			case door_open:
			case moving:
				if target_floor == Elev_state.last_passed_floor {
					driver.SetMotorDirection(driver.DirnStop)
					door_timer.Reset(door_period)
					fmt.Print("Target floor is current floor\n")
					driver.SetDoorOpenLamp(1)
					FSM_state = door_open
				}
			case stuck:
			}

		case floor := <-floor_event:
			driver.SetFloorIndicator(floor)

			if (floor == 0) || (floor == driver.NumFloors-1) {
				Elev_state.current_direction = driver.DirnStop
			}

			Elev_state.last_past_floor = floor
			switch FSM_state {
			case idle:
				fmt.Printf("STATE [idle]; reached floor %d", floor+1)

			case door_open:
				fmt.Printf("STATE [doors open]; reached floor %d\n", floor+1)

			case moving:
				fmt.Printf("STATE [moving]; reached floor %d", floor+1)
				stuck_timer.Reset(5 * time.Second)

				if target_floor == -1 {
					break
				} else if target_floor < Elev_state.last_passed_floor {
					FSM_state = moving
					driver.SetMotorDirection(driver.DirnDown)
					Elev_state.current_direction = driver.DirnDown
				} else if target_floor > Elev_state.last_passed_floor {
					FSM_state = moving
					driver.SetMotorDirection(driver.DirnUp)
					Elev_state.current_direction = driver.DirnUp
				} else {
					driver.SetMotorDirection(driver.DirnStop)
					door_timer.Reset(door_period)
					driver.SetDoorOpenLamp(1)
					stuck_timer.Stop()
					FSM_state = door_open
					fmt.Printf("Opening doors\n")
				}

			case stuck:
				fmt.Printf("STAte [stuck]; cannot reach floor %d\n", floor+1)

			}
		}
	}
}
