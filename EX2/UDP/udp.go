package main

import (
	"flag"
	"fmt"
	"net"
	"time"
)

const (
	workspaceNumber = 8
)

var wfh = flag.Bool("wfh", false, "Work from home mode")

// Config determines if running in work-from-home mode
func isWorkFromHome() bool {
	return *wfh
}

func main() {
	flag.Parse()

	var serverIP string

	if isWorkFromHome() {
		// WFH mode: use localhost (server is on same machine)
		serverIP = "127.0.0.1"
		fmt.Printf("Work-from-home mode: using server at %s\n", serverIP)
	} else {
		// Lab mode: discover server via broadcast
		serverIPChan := make(chan string)
		go receiveServerIP(serverIPChan)

		serverIP = <-serverIPChan
		fmt.Printf("Found server at: %s\n", serverIP)
	}

	// Start sending goroutine
	go sendMessages(serverIP)

	// Keep main alive and listen for server responses
	listenForResponses()
}

// receiveServerIP listens on port 30000 for the server's broadcast
func receiveServerIP(serverIPChan chan string) {
	addr, err := net.ResolveUDPAddr("udp", "0.0.0.0:30000")
	if err != nil {
		fmt.Println("Error resolving address:", err)
		return
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Println("Error listening:", err)
		return
	}
	defer conn.Close()

	fmt.Println("Listening for server broadcast on port 30000...")

	buffer := make([]byte, 1024)
	for {
		n, remoteAddr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println("Error reading:", err)
			continue
		}

		message := string(buffer[:n])
		fmt.Printf("Received broadcast from %s: %s\n", remoteAddr.IP.String(), message)

		// Send the server IP through channel (only once)
		select {
		case serverIPChan <- remoteAddr.IP.String():
			fmt.Printf("Server IP sent to channel: %s\n", remoteAddr.IP.String())
			return
		default:
			// Channel already has a value
			return
		}
	}
}

// sendMessages sends messages to the server and listens for responses
func sendMessages(serverIP string) {
	var port string
	if isWorkFromHome() {
		port = "20000" // WFH server uses port 20000
		fmt.Println("Work-from-home mode: sending to port 20000")
	} else {
		port = fmt.Sprintf("%d", 20000+workspaceNumber) // Lab uses 20000+n
		fmt.Printf("Lab mode: sending to port %s\n", port)
	}
	serverAddr := fmt.Sprintf("%s:%s", serverIP, port)

	fmt.Printf("Connecting to server at %s\n", serverAddr)

	udpAddr, err := net.ResolveUDPAddr("udp", serverAddr)
	if err != nil {
		fmt.Println("Error resolving server address:", err)
		return
	}

	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		fmt.Println("Error dialing:", err)
		return
	}
	defer conn.Close()

	// Send messages with sleep to be nice to the network
	for i := 0; i < 10; i++ {
		message := fmt.Sprintf("Hello from client #%d (message %d)", workspaceNumber, i)
		_, err := conn.Write([]byte(message))
		if err != nil {
			fmt.Println("Error sending:", err)
			continue
		}

		fmt.Printf("Sent: %s\n", message)
		time.Sleep(1 * time.Second) // Don't spam
	}
}

// listenForResponses listens on a local port for server responses
func listenForResponses() {
	var port int
	if isWorkFromHome() {
		port = 20001 // WFH server replies on port 20001
	} else {
		port = 20000 + workspaceNumber // Lab replies on same port
	}

	localAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%d", port))
	if err != nil {
		fmt.Println("Error resolving local address:", err)
		return
	}

	conn, err := net.ListenUDP("udp", localAddr)
	if err != nil {
		fmt.Println("Error listening for responses:", err)
		return
	}
	defer conn.Close()

	fmt.Printf("Listening for responses on port %d\n", port)

	buffer := make([]byte, 1024)
	for {
		n, remoteAddr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println("Error reading response:", err)
			continue
		}

		message := string(buffer[:n])
		fmt.Printf("Response from %s: %s\n", remoteAddr.IP.String(), message)
	}
}
