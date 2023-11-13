package p2p

import (
	"fmt"
	"net"
)

type Config struct {
	Host string
	Port string
}

func (config *Config) ServerAdress() string {
	serverAddr := fmt.Sprintf("%s:%s", config.Host, config.Port)
	return serverAddr
}

func ParseConfig(serverAddr string) (*Config, error) {
	var parsedConfig Config

	// Assuming the serverAddr has the format "host:port"
	host, port, err := net.SplitHostPort(serverAddr)
	if err != nil {
		return nil, err
	}

	parsedConfig.Host = host
	parsedConfig.Port = port

	return &parsedConfig, nil
}
