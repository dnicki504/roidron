package main

import (
	"fmt"
	"net"
	// "net/netip"
)

func main() {
	go func() {
		conn, err := net.ListenPacket("udp", ":7099")
		if err != nil {
			fmt.Println("err ", err)
		}
		buf := make([]byte, 20)
		for true {
			n, dst, _ := conn.ReadFrom(buf)
			fmt.Println("serv recv", string(buf[:n]))
			fmt.Printf("serv dst %+v \n", dst)
			// time.Sleep(10 * time.Second)
			// conn.WriteTo(buf, dst)
		}
	}()

	for true {

	}
}
