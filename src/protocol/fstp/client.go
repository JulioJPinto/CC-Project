package fstp

import (
	"fmt"
	"net"

	"github.com/fatih/color"
)

type FSTPclient struct {
	Conn     net.Conn
	UDP_PORT int
}

const DefaultUDPPort = 9090

func NewClient(config Config) (*FSTPclient, error) {
	conn, err := net.Dial("tcp", config.ServerAdress())
	if err != nil {
		return nil, err
	}
	return &FSTPclient{Conn: conn, UDP_PORT: DefaultUDPPort}, nil
}

func (client *FSTPclient) Close() {
	client.Conn.Close()
}

func (client *FSTPclient) Request(request Request) (*Response, error) {
	req_msg := Message(request)
	s, _ := req_msg.Serialize()

	s_data := fmt.Sprint("\nsending: ", s)
	s_str := fmt.Sprint("\nAKA \n\t<", HeaderType(int(s[0])), ">\n\tPayload: ", string(s[HeaderSize:]))
	color.Green(s_data)
	color.Blue(s_str)

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
	data := fmt.Sprint("\nrecieved: ", recieved_data)
	str := fmt.Sprint("\nAKA: \n\t<", HeaderType(int(recieved_data[0])), ">\n\tPayload: ", string(recieved_data[HeaderSize:]))
	color.Green(data)
	color.Blue(str)
	if err != nil {
		client.Conn.Close()
		return nil, err
	}
	resp_msg := &Message{}
	resp_msg.Deserialize(recieved_data)

	resp := Response(*resp_msg)
	return &resp, nil
}

func IHaveFileRequest(props IHaveFileReqProps) Request {
	header := Header{
		IHaveFileReq,
	}
	return Request{header, &props}
}

func AllFilesRequest() Request {
	return Request{Header{AllFilesReq}, nil}
}
