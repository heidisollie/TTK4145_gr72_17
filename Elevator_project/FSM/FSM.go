package FSM

import (
	"../driver"
	"../structs"
	"log"
	"time"
)

const doorPeriod = 3 * time.Second

type FSM_state int

const (
	idle 		= 0
	door_open 	= 1
	moving 		= 2
	stuck 		= -1
)



func init(
	floor_event <- chan int
	//FSM -> elev_driver
	setmotordir chan<- int,
	doorstate chan<- bool,
	//Order_distr. -> FSM
	ordermotordir <-chan int,
	//Order_handler -> FSM
	newtargetfloor <-chan int,
	isorder <-chan bool,
	timer <-chan int,
	//FSM -> order_distr.
	state chan<- ElevState
	){
	
	target_floor := -1
	state := idle
	
	stuck_timer := time.NewTimer(time.Second*5)
	
	for{
		select {
			case <- stuck:
				FSM_state = stuck
				//Kø til heis skal merges med andre heiser
				//Stuck sendes på channel til order distribution 
			
			case motordir := <- ordermotordir:
				fmt.Printf("Setting motor direction: ", motordir)
				driver.SetMotorDirection(motordir)
				FSM_state = moving
				Elev_state.crrent_direction = motordir
	

			case <- isorder:
				driver.SetMotorDirection(driver.DirnStop)
				doorstate <- true
				FSM_state = doorOpen
				sleep(doorPeriod)
				doorstate <- false
				FSM_state = idle
			
			case <- stuck_timer.C:  //Går inn her om timeren går ut
				fmt.Printf("Elevator timed out. state=stuck")
				FSM_state = stuck
			
			case floor := <- new_target_floor:
				if floor != target_floor{
					// funksjon for å resette timer (for å sjekke stuck)
				}
				target_floor = floor
				fmt.Printf("New target floor is %d", target_floor+1)
				
				switch FSM_state {
					case idle:
						if target_floor == -1 {
							break
						}
						else if target_floor < Elev_state.LastPassedFloor{
							driver.SetMotorDirection(driver.DirnDown)
							Elev_state.CurrentDirection = driver.DirnDown
						}
						else if target_floor > Elev_state.LastPassedFloor{
							driver.SetMotorDirection(driver.DirnUp)
							Elev_state.CurrentDirection = driver.DirnUp
						}
						else {
							driver.SetMotorDirection(driver.DirnStop)
							FSM_state = door_open
						}
					case moving:
						// skal ikke gjøre noe som helst
					case door_open:
						driver.SetDoorOpenLamp(1)
						FSM_state = doorOpen
						sleep(doorPeriod)
						driver.SetDoorOpenLamp(0)
						FMS_state = idle
					case stuck:
									
					
			case floor := <- floor_event:
				
				driver.SetFloorIndicator(floor)
				
				//Sjekke om heisen er stuck
				
				if (floor == 0) || (floor == driver.numFloors-1){
						Elev_state.current_direction = driver.DirnStop
					}
					
				Elev_state.last_past_floor = floor
				switch FSM_state {
					case idle:
						fmt.Printf("STATE [idle]; reached floor %d", floor+1)
						
					case moving:
						fmt.Printf("STATE [moving]; reached floor %d", floor+1)
						stuck_timer.Reset(5*time.Second) //resetter timeren
						if target_floor == -1 {
							break
						}
						else if target_floor < Elev_state.LastPassedFloor{
						
							driver.SetMotorDirection(driver.DirnDown)
							Elev_state.current_direction = driver.DirnDown
						}
						else if target_floor > Elev_state.LastPassedFloor{
							driver.set_motor_direction(driver.DirnUp)
							Elev_state.current_direction = driver.DirnUp
						}
						else {
							driver.SetMotorDirection(driver.DirnStop)
							FSM_state = door_open
						}
					
						
					case door_open:
						fmt.Printf("STATE [doors open]; reached floor %d", floor+1)
					case stuck:
								
			}		
		}
	}	
}

	
