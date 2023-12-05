package lib

import (
	"cc_project/protocol"
	"cc_project/protocol/p2p"
	"fmt"
	"net"
	"os"

	"github.com/fatih/color"
)

func (node *Node) ListenOnUDP() error {
	// fmt.Sprintf("%s:%s", node.P2PConfig.Host, node.P2PConfig.Port)
	serverAddr, err := net.ResolveUDPAddr("udp", "0.0.0.0:9090")
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
		s := fmt.Sprintf("\n\nReceived %d bytes from %s: %s\n\n", n, addr.String(), string(data))
		color.Green(s)
		node.handleUDPMessage(addr, data)
	}
}

func (node *Node) handleUDPMessage(addr *net.UDPAddr, packet []byte) error {
	message := p2p.Message{}
	if err := message.Deserialize(packet); err != nil {
		return err
	}
	if message.IsRequest {

		color.Green("itssa requestttt")
		go node.HandleP2PRequest(addr, message)
	} else {
		color.Green("itssa responsss")

		hash := message.FileId
		queue, ok := node.Chanels.Get(hash)
		if !ok {
			color.Red("channel dont exist")
			fmt.Println(node.Chanels.Keys())
			return nil
		}
		queue <- message
	}
	return nil
}

func (node *Node) HandleP2PRequest(addr *net.UDPAddr, msg p2p.Message) {
	f_path, ok := node.MyFiles[msg.FileId]
	if !ok {
		return
	}
	segment, err := getSegment(f_path, msg.Header.SegmentOffset)
	segment_data := node.KnownFiles[msg.FileId]
	hash := segment_data.SegmentHashes[msg.Header.SegmentOffset]
	f := protocol.FileSegment{
		BlockOffset: int64(msg.SegmentOffset),
		FileHash:    msg.FileId,
		Hash:        hash}
	if err != nil {
		return
	}
	addr.Port = 9090
	ret_msg := p2p.GivYouFileSegmentResponse(f, segment, 0)
	bytes, _ := ret_msg.Serialize()
	node.sender.Send(*addr, bytes)
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
