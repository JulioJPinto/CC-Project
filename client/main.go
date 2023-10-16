package main

import (
	"fmt"
	"net"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}
	defer conn.Close()

	fmt.Println("Connected to server at localhost:8080")

	request(conn)
	// Send and receive data with the server
}
func request(conn net.Conn) {

	request := "This is a dummy request"
	_, err := conn.Write([]byte(request))
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	fmt.Println("Sent request:", request)

	// Receive and print the server's response
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Error receiving response:", err)
		return
	}
	response := string(buffer[:n])
	fmt.Println("Received response:", response)
}
