package lib

import (
	helpers_sync "cc_project/helpers/sync"
	"cc_project/protocol"
	"math/rand"
	"math"
	"time"
)

func GenerateRandomFloat() float64 {
	rand.Seed(time.Now().UnixNano())
	return rand.Float64()
}

func FindClosestNode(peerStats map[protocol.DeviceIdentifier]float64, targetValue float64) protocol.DeviceIdentifier {
	var closestNode protocol.DeviceIdentifier
	closestDistance := -1.0

	for node, weight := range peerStats {
		distance := math.Abs(weight - targetValue)

		if closestDistance == -1.0 || distance < closestDistance {
			closestNode = node
			closestDistance = distance
		}
	}

	return closestNode
}

func SelectNode(peerStats *helpers_sync.Map[protocol.DeviceIdentifier, *PeerStats]) protocol.DeviceIdentifier {
	converted := make(map[protocol.DeviceIdentifier]float64)
	var totalWeight float64

	peerStats.Range(func(deviceID protocol.DeviceIdentifier, stats *PeerStats) bool {
		weight := PeerWeight(*stats)
		converted[deviceID] = weight
		totalWeight += weight
		return true
	})

	randomWeight := GenerateRandomFloat() * totalWeight
	selectedNode := FindClosestNode(converted, randomWeight)

	return selectedNode
}
