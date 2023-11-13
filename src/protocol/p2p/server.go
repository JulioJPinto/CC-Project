package p2p

import (
	"fmt"
	"net"
)

type Header struct {
}

type Request Message
type Response Message

type Message struct {
	Header  Header
	Payload []byte //JSON Serializable
}


type Handler interface {
	HandleRequest(net.Conn, Request) Response
}

type Server struct {
	host string
	port string
	handler Handler
}

const P2P_PacketSize = 2048

type p2p_routine struct {
	conn net.Conn
	handler Handler
}

func NewP2PServer(config Config,handler Handler) *Server {
	return &Server{host: config.Host, port: config.Port,handler: handler}

}

func (server *Server) Run() error {
	serverAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%s", server.host, server.port))
	if err != nil {
		return err
	}
	conn, err := net.ListenUDP("udp", serverAddr)
	if err != nil {
		return err
	}
	defer conn.Close()
	buffer := make([]byte, P2P_PacketSize)
	for {
		n, addr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			return err
		}
		fmt.Printf("just read %d bytes from %s: %v\n", n, addr.String(), buffer[:n])

		instance := &p2p_routine{conn: conn}
		go instance.handleClient(addr, buffer[:n])	 	
	}

}

func (instance *p2p_routine) handleClient(addr *net.UDPAddr, buffer []byte) {
	fmt.Println(instance.conn.RemoteAddr())
	fmt.Printf("Boutta handle all this %v\n", buffer)
}
