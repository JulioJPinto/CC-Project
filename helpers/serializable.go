package helpers

import "encoding/json"

type Serializable interface {
	Serialize() ([]byte, error)
	Deserialize([]byte) error
}

type SerializableMap map[string]interface{}

func (sm SerializableMap) Serialize() ([]byte, error) {
	return json.Marshal(sm)
}

func (sm *SerializableMap) Deserialize(data []byte) error {
	return json.Unmarshal(data, sm)
}
