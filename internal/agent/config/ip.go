package config

import (
	"log"
	"net"
)

func GetLocalIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		cErr := conn.Close()
		if cErr != nil {
			log.Fatal("failed to close ip connection")
		}
	}()

	localAddress := conn.LocalAddr().(*net.UDPAddr)

	return localAddress.IP
}
