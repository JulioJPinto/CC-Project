package p2p

import (
	"fmt"
	"net"
)

type Handler interface {
	HandleRequest(Request) Response
}

type Listener struct {
	host string
	port string
}

type p2p_routine struct {
	conn    net.Conn
	handler Handler
}

func P2PConn(config Config) *Listener {
	return &Listener{host: config.Host, port: config.Port}
}

func ListenOnUDP(host string,port string) error {

	serverAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%s", host, port))
	if err != nil {
		return err
	}
	conn, err := net.ListenUDP("udp", serverAddr)
	if err != nil {
		return err
	}
	defer conn.Close()
	buffer := make([]byte, PacketSize)
	for {
		n, addr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			return err
		}
		fmt.Printf("just read %d bytes from %s: %v\n", n, addr.String(), buffer[:n])
	}
}

func (instance *p2p_routine) handlePacket(addr *net.UDPAddr, buffer []byte) {
	fmt.Println(instance.conn.RemoteAddr())
	fmt.Printf("Boutta handle all this %v\n", buffer)
	var req Message
	req.Deserialize(buffer)
	if req.Header.is_request {
		instance.handler.HandleRequest(Request(req))
	} // else {
	// 	instance.handler.HandleResponse(Response(req))
	// }
}
