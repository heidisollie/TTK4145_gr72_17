package FSM

import (
	"../driver"
	"log"
	"time"
)

const doorPeriod = 3 * time.Second

type FSM_state int

const (
	idle 		= 0
	doorOpen 	= 1
	moving 		= 2
	stuck 		= -1
)

type ElevState struct{
	LastPassedFloor  int
	CurrentDirection driver.MotorDirection
	id int
}



func init(
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
	
	state := idle

	
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
				ElevState.CurrentDirection = motordir
	

			case <- isorder:
				driver.SetMotorDirection(driver.DirnStop)
				doorstate <- true
				FSM_state = doorOpen
				sleep(doorPeriod)
				doorstate <- false
				FSM_state = idle
			
			case floor := <- newtargetfloor:
				switch FSM_state {
					case idle:
						if floor == -1 {
							break
						}
						else if floor < ElevState.LastPassedFloor{
							driver.SetMotorDirection(driver.DirnDown)
							ElevState.CurrentDirection = driver.DirnDown
						}
						else if floor > ElevState.LastPassedFloor{
							driver.SetMotorDirection(driver.DirnUp)
							ElevState.CurrentDirection = driver.DirnUp
						}
						else {
							driver.SetMotorDirection(driver.DirnStop)
							driver.SetDoorOpenLamp(1)
							FSM_state = doorOpen
							sleep(doorPeriod)
							driver.SetDoorOpenLamp(0)
							FMS_state = idle
						}
					case moving:
					case dooropen:
					case stuck:
					
				}		
		}
	}	
}

	
