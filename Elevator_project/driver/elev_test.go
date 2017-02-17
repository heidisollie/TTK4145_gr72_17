package main

import (

	"fmt"
	"./driver"
	"time"
)

const NumFloors = driver.NumFloors
const NumButtons = driver.NumButtons

//scp -r $GOPATH/src/github.com/Gruppe72 student@129.241.187.150:go/src/github.com

//ssh student@129.241.187.150

//Channel for knappetrykk, fra elev_driver til order_handler (ekstern)
buttonEvent = make(chan OrderButton)
//Channel for heis ved etasje, fra elev_driver til order_handler
floorEvent = make(chan int)
//Channel, fra network til order_handler
order = make(chan Order)
//Channel for remove order, fra order_handler til network 
removeOrder = make(chan Order)
//Channel for new order, fra elev_driver til network
newOrder = make(chan order)
//Channel for cost function value, fra order_dist. til network
sendCostValue = make(chan int)
//Channel, fra network til order_distribution
receiveCostValue = make(chan int)
//Channel for status of peer to peer, fra network til order_distribution
peers = make(chan int)
//Channel for motor direction, from order_distribution to FSM, FSM calls function
setMotorDir = make(chan int)




func ButtonLights(){
	for floor := 0; floor < NumFloors; floor ++{
		for button := 0; button < NumButtons; button++ {
			driver.SetButtonLamp(driver.ButtonType(button), floor, driver.GetButtonSignal(driver.ButtonType(button), floor))
		}
	}
}

func easyElev(buttonEvent chan OrderButton, floorEvent chan int){
	button := <-buttonEvent
	floor := <- floorEvent
	select {
		case button := <- buttonEvent:
			
		case floor := <- floorEvent:
	
	}
}


func doorsOpen(){
	SetDoorOpenLamp(1)
	sleep(3000* time.Milliseconds)
	SetDoorOpenLamp(0)
}

func main(){
	fmt.Printf("Elev_driver started\n")
	driver.ElevInit()
	buttonEvent := make(chan driver.OrderButton)
	floorEvent := make(chan int)
	for {
		ButtonLights()
	}
	go driver.EventListener(buttonEvent, floorEvent) 
	fmt.Printf("Elev_driver ended\n")
	
}


