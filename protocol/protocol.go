package protocol

import (
	"net"
)

type Protocol interface {
	CreateConn(net.IP) error
	CloseConn() error
}
