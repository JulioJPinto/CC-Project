package fstp

import (
	// "cc_project/helpers"
	// "reflect"

	"fmt"
	"reflect"
	"testing"
)

func TestIHavePropsSerialization(t *testing.T) {
	x := IHaveProps{Files: []FileInfo{{Id: 1}}}
	x1 := x
	bytes, err := x.Serialize()
	fmt.Println(bytes)
	if err != nil {
		t.Errorf("Deserialization error: %v", err)
		return
	}
	x.Deserialize(bytes)
	if reflect.DeepEqual(x, x1) {
		fmt.Println("x: ", x)
		fmt.Println("x1: ", x1)
		t.Errorf("x!=x1")
		return
	}
}

func TestIHaveMessagePropsSerialization(t *testing.T) {
	x := IHaveProps{Files: []FileInfo{{Id: 1}}}
	m := FSTPmessage{Header: FSTPHeader{1}, Payload: &x}
	bytes, err := m.Serialize()
	fmt.Println(bytes)
	fmt.Println(bytes[FSTPHEaderSize:])
	if err != nil {
		t.Errorf("Deserialization error: %v", err)
		return
	}
}

// 	err := json.Unmarshal([]byte("{\"file\":{\"files\":[{\"id\":1}]}}"), &x)

// 	if err != nil {
// 		t.Errorf("Deserialization error: %v", err)
// 		return
// 	}
// }

// 	// Create an instance of your FSTPRequest struct with test data
// 	originalMessage := &FSTPrequest{
// 		header:  FSTPHeader{flags: IHave},     // Replace with appropriate test data
// 		payload: &helpers.SerializableMap{}, // Replace with appropriate test data
// 	}

// 	// Serialize the original message
// 	serializedData, err := originalMessage.Serialize()
// 	if err != nil {
// 		t.Errorf("Serialization error: %v", err)
// 		return
// 	}

// 	// Deserialize the serialized data
// 	deserializedMessage := &FSTPrequest{}
// 	err = deserializedMessage.Deserialize(serializedData)
// 	if err != nil {
// 		t.Errorf("Deserialization error: %v", err)
// 		return
// 	}

// 	// Compare the original and deserialized messages to check if they are equal
// 	if !reflect.DeepEqual(originalMessage, deserializedMessage) {
// 		t.Errorf("Deserialized message is not equal to the original message")
// 	}
// }
