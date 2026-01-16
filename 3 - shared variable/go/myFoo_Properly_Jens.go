// Use `go run foo.go` to run your program

package main

import (
	. "fmt"
	"runtime"
	"time"
)

var i = 0

func incrementing() {
	//TODO: increment i 1000000 times
	for k:=0; k < 1000000; k++ {
		i++
	}
	quit <- 0
}

func decrementing() {
	//TODO: decrement i 1000000 times
	for k:=0; k < 1000000; k++ {
		i--
	}
	quit <- 0
}

func main() {
	// What does GOMAXPROCS do? What happens if you set it to 1?
	// This controls the max number of parallell threads, which must be greater than 1 to
	// achieve the "desired" error magic number
	runtime.GOMAXPROCS(2)

	// Making two channels
	ch := make(chan int)
	quit := make(chan int) // when two 0's is sent, we are done

	
	// TODO: Spawn both functions as goroutines
	// This is GO's way of starting a new goroutine, which is equivalent to a thread.
	go incrementing()
	go decrementing()


	// We have no direct way to wait for the completion of a goroutine (without additional synchronization of some sort)
	// We will do it properly with channels soon. For now: Sleep.
	time.Sleep(500*time.Millisecond)
	Println("The magic number is:", i)
}
