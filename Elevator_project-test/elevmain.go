package main

import (
	"./FSM"
	"./backup"
	def "./definitions"
	"./driver"
	"./network"
	"./network/peers"
	"./orderDistribution"
	"./orderHandler"
	"fmt"
	"log"
	"os"
	//"os/exec"
)

func main() {
	fmt.Printf("Elevator driver started\n")
	driver.ElevInit()
	var localIP string = network.GetIP()

	//----------------------------------CHANNELS----------------------------------------------------------------
	newTargetFloor := make(chan int, 100)
	floorCompleted := make(chan int, 100)

	floorEvent := make(chan int, 100)
	floorEventFSM := make(chan int, 100)
	floorEventOrderHandler := make(chan int, 100)

	buttonEvent := make(chan def.OrderButton, 100)
	processNewOrder := make(chan def.Order, 100)
	assignedNewOrder := make(chan def.Order, 100)

	elevSendState := make(chan def.Elevator, 100)
	elevSendNewOrder := make(chan def.Order, 100)
	elevSendRemoveOrder := make(chan def.Order, 100)

	elevReceiveState := make(chan def.Elevator, 100)
	elevReceiveNewOrder := make(chan def.Order, 100)
	elevReceiveRemoveOrder := make(chan def.Order, 100)

	elevatorLost := make(chan string, 100)
	peers := make(chan peers.PeerUpdate, 100)
	//-----------------------------------------------------------------------------------------------------------

	driver.ElevInit()

	//Creating backup that waits
	if _, err := os.Open(def.Filename); err == nil {
		fmt.Printf("Backup waiting\n")
		//backup.BackUp()
		fmt.Printf("Backup starting\n")
	} else {
		fmt.Printf("First\n")
		if _, err := os.Create(def.Filename); err != nil {
			log.Fatal("Cannot create a file\n")
		}
	}

	//Spawning backup terminal
	//backupTerminal := exec.Command("gnome-terminal", "-x", "sh", "-c", "go run elevmain.go")
	//backupTerminal.Run()

	go backup.AliveSpammer(def.Filename)

	go driver.DriverInit(buttonEvent, floorEvent)

	go FSM.FSMInit(floorEventFSM, newTargetFloor, floorCompleted, elevSendState)

	go orderHandler.OrderHandlerInit(localIP,
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

	go orderDistribution.OrderDistInit(localIP,
		processNewOrder,
		assignedNewOrder,
		elevReceiveState,
		elevatorLost,
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
