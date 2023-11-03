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

