package main

import (
	"fmt"
	"sync"
)

// numberServer owns the integer `i` and performs actions on it
// based on messages received on the channels.
func numberServer(inc chan struct{}, dec chan struct{}, get chan chan int, quit chan struct{}) {
	i := 0
	for {
		select {
		case <-inc:
			i++
		case <-dec:
			i--
		case resp := <-get:
			resp <- i
		case <-quit:
			return
		}
	}
}

// Incrementer sends `n` increment requests to the number server.
// It expects the caller to start it as a goroutine (e.g. `go Incrementer(...)`).
func Incrementer(wg *sync.WaitGroup, inc chan<- struct{}, n int) {
	defer wg.Done()
	for j := 0; j < n; j++ {
		inc <- struct{}{}
	}
}

// Decrementer sends `n` decrement requests to the number server.
// It expects the caller to start it as a goroutine (e.g. `go Decrementer(...)`).
func Decrementer(wg *sync.WaitGroup, dec chan<- struct{}, n int) {
	defer wg.Done()
	for j := 0; j < n; j++ {
		dec <- struct{}{}
	}
}

func main() {
	inc := make(chan struct{})
	dec := make(chan struct{})
	get := make(chan chan int)
	quit := make(chan struct{})

	var wg sync.WaitGroup
	wg.Add(2)
	// Start workers and pass the local channels and WaitGroup
	go Incrementer(&wg, inc, 1000000)
	go Decrementer(&wg, dec, 1000000)

	// Start the number-server
	go numberServer(inc, dec, get, quit)

	// Wait for workers to finish sending requests
	wg.Wait()

	// Ask the server for the final value
	resp := make(chan int)
	get <- resp
	final := <-resp

	// Tell server to quit and print result
	close(quit)
	fmt.Println("Final value of i:", final)
}
