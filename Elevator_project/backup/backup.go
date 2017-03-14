package backup

import (
	"../structs"
	"encoding/json"
	//"fmt"
	"../network/bcast"
	//"../network"
	"io/ioutil"
	"log"
	"time"
	//"os"
)

func WriteQueueToFile(queue []structs.Order, filename string) {
	data, err := json.Marshal(queue)
	ioutil.WriteFile(filename, data, 0644)
	if err != nil {
		log.Fatal(data, err)
	}
}

func ReadQueueFromFile(queue *[]structs.Order, filename string) {
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
