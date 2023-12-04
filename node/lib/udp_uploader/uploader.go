package udp_uploader

import (
	"fmt"
	"net"
)

type Uploader struct {
	queue chan struct {
		Address net.UDPAddr
		Data    []byte
	}
}

func NewUploader(n int) *Uploader {
	uploader := &Uploader{}
	uploader.queue = make(chan struct {
		Address net.UDPAddr
		Data    []byte
	})

	for i := 0; i < n; i++ {
		go uploader.sender()
	}

	return uploader
}

// func (u *Uploader) Send(address string, data []byte) {
// 	u.queue <- struct {
// 		Address net.UDPAddr
// 		Data    []byte
// 	}{address, data}

// }

func (u *Uploader) Send(address net.UDPAddr, data []byte) {

}

func (u *Uploader) sender() {
	for {
		message := <-u.queue
		// UDP sending logic
		// destination, err := net.ResolveUDPAddr("udp", message.Address)
		// if err != nil {
		// 	fmt.Println("Error resolving server address:", err)
		// 	continue
		// }

		conn, err := net.DialUDP("udp", nil, &message.Address)
		if err != nil {
			fmt.Println("Error creating UDP connection:", err)
			continue
		}

		_, err = conn.Write(message.Data)
		if err != nil {
			fmt.Println("Error sending UDP packet:", err)
			continue
		}

		// fmt.Printf("Sent data to %s: %s\n", message.Address, string(message.Data))

	}
}
