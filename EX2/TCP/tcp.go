package main

import (
	"flag"
	"fmt"
	"net"
	"time"
)

const (
	workspaceNumber = 8
	myIP            = ""
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

// getLocalIP returns the client's local IP address by connecting to the server
func getLocalIP(serverIP string) string {
	connection, err := net.Dial("udp", serverIP+":80")
	if err != nil {
		fmt.Println("Error getting local IP:", err)
		return "127.0.0.1"
	}
	defer connection.Close()
	localAddr := connection.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String()
}

// sendMessages sends messages to the server and reads responses on the same connection
func sendMessages(serverIP string) {
	localIP := getLocalIP(serverIP)
	listenPort := 12345

	serverAddr := fmt.Sprintf("%s:33546", serverIP)
	fmt.Printf("Connecting to server at %s\n", serverAddr)

	tcpAddr, err := net.ResolveTCPAddr("tcp", serverAddr)
	if err != nil {
		fmt.Println("Error resolving server address:", err)
		return
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}
	defer conn.Close()

	// Read welcome message
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err == nil && n > 0 {
		fmt.Printf("[Welcome] %s\n", string(buffer[:n]))
	}

	// Send connect message
	connectMsg := fmt.Sprintf("Connect to: %s:%d\x00", localIP, listenPort)
	_, err = conn.Write([]byte(connectMsg))
	if err != nil {
		fmt.Println("Error sending connect message:", err)
		return
	}
	fmt.Printf("Sent: %s\n", connectMsg)

	// Read echo for connect message
	n, err = conn.Read(buffer)
	if err == nil && n > 0 {
		fmt.Printf("Echo: %s\n", string(buffer[:n]))
	}

	// Send messages with sleep to be nice to the network
	for i := 0; i < 10; i++ {
		message := fmt.Sprintf("Hello from client #%d (message %d)\x00", workspaceNumber, i)
		_, err := conn.Write([]byte(message))
		if err != nil {
			fmt.Println("Error sending:", err)
			continue
		}
		fmt.Printf("Sent: %s\n", message)

		// Read echo for each message
		n, err := conn.Read(buffer)
		if err == nil && n > 0 {
			fmt.Printf("Echo: %s\n", string(buffer[:n]))
		}

		time.Sleep(1 * time.Second) // Don't spam
	}
}

// listenForResponses listens on a local port for server callback connections
func listenForResponses() {
	listenPort := 12345

	localAddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf(":%d", listenPort))
	if err != nil {
		fmt.Println("Error resolving local address:", err)
		return
	}

	listener, err := net.ListenTCP("tcp", localAddr)
	if err != nil {
		fmt.Println("Error listening for callbacks:", err)
		return
	}
	defer listener.Close()

	fmt.Printf("Listening for callbacks on port %d\n", listenPort)

	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		go func(c *net.TCPConn) {
			defer c.Close()
			fmt.Printf("[Callback] Connection from %s\n", c.RemoteAddr())

			buffer := make([]byte, 1024)
			for {
				n, err := c.Read(buffer)
				if err != nil {
					fmt.Println("[Callback] Connection closed")
					return
				}
				message := string(buffer[:n])
				fmt.Printf("Callback from %s: %s\n", c.RemoteAddr().String(), message)
			}
		}(conn)
	}
}
