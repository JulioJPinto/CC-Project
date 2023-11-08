package fstp

import (
	"fmt"
	"log"
	"net"
)

// fstp_routine ...
type fstp_routine struct {
	conn    net.Conn
	handler FSTPHandler
}

type FSTPHandler interface {
	HandleRequest(net.Conn, FSTPRequest) FSTPresponse
	HandleShutdown(net.Conn, error)
}

// FSTPServer ...
type FSTPServer struct {
	host    string
	port    string
	handler FSTPHandler
}

// New ...
func New(config *Config, handler FSTPHandler) *FSTPServer {
	return &FSTPServer{
		host:    config.Host,
		port:    config.Port,
		handler: handler,
	}
}

// Run ...
func (server *FSTPServer) Run() {
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
		header := make([]byte, FSTPHEaderSize)
		n, err := instance.conn.Read(header)
		if n != FSTPHEaderSize || err != nil {
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

		fmt.Printf("header: %x payload: %s \n", buffer[0:FSTPHEaderSize], buffer[FSTPHEaderSize-1:])

		req_msg := FSTPmessage{}
		req_msg.Deserialize(buffer)
		req := FSTPRequest(req_msg)
		resp := instance.handler.HandleRequest(instance.conn, req)

		resp_msg := FSTPmessage(resp)
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
