package state_manager

import (
	"cc_project/helpers"
	"fmt"
)

func (s *StateManager) Files() helpers.StatusMessage {
	ret := helpers.NewStatusMessage()

	s_n := s.SegmentsNodes()
	for segment, nodes := range s_n {
		for _, node := range nodes {
			ret.AddMessage(nil, fmt.Sprint(string(node), " has ", segment.String()))

		}
	}
	return ret
}
