package main

import (
	"fmt"
	"net"
	"time"
)

func sendMessage(serverIP net.IP, port int, msg string) error {
	addr := net.UDPAddr{
		IP:   serverIP,
		Port: port,
	}

	conn, err := net.DialUDP("udp", nil, &addr)
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = conn.Write([]byte(msg))
	return err
}

func main() {
	// Creates new empty IP st. we listen to ALL IP's on port 30000.
	addr := net.UDPAddr{
		IP:   net.IPv4zero, // 0.0.0.0
		Port: 30000,
	}

	// Starts listening on every IP
	conn, err := net.ListenUDP("udp", &addr)

	if err != nil {
		panic(err)
	}
	defer conn.Close()

	// Creates a buffer we will populate with data from remote server
	buffer := make([]byte, 1024)

	serverAddr := net.UDPAddr{
		IP:   net.IPv4zero,
		Port: 20000, // Must be 20000
	}

	// Fetches server IP
	for {
		n, sender, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println("Read error:", err)
			continue
		}

		message := string(buffer[:n])
		fmt.Printf("Received from %s: %s\n", sender.IP, message)
		serverAddr.IP = sender.IP
		break
		// Breaks loop after obtaining server IP
	}

	// We now have server IP, and can now handle messages any way we want
	for {
		// Sends message to server
		time.Sleep(1 * time.Second)
		sendMessage(serverAddr.IP, serverAddr.Port, "yalla habibi")
		fmt.Printf("Tried sending to %s\n", &serverAddr.IP)
	}
}
