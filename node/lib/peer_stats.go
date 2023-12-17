package lib

import "fmt"

// import "cc_project/helpers"

type PeerStats struct {
	IP_RTT          uint32
	P2P_RTT         uint32
	Load            uint8 // percentagem
	NPackets        uint32
	NDroppedPackets uint32
}

func (p PeerStats) String() string {
	return fmt.Sprintf("IP_RTT: %d, P2P_RTT: %d, Load: %d%%, NPackets: %d, NDroppedPackets: %d", p.IP_RTT, p.P2P_RTT, p.Load, p.NPackets, p.NDroppedPackets)
}
