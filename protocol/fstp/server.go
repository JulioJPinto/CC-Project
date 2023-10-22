package fstp

import (
	"log"
	"fmt"
	"net"
)
// fstp_instance ...
type fstp_instance struct {
	conn net.Conn
	handler FSTPHandler
}

type FSTPHandler interface {
	HandleRequest(FSTPrequest) FSTPresponse
} 

// FSTPServer ...
type FSTPServer struct {
	host string
	port string
	handler FSTPHandler
}


// New ...
func New(config *FSTPConfig,handler FSTPHandler) *FSTPServer {
	return &FSTPServer{
		host: config.Host,
		port: config.Port,
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

		instance := &fstp_instance{
			conn: conn,
			handler: server.handler,
		}
		go instance.handleClient()
	}
}

const buffer_limit = 1024

func (instance *fstp_instance) handleClient() {
	defer instance.conn.Close()
	fmt.Println("Accepted connection from", instance.conn.RemoteAddr())

	var recieved_data []byte
	buffer := make([]byte, buffer_limit) // Create a buffer to store incoming data

	
	for {
		for {
			n, err := instance.conn.Read(buffer)
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

		req := FSTPrequest{}
		req.Deserialize(recieved_data)
		resp := instance.handler.HandleRequest(req)
		// If you want to send a response, you can use conn.Write
		response,_ := resp.Serialize()
		_, err := instance.conn.Write(response)
		if err != nil {
			fmt.Println("Error writing:", err)
			return
		}
	}
}
