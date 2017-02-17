package main

import (
	"fmt"
	"runtime"
//	"time"
)

var i int = 0;


func thread_1(){

	for j:= 0; j<100; j++{	
		temp := <- channel
		temp++
		channel <- temp
		fmt.Printf("A: %v\n", temp)
	}
	done <- 1
}

func thread_2(){
	for n:= 0; n<100; n++{
		temp := <- channel
		temp--
		channel <- temp
		fmt.Printf("B: %v\n", temp)
	}

}




func main(){
	runtime.GOMAXPROCS(runtime.NumCPU()) //Gjør det mulig å kjøre trådene parallelt
	channel <- i
	
	go thread_1()
	go thread_2()
	go thread_3()
	
	<- done
	<- done
	<- done
	
	fmt.Println(<- channel)
}
