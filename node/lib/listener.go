package lib

import (
	"cc_project/helpers"
	"cc_project/protocol"
	"cc_project/protocol/p2p"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"time"

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
	buffer := make([]byte, 2048) // Adjust the buffer size based on your needs

	for {
		n, addr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println("Error reading from UDP:", err)
			continue
		}

		data := buffer[:n]
		// s := fmt.Sprintf("\n\nReceived %d bytes from %s: %s\n\n", n, addr.String(), string(data))
		s := fmt.Sprintf("\n\nReceived %d bytes from %s\n\n", n, addr.String())
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
		s, _ := json.Marshal(message.Header)
		color.Green(string(s))
		go node.HandleP2PRequest(addr, message)
	} else {
		now := helpers.TrunkI64(time.Now().UnixMilli())
		delay := now - int32(message.TimeStamp)
		s := fmt.Sprintf("responded in %d miliseconds", now-int32(message.TimeStamp))
		color.Green(s)

		hash := message.FileId
		downloader, ok := node.Downloads.Load(hash)
		if !ok {
			color.Red("channel dont exist")
			// fmt.Println(node.Chanels.())
			return nil
		}
		peer := protocol.DeviceIdentifier(addr.IP.String())
		stats, ok := node.PeerStats.Load(peer)
		if !ok {
			stats = &PeerStats{}
		}
		stats.P2P_RTT = uint32(delay)
		stats.NPackets++
		node.PeerStats.Store(peer, stats)
		downloader.ForwardMessage(message, peer)
	}
	return nil
}

func (node *Node) HandleP2PRequest(addr *net.UDPAddr, msg p2p.Message) {

	f_path, ok := node.MyFiles[msg.FileId]
	if !ok {
		return
	}
	segment, err := getSegment(f_path, msg.Header.SegmentOffset)
	if err != nil {
		color.Red("ERROR GETTING SEGMENT")
	} else {
		color.Yellow("segment:")
		color.Yellow(string(segment))
	}
	segment_data, _ := node.KnownFiles.Load(msg.FileId) // WARN !!!! NOT CHECKING THIS
	hash := segment_data.SegmentHashes[msg.Header.SegmentOffset]
	f := protocol.FileSegment{
		BlockOffset: int64(msg.SegmentOffset),
		FileHash:    msg.FileId,
		Hash:        hash}
	if err != nil {
		return
	}
	addr.Port = 9090
	ret_msg := p2p.GivYouFileSegmentResponse(f, segment, msg.TimeStamp)
	m := p2p.Message(ret_msg)
	bytes, _ := m.Serialize()
	color.Cyan(string(ret_msg.Payload))
	node.sender.Send(*addr, bytes)
}

func getSegment(f_path string, segmentOffset uint32) ([]byte, error) {
	file, err := os.Open(f_path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	startingByte := segmentOffset * protocol.SegmentMaxLength // Replace with the actual starting byte
	_, err = file.Seek(int64(startingByte), 0)
	if err != nil {
		return nil, err
	}

	buffer := make([]byte, protocol.SegmentMaxLength)
	n, err := file.Read(buffer)
	if err != nil {
		return nil, err
	}
	return buffer[:n], nil
	// Print the read bytes
}
