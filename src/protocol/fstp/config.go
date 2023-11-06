package fstp

import "fmt"

// FSTPConfig ...
type FSTPConfig struct {
	Host string
	Port string
}

func (config *FSTPConfig) ServerAdress() string {
	serverAddr := fmt.Sprintf("%s:%s", config.Host, config.Port)
	return serverAddr
}
