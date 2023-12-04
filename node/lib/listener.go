package lib

import (
	"cc_project/node/p2p"
	"fmt"
	"net"
)

type Handler struct{}

func (node *Node) ListenOnUDP() error {
	serverAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%s", node.P2PConfig.Host, node.P2PConfig.Port))
	if err != nil {
		return err
	}
	conn, err := net.ListenUDP("udp", serverAddr)
	if err != nil {
		return err
	}
	buffer := make([]byte, 1024) // Adjust the buffer size based on your needs

	for {
		n, addr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println("Error reading from UDP:", err)
			continue
		}

		data := buffer[:n]
		fmt.Printf("Received %d bytes from %s: %s\n", n, addr.String(), string(data))
		node.handleUDPMessage(addr, data)
	}
}

func (node *Node) handleUDPMessage(addr *net.UDPAddr, packet []byte) error {
	message := p2p.Message{}
	if err := message.Deserialize(packet); err != nil {
		return err
	}
	if message.IsRequest {
		fmt.Println(message.Payload)
		// g.HandleP2PRequest()
	} else {
		hash := message.FileId
		queue, ok := node.Chanels.Get(hash)
		if !ok {
			return nil
		}
		queue <- message
	}
	return nil
}
