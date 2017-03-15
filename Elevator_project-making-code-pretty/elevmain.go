package main

import (
	"./FSM"
	"./backup"
	"./driver"
	"./network"
	"./network/bcast"
	"./network/peers"
	"./order_distribution"
	"./order_handler"
	def "./definitions"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"
)


func main() {
	fmt.Printf("Elev_driver started\n")

	driver.ElevInit()

	const NumFloors = driver.NumFloors
	const NumButtons = driver.NumButtons

	var localIP string = network.GetIP()

	//----------------------------------CHANNELS----------------------------------------------------------------
	newTargetFloor := make(chan int)						
	floorCompleted := make(chan int)				

	floorEvent := make(chan int) 							
	floorEventFSM := make(chan int)							
	floorEventOrderHandler := make(chan int)				

	buttonEvent := make(chan driver.OrderButton) 			
	processNewOrder := make(chan def.Order)  			
	assignedNewOrder := make(chan def.Order) 			

	elevSendState := make(chan def.ElevState) 			
	elevSendNewOrder := make(chan def.Order) 			
	elevSendRemoveOrder := make(chan def.Order) 		

	elevReceiveState := make(chan def.ElevState)		
	elevReceiveNewOrder := make(chan def.Order)			
	elevReceiveRemoveOrder := make(chan def.Order)	

	elevatorLost := make(chan string)							

	peers := make(chan peers.PeerUpdate)					
	//-----------------------------------------------------------------------------------------------------------

	driver.ElevInit()			

	//Creating backup that waits
	if _, err := os.Open(def.Filename); err == nil {
		fmt.Printf("Backup waiting\n")
		backup.BackUp()
		fmt.Printf("Backup starting\n")
	} else {
		fmt.Printf("First\n")
		if _, err := os.Create(def.Filename); err != nil {
			log.Fatal("Cannot create a file\n")
		}

		go order_handler.OrderHandlerInit(localIP,
			floorCompleted,
			buttonEvent,
			assignedNewOrder,
			processNewOrder,
			elevSendNewOrder,
			elevSendRemoveOrder,
			elevReceiveNewOrder,
			elevReceiveRemoveOrder,
			elevatorLost,
			newTargetFloor,
			floorEventOrderHandler)
	}

	//Spawning backup terminal
	backupTerminal := exec.Command("gnome-terminal", "-x", "sh", "-c", "go run elevmain.go")
	backupTerminal.Run()


	go backup.AliveSpammer(def.Filename)

	go driver.driverInit(buttonEvent, floorEvent)

	go FSM.FSMInit(floorEventFSM, newTargetFloor, floorCompleted, elevSendState)

	go order_handler.OrderHandlerInit(localIP,
		floorCompleted,
		buttonEvent,
		assignedNewOrder,
		processNewOrder,
		elevSendNewOrder,
		elevSendRemoveOrder,
		elevReceiveNewOrder,
		elevReceiveRemoveOrder,
		elevatorLost,
		newTargetFloor,
		floorEventOrderHandler)

	go order_distribution.OrderDistInit(localIP,
		processNewOrder,
		assignedNewOrder,
		elevReceiveState,
		elevLost,
		peers)

	go network.NetworkInit(localIP,
		elevSendState,
		elevSendNewOrder,
		elevSendRemoveOrder,
		elevReceiveState,
		elevReceiveNewOrder,
		elevReceiveRemoveOrder,
		peers)

	//Sends information from floorEvent to two channels
	go network.Repeater(floorEvent, floorEventFSM, floorEventOrderHandler)

	select {}
}

