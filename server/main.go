package main

import (
	"cc_project/server/db"
	"fmt"
	"net"
)

func main() {

	var database = db.NewJSONDatabase("db.json")
	fmt.Println(database.Connect())

	ip := net.ParseIP("127.0.0.1")
	a := db.DeviceData{ip}
	database.RegisterDevice(a)
	database.Close()

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

	buffer := make([]byte, 8) // Create a buffer to store incoming data

	for {
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Error reading:", err)
			return
		}

		// Process the received data
		payload := buffer[:n]
		fmt.Println("Received:", string(payload))

		// If you want to send a response, you can use conn.Write
		response := []byte("Hello from the server")
		_, err = conn.Write(response)
		if err != nil {
			fmt.Println("Error writing:", err)
			return
		}
	}
}
