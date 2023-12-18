package state_manager

import (
	"cc_project/helpers"
	"encoding/json"
	"fmt"
)

func (s *StateManager) Files() helpers.StatusMessage {
	ret := helpers.NewStatusMessage()

	// s_n := s.SegmentsNodes()
	files := s.GetAllFiles()
	for _, f := range files {
		x, _ := json.Marshal(f)
		ret.AddMessage(nil, fmt.Sprint(string(x)))
	}
	// for segment, nodes := range s_n {
	// 	for _, node := range nodes {
	// 		ret.AddMessage(nil, fmt.Sprint(string(node), " has ", segment.String()))

	// 	}
	// }
	return ret
}
