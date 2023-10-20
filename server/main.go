package main

import (
	"cc_project/protocol/client_tracker"
	"fmt"
	"net"
)

const buffer_limit = 8

func main() {
	listener, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Error listening:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Server is listening on localhost:8080")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	defer conn.Close()
	fmt.Println("Accepted connection from", conn.RemoteAddr())

	var recieved_data []byte
	buffer := make([]byte, buffer_limit) // Create a buffer to store incoming data

	for {
		for {
			n, err := conn.Read(buffer)
			if err != nil {
				fmt.Println("Error reading:", err)
				return
			}

			recieved_data = append(recieved_data, buffer[:n]...)
			if n != buffer_limit {
				break
			}
		}
		fmt.Println("Received:", string(recieved_data))
		handleRequest(recieved_data, conn.RemoteAddr())
		// If you want to send a response, you can use conn.Write
		response := []byte("Hello from the server")
		_, err := conn.Write(response)
		if err != nil {
			fmt.Println("Error writing:", err)
			return
		}
	}
}

func handleRequest(data []byte, remote net.Addr) {
	header := client_tracker.Deserialize(data)
	fmt.Println(remote, header)
}
