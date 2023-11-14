package lib

import (
	"cc_project/protocol/p2p"
	"fmt"
	"net"
)

type Handler struct{}

func (*Handler) HandleRequest(p2p.Request) p2p.Response {
	return p2p.Response(p2p.Message{})
}

func ListenOnUDP(config p2p.Config) error {
	serverAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%s", config.Host, config.Port))
	if err != nil {
		return err
	}
	conn, err := net.ListenUDP("udp", serverAddr)
	if err != nil {
		return err
	}
	defer conn.Close()
	buffer := make([]byte, p2p.PacketSize)
	for {
		n, addr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			return err
		}
		fmt.Printf("just read %d bytes from %s: %v\n", n, addr.String(), buffer[:n])
	}
}
