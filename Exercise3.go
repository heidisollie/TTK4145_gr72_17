
package main

import (
	"fmt"
	"net"
	//"log"
	"time"
	//"strconv"
)


const bcast1 = ":30000"
const bcast2 = "129.241.187.255:20002"
const udpPort = "20002"



func CheckError (err error){
	if err != nil {
		fmt.Println("Error: ", err)
	}
}



func main(){
	InternetAdr, err := net.ResolveUDPAddr("udp", bcast1)
	CheckError(err)
	
	ClientAdr, err := net.ResolveUDPAddr("udp", bcast2)
	CheckError(err)
	
	
	Conn1, err := net.DialUDP("udp", InternetAdr, ClientAdr)
	CheckError(err)
	
	//Conn2, err := net.ListenUDP("udp", InternetAdr)
	//CheckError(err)
	
	defer Conn1.Close()
	//defer Conn2.Close()
	

	for {

		ClientBuf := make([]byte, 1024)
		Message := []byte("Hello")
		
		fmt.Printf("Hello\n")
		_, err := Conn1.Write(Message)
		CheckError(err)
	
		
		time.Sleep(1000 * time.Millisecond)
		
		
		n, addr, err := Conn1.ReadFrom(ClientBuf)
		CheckError(err)
		fmt.Printf("Message received from %s : %s",addr, string(ClientBuf[0:n]))
		
		

	
	
	}
	
	
}
