package fstp

import (
	"cc_project/helpers"
	"fmt"
	"net"
)

type FSTPclient struct {
	conn net.Conn
}

func NewFSTPClient(config FSTPConfig) (*FSTPclient, error) {
	conn, err := net.Dial("tcp", config.ServerAdress())
	if err != nil {
		return nil, err
	}
	return &FSTPclient{conn}, nil
}

func (client *FSTPclient) Close() {
	client.conn.Close()
}

func (client *FSTPclient) Request(request FSTPrequest) (*FSTPresponse, error) {
	req_msg := FSTPmessage(request)
	s, _ := req_msg.Serialize()
	client.conn.Write(s)

	var recieved_data []byte
	buffer := make([]byte, buffer_limit) // Create a buffer to store incoming data
	var err error
	for {
		n, err := client.conn.Read(buffer)
		if err != nil {
			fmt.Println("Error reading:", err)
			break
		}

		recieved_data = append(recieved_data, buffer[:n]...)
		if n != buffer_limit {
			break
		}
	}
	fmt.Println("recieved:", recieved_data)
	if err != nil {
		client.conn.Close()
		return nil, err
	}
	resp_msg := &FSTPmessage{}
	resp_msg.Deserialize(recieved_data)
	resp := FSTPresponse(*resp_msg)
	return &resp, nil
}

func (client *FSTPclient) WhoHasRequest(info FileInfo) FSTPrequest {

	header := FSTPHeader{
		IHave,
	}

	payload := helpers.SerializableMap(map[string]any{"file": info})
	s, _ := payload.Serialize()
	client.conn.Write(s)
	return FSTPrequest{header, &payload}
}

func IHaveRequest(props IHaveProps) FSTPrequest {

	header := FSTPHeader{
		IHave,
	}	
	return FSTPrequest{header, &props}
}
