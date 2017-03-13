package main

import (
	"./FSM"
	"./driver"
	"./network"
	"./network/peers"
	"./order_distribution"
	"./order_handler"
	"./structs"
	"fmt"
)

//------------Oppdatering------------
// Heisen mister fortsatt noen bestillinger, hvis man trykker for fort
// Vi vet ikke hvorfor dette skjer

// Burde forsøke å gå bort fra new_target_floor tenking og heller bare velge retning
// Problemet da er at vi bare har en funksjon som sjekker other_orders_in_dir
// Den plukker ikke opp ordre som ikke "passer med dette" men den skal jo fortsatt ta første ordre i køen
// Den må fortsatt ha et mål

//Fikse backup slik at terminal åpner seg igjen hvis man dreper programmet.
//Da skal den åpen seg med de aktuelle ordrene (command button)

// X Miste nettverk, eksterne ordre må delegeres til andre heiser
// X Deklareres stuck, eksterne ordre må delegeres til andre heiser

// X Miste nettverk, samme greia med da må order_dist få tilgang til IP adressen til mistet heis

// -X -- Lage funksjon som restribuerer

//Må teste packetloss? Mulig å implementere dette kjapt, eller er det bare noen poeng som mistes

//Teste at order_distribution fungerer som den skal
//Teste at tre heiser kjører selvstendig når nettverket går bort

//Når heisen er i idle og får en bestilling men ikke klarer å bevege seg skal den deklareres stuck.

func main() {
	fmt.Printf("Elev_driver started\n")
	driver.ElevInit()

	const NumFloors = driver.NumFloors
	const NumButtons = driver.NumButtons

	var localIP string
	localIP = network.GetIP()

	newTargetFloor := make(chan int, 100)
	floorEvent := make(chan int, 100)                 //heis ved etasje til order_handler
	buttonEvent := make(chan driver.OrderButton, 100) //knappetrykk til order_handler
	processNewOrder := make(chan structs.Order, 100)  //ekstern ordre fra order handler for kost funksjonen
	assignedNewOrder := make(chan structs.Order, 100) //ekstern ordre fra order_dist. , med heisID

	elevSendState := make(chan structs.ElevState, 100)
	elevSendNewOrder := make(chan structs.Order, 100)
	elevSendRemoveOrder := make(chan structs.Order, 100)

	elevReceiveState := make(chan structs.ElevState, 100)
	elevReceiveNewOrder := make(chan structs.Order, 100)
	elevReceiveRemoveOrder := make(chan structs.Order, 100)

	elevLost := make(chan string, 100)

	floorCompleted := make(chan int, 100)
	peers := make(chan peers.PeerUpdate, 100)

	driver.ElevInit()
	go driver.EventListener(buttonEvent, floorEvent)
	go FSM.FSMInit(floorEvent, newTargetFloor, floorCompleted, elevSendState)

	go order_handler.OrderHandlerInit(localIP,
		floorCompleted,
		buttonEvent,
		assignedNewOrder,
		processNewOrder,
		elevSendNewOrder,
		elevSendRemoveOrder,
		elevReceiveNewOrder,
		elevReceiveRemoveOrder,
		elevLost,
		newTargetFloor)

	go network.NetworkInit(localIP,
		elevSendState,
		elevSendNewOrder,
		elevSendRemoveOrder,
		elevReceiveState,
		elevReceiveNewOrder,
		elevReceiveRemoveOrder,
		peers)

	go order_distribution.OrderDistInit(localIP,
		processNewOrder,
		assignedNewOrder,
		elevReceiveState,
		elevLost,
		peers)

	select {}
}
