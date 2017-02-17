package main

import "time"

func main() {
    timer := time.NewTimer(time.Second * 2)
    <- timer.C
    println("Timer expired")
}

