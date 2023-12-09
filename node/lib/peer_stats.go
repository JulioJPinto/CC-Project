package lib

// import "cc_project/helpers"

type PeerStats struct{
	RTT int32
	Load uint8 // percentagem
	NPackets uint32
	NDroppedPackets uint32
}