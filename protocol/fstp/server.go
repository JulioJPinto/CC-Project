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
	HandleRequest(FSTPrequest) FSTPresponse
	HandleIHaveRequest(net.Conn, FSTPHeader, IHaveProps) FSTPresponse
}

// FSTPServer ...
type FSTPServer struct {
	host    string
	port    string
	handler FSTPHandler
}

// New ...
func New(config *FSTPConfig, handler FSTPHandler) *FSTPServer {
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
		fmt.Printf("header: %x payload: %s \n", recieved_data[0:FSTPHEaderSize], recieved_data[FSTPHEaderSize:])
		// FSTPHeader
		msg_type := MessageType(recieved_data)

		switch msg_type {
		case IHave:
			req_payload := IHaveProps{}
			err := (&req_payload).Deserialize(recieved_data[FSTPHEaderSize:])
			fmt.Println("erro:", err)
			fmt.Println("req_payload", req_payload)
			instance.handler.HandleIHaveRequest(instance.conn, FSTPHeader{msg_type}, req_payload)
		}
		// if err != nil {
		// 	fmt.Println("Error serializing:", err)
		// 	return
		// }
		// json.Unmarshal(recieved_data, &req_msg)
		// fmt.Println("dá erro?? ", err)
		// fmt.Println("msg", req_msg.Payload)
		// req := FSTPrequest(req_msg)
		// resp := instance.handler.HandleRequest(req)

		// resp_msg := FSTPmessage(resp)
		// response, err := resp_msg.Serialize()
		// _, err = instance.conn.Write(response)
		// if err != nil {
		// 	fmt.Println("Error writing:", err)
		// 	return
		// }
	}
}
