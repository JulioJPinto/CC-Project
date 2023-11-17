package lib

import (
	"cc_project/node/p2p"
	"fmt"
	"net"
)

type Handler struct{}

func (*Handler) HandleRequest(p2p.Request) p2p.Response {
	return p2p.Response(p2p.Message{})
}

func (g *Gaijo) ListenOnUDP() error {
	serverAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%s", g.P2PConfig.Host, g.P2PConfig.Port))
	if err != nil {
		return err
	}
	conn, err := net.ListenUDP("udp", serverAddr)
	if err != nil {
		return err
	}
	g.udp_conn = conn
	defer conn.Close()
	buffer := make([]byte, p2p.PacketSize)
	for {
		n, addr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			return err
		}
		fmt.Printf("just read %d bytes from %s: %v\n", n, addr.String(), buffer[:n])
		go g.HandleUDPMessage(conn, buffer[:n])
	}
}

func (g *Gaijo) HandleUDPMessage(conn *net.UDPConn, packet []byte) error {
	message := p2p.Message{}
	if err := message.Deserialize(packet); err != nil {
		return err
	}
	if message.IsRequest {

	} else {
		hash := message.FileId
		queue, ok := g.Chanels.Get(hash)
		if !ok {
			return nil
		}
		queue <- message
	}
	return nil
}
