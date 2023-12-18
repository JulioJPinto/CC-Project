package fstp

import (
	"fmt"
	"log"
	"net"

	"github.com/fatih/color"
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
	debug   bool
}

// NewServer ...
func NewServer(config *Config, handler Handler, debug bool) *Server {
	return &Server{
		host:    config.Host,
		port:    config.Port,
		handler: handler,
		debug:   debug,
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
		go instance.handleClient(server.debug)
	}
}

const buffer_limit = 1024

func (instance *fstp_routine) handleClient(debugging bool) {
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

		if debugging {
			s_data := fmt.Sprint("\nsending: ", buffer)
			s_str := fmt.Sprint("\nAKA \n\t<", HeaderType(int(buffer[0])), ">\n\tPayload: ", string(buffer[HeaderSize:]))
			color.Green(s_data)
			color.Blue(s_str)
		}

		req_msg := Message{}
		req_msg.Deserialize(buffer)
		req := Request(req_msg)
		resp := instance.handler.HandleRequest(instance.conn, req)

		resp_msg := Message(resp)
		// fmt.Printf("response payload: %s \n", resp.Payload)

		response, err := resp_msg.Serialize()
		if err != nil {
			fmt.Println("Error serializing:", err)
			return
		}
		if debugging {

			str := fmt.Sprint("\nAKA: \n\t<", HeaderType(int(response[0])), ">\n\tPayload: ", string(response[HeaderSize:]))
			color.Blue(str)
		}
		_, err = instance.conn.Write(response)
		if err != nil {
			fmt.Println("Error writing:", err)
			return
		}
	}

}
