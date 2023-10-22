package fstp

import (
	"bufio"
	"cc_project/helpers"
	"net"
)

const (
	IHave  = 0b0001
	WhoHas = 0b0010
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
	s, _ := request.Serialize()
	client.conn.Write(s)

	reader := bufio.NewReader(client.conn)
	message, err := reader.ReadBytes('\r')
	if err != nil {
		client.conn.Close()
		return nil, err
	}
	resp := &FSTPresponse{}
	resp.Deserialize(message)
	return resp, nil
}

func (client *FSTPclient) WhoHasRequest(info File_info) FSTPrequest {

	header := FSTPHeader{
		IHave,
		0, // set as 0 by default
	}

	payload := helpers.SerializableMap(map[string]any{"file": info})
	s, _ := payload.Serialize()
	client.conn.Write(s)
	return FSTPrequest{header, &payload}
}

type IHaveProps struct {
	Files []File_info
}

func IHaveRequest(props IHaveProps) FSTPrequest {

	header := FSTPHeader{
		IHave,
		0, // set as 0 by default
	}
	payload := helpers.SerializableMap(map[string]any{"file": props})
	return FSTPrequest{header, &payload}
}
