package fstp

import (
	"fmt"
	"log"
	"net"
)

// fstp_routine ...
type fstp_routine struct {
	conn    net.Conn
	handler Handler
}

type Handler interface {
	HandleRequest(net.Conn, Request) Response
	HandleShutdown(net.Conn, error)
}

// Server ...
type Server struct {
	host    string
	port    string
	handler Handler
}

// NewServer ...
func NewServer(config *Config, handler Handler) *Server {
	return &Server{
		host:    config.Host,
		port:    config.Port,
		handler: handler,
	}
}

// Run ...
func (server *Server) Run() {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", server.host, server.port))
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}

		instance := &fstp_routine{
			conn:    conn,
			handler: server.handler,
		}
		go instance.handleClient()
	}
}

const buffer_limit = 1024

func (instance *fstp_routine) handleClient() {
	var err error

	defer instance.handler.HandleShutdown(instance.conn, err)
	defer instance.conn.Close()

	fmt.Println("Accepted connection from", instance.conn.RemoteAddr())

	for {
		header := make([]byte, HeaderSize)
		n, err := instance.conn.Read(header)
		if n != HeaderSize || err != nil {
			fmt.Println("Error reading header", err)
			return
		}
		payload_size := PayloadSize(header)
		buffer := make([]byte, payload_size) // Create a buffer to store incoming data
		_, err = instance.conn.Read(buffer)
		if err != nil {
			fmt.Println("Error reading:", err)
			return
		}

		buffer = append(header, buffer...)

		fmt.Printf("header: %x payload: %s \n", buffer[0:HeaderSize], buffer[HeaderSize-1:])

		req_msg := Message{}
		req_msg.Deserialize(buffer)
		req := Request(req_msg)
		resp := instance.handler.HandleRequest(instance.conn, req)

		resp_msg := Message(resp)
		fmt.Printf("response payload: %s \n", resp.Payload)

		response, err := resp_msg.Serialize()
		if err != nil {
			fmt.Println("Error serializing:", err)
			return
		}
		_, err = instance.conn.Write(response)
		if err != nil {
			fmt.Println("Error writing:", err)
			return
		}
	}

}
