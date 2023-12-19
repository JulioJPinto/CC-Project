package lib

import (
	"fmt"
	"math"
)

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

func (stats *PeerStats) Weight2() float64 {

	ipRTTWeight := 0.25
	p2pRTTWeight := 0.25
	packetsRatioWeight := 0.25

	return ((float64(stats.IP_RTT) * ipRTTWeight) +
		(float64(stats.P2P_RTT) * p2pRTTWeight) +
		((float64(stats.NDroppedPackets) / (float64(stats.NPackets))) * packetsRatioWeight))

}

func (stats *PeerStats) Weight() float64 {

	loss_rate := (float64(stats.NDroppedPackets) / (float64(stats.NPackets)))
	lossweight := math.Pow(10, loss_rate)
	return (float64(stats.P2P_RTT) * lossweight)

}
