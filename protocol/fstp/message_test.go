package fstp

// import (
// 	"cc_project/helpers"
// 	"reflect"
// 	"testing"
// )

// func TestSerializationDeserialization(t *testing.T) {
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
