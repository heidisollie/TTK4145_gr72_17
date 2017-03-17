package backup

import (
	def "../definitions"
	"../network/bcast"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"time"
)

func WriteQueueToFile(queue []def.Order, filename string) {
	data, err := json.Marshal(queue)
	ioutil.WriteFile(filename, data, 0644)
	if err != nil {
		log.Fatal(data, err)
	}
}

func ReadQueueFromFile(queue *[]def.Order, filename string) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal("error:: ", err)
	} else {
		err = json.Unmarshal(data, &queue)
	}
}

func AliveSpammer(filename string) {
	aliveChannel := make(chan string)
	alivePeriod := 50 * time.Millisecond
	alivePort := 37718
	msg := "alive"
	aliveTicker := time.NewTicker(alivePeriod)

	go bcast.LocalTransmitter(alivePort, aliveChannel)
	for {
		<-aliveTicker.C
		aliveChannel <- msg
	}
}

func BackUp() {
	alivePeriod := 200 * time.Millisecond
	alivePort := 37718
	var resetChannel = make(chan string)
	isAliveTimer := time.NewTimer(alivePeriod)
	isAliveTimer.Stop()

	go bcast.Receiver(alivePort, resetChannel)

	for {
		select {
		case msg := <-resetChannel:
			if msg == "alive" {
				isAliveTimer.Reset(alivePeriod)
			}

		case <-isAliveTimer.C:
			fmt.Printf("Program timed out, backup starting\n")
			return
		}
	}
}
