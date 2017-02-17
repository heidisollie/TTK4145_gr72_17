package main

import (
	"fmt"
	"time"
	"os/exec"
	"../network/bcast"
)

const udpPort int = 20002
//const otherIP = "129.241.187.255"
	
	
type Counter struct{
	State int
}

type Message struct{
	Data int
}

func main(){
	counter := Counter{0}
	fmt.Print("We are counting\n")
	
	spawnBackup := exec.Command("gnome-terminal", "-x", "sh", "-c", "go run backup.go")
	
	spawnBackup.Start()
	
	toBackup := make(chan Message, 1)
	go bcast.Transmitter(udpPort, toBackup)
	
	for {
		fmt.Printf("Counter: %d \n", counter.State)
		msg := Message{counter.State}
		toBackup <- msg
		counter.State++
		time.Sleep(1*time.Second)
	
	}
	
} 


