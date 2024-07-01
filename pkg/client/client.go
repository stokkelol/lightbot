package client

import (
	"log"
	"net"
)

// GetIP preferred outbound ip of this machine
func GetIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	return conn.RemoteAddr().(*net.UDPAddr).String()
}
