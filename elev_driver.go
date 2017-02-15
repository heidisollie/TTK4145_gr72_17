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

//Channel for knappetrykk, fra elev_driver til order_handler
buttonEvent = make(chan OrderButton)
//Channel for heis ved etasje, fra elev_driver til order_handler
floorEvent = make(chan int)
//Channel for  ---, fra order_dist. til elev_driver
motorDir = make(chan int)
//Channel, fra network til order_handler
order = make(chan Order)
//Channel, fra network til order_distribution
costValue = make(chan int)
//Channel for motor direction, fra FSM til elev_driver
motorDir = make(chan int)
//Channel for open/close door, fra FSM til elev_driver
doorState = make(chan int)
//Channel for motor direction, from order_distribution to FSM
motorDir = make(chan int)
//Channel for timer, from elev_driver to FSM
timer = make(chan int)




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
	for floor := 0; floor < NumFloors; floor ++{
		for button := 0; button < NumButtons; button++ {
			buttonSignal = driver.GetButtonSignal(driver.ButtonType(button), floor)
			if (buttonSignal){
				
			}
		}
	}
}


func doorsOpen(){
	SetDorrOpenLamp(1)
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


