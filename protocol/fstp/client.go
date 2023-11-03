package fstp

import (
	"fmt"
	"net"
)

type FSTPclient struct {
	Conn net.Conn
}

func NewFSTPClient(config FSTPConfig) (*FSTPclient, error) {
	conn, err := net.Dial("tcp", config.ServerAdress())
	if err != nil {
		return nil, err
	}
	return &FSTPclient{conn}, nil
}

func (client *FSTPclient) Close() {
	client.Conn.Close()
}

func (client *FSTPclient) Request(request FSTPrequest) (*FSTPresponse, error) {
	req_msg := FSTPmessage(request)
	s, _ := req_msg.Serialize()
	client.Conn.Write(s)

	var recieved_data []byte
	buffer := make([]byte, buffer_limit) // Create a buffer to store incoming data
	var err error
	for {
		n, err := client.Conn.Read(buffer)
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
		client.Conn.Close()
		return nil, err
	}
	resp_msg := &FSTPmessage{}
	resp_msg.Deserialize(recieved_data)
	resp := FSTPresponse(*resp_msg)
	return &resp, nil
}

func IHaveRequest(props IHaveProps) FSTPrequest {
	header := FSTPHeader{
		IHaveReq,
	}
	return FSTPrequest{header, &props}
}

func IHaveFileRequest(props IHaveFileProps) FSTPrequest {
	header := FSTPHeader{
		IHaveFileReq,
	}
	return FSTPrequest{header, &props}
}
