package main

import (
	"fmt"
	"net"
	"time"
)

func fetchIP() net.IP {
	// GETS IP USING UDP
	// Creates new empty IP st. we listen to ALL IP's on port 30000.
	addrUDP := net.UDPAddr{
		IP:   net.IPv4zero, // 0.0.0.0
		Port: 30000,
	}

	// Starts listening on every IP
	conn, err := net.ListenUDP("udp", &addrUDP)

	if err != nil {
		panic(err)
	}
	defer conn.Close()

	// Creates a buffer we will populate with data from remote server
	buffer := make([]byte, 1024)

	// Fetches server IP
	for {
		n, sender, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println("Read error:", err, n)
			continue
		}

		return sender.IP
	}
}

func main() {

	serverIP := fetchIP().String()
	serverPort := "34933" // Forces fixed-size data of 1024

	conn, err := net.Dial("tcp", serverIP+":"+serverPort)
	if err != nil {
		fmt.Println("Connection error: ", err)
		return
	}
	defer conn.Close()

	fmt.Println(("Connected to server"))

	for {
		// Prepares the 1024 byte data to be sent
		msg := make([]byte, 1024)
		copy(msg, []byte("Inshallah"))

		_, err = conn.Write(msg)
		if err != nil {
			fmt.Println("Write error: ", err)
			return
		}

		// Reads the 1024 byte echo from the server
		buf := make([]byte, 1024)

		echo, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Read error: ", err)
			return
		}
		fmt.Println("Echo: ", echo)
		time.Sleep(1 * time.Second)
	}

}
