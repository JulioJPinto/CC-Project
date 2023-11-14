package p2p


import (
	// "fmt"
	// "net"
)

type ResponseHandler interface {
	HandleResponse(Response)
}

type Client struct {
	host    string
	port    string
	handler Handler
}


type p2p_client_routine struct {
	handler Handler
}

func NewP2PClient(config Config, handler Handler) *Client {
	return &Client{host: config.Host, port: config.Port, handler: handler}

}