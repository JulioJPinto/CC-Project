package lib

import (
	"cc_project/helpers"
	"cc_project/protocol"
	"time"

	"github.com/go-ping/ping"
)

// import "cc_project/helpers"

type Stats struct {
	RTT     time.Duration
	Latency int32
	Load    uint8 // percentagem
}

type StatsManager struct {
	helpers.SyncMap[protocol.DeviceIdentifier, Stats]
}

// usar numa routine
func (m *StatsManager) UpdateRTT(peer protocol.DeviceIdentifier) error {
	pinger, err := ping.NewPinger(string(peer))
	if err != nil {
		return err
	}
	pinger.Count = 3
	err = pinger.Run() // Blocks until finished.
	if err != nil {
		return err
	}
	ping_stats := pinger.Statistics()

	peer_stats, _ := m.Get(peer)
	peer_stats.RTT = ping_stats.MaxRtt
	m.Set(peer, peer_stats)
	return nil
}
