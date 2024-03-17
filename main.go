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
			fmt.Println("serv recv s", string(buf[:n]))
			fmt.Println("serv recv b", buf[:n])
			fmt.Printf("serv dst %+v \n", dst)
			// time.Sleep(10 * time.Second)
			dst, _ = net.ResolveUDPAddr("udp", "192.168.0.51:7088")
			conn.WriteTo(buf, dst)
		}
	}()
	fmt.Println("started")
	for true {

	}
}
