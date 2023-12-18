package state_manager

import (
	"cc_project/helpers"
	"encoding/json"
	"fmt"
	"log"
)

func (s *StateManager) Files() helpers.StatusMessage {
	ret := helpers.NewStatusMessage()

	// s_n := s.SegmentsNodes()
	files := s.GetAllFiles()
	for _, f := range files {
		x, _ := json.Marshal(f)
		ret.AddMessage(nil, fmt.Sprint(string(x)))
	}
	return ret
}


func Shutdown() helpers.StatusMessage {
	fmt.Println("would be a good time to save state to disk")
	log.Fatal("shuting down ...")
	return helpers.NewStatusMessage()
}
