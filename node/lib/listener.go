package lib

import (
	"cc_project/protocol"
	"cc_project/protocol/p2p"
	"fmt"
	"net"
	"os"
)

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
		go node.HandleP2PRequest(addr, message)
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

func (node *Node) HandleP2PRequest(addr *net.UDPAddr, msg p2p.Message) {
	f_path, ok := node.MyFiles[msg.Header.FileId]
	if !ok {
		return
	}
	segment, err := getSegment(f_path, msg.Header.SegmentOffset)
	if err != nil {
		return
	}
	node.sender.Send(*addr, segment)
}

func getSegment(f_path string, segmentOffset uint32) ([]byte, error) {
	file, err := os.Open(f_path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Seek to the starting position (byte X)
	startingByte := 10 // Replace with the actual starting byte
	_, err = file.Seek(int64(startingByte), 0)
	if err != nil {
		return nil, err
	}

	// Read the desired number of bytes (from byte X to byte Y)
	buffer := make([]byte, protocol.SegmentLength)
	bytesRead, err := file.Read(buffer)
	if err != nil {
		return nil, err
	}
	return buffer[:bytesRead], nil
	// Print the read bytes
}
