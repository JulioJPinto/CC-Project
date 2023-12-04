package p2p

import (
	"fmt"
	// "net"
)

type Config struct {
	Host string
	Port string
}

func (config *Config) ServerAdress() string {
	serverAddr := fmt.Sprintf("%s:%s", config.Host, config.Port)
	return serverAddr
}
