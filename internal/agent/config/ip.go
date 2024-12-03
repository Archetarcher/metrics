package config

import (
	"fmt"
	"log"
	"net"
)

func GetLocalIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddress := conn.LocalAddr().(*net.UDPAddr)

	fmt.Println(localAddress.IP)
	return localAddress.IP
}
