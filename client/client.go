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

    // Send and receive data with the server
}