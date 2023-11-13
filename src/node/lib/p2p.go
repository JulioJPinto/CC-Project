package lib

import (
	"cc_project/protocol/p2p"
	"net"
)

type Handler struct {}

func (*Handler)HandleRequest(net.Conn,p2p.Request) p2p.Response{
	return p2p.Response(p2p.Message{})
}