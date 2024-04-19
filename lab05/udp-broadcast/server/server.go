package main

import (
	"flag"
	"fmt"
	"net"
	"time"
)

var (
	address = flag.String("address", "localhost", "broadcast address")
	port    = flag.Int("port", 8829, "broadcast port")
)

func main() {
	flag.Parse()

	listenAddr, err := net.ResolveUDPAddr("udp4", fmt.Sprintf(":%d", *port))
	if err != nil {
		panic(err)
	}
	list, err := net.ListenUDP("udp4", listenAddr)
	if err != nil {
		panic(err)
	}
	defer list.Close()

	addr, err := net.ResolveUDPAddr("udp4", fmt.Sprintf("%s:%d", *address, *port))
	if err != nil {
		panic(err)
	}

	for {
		t := time.Now()
		_, err := list.WriteTo([]byte(fmt.Sprintf("Time: %s", t.String())), addr)
		time.Sleep(time.Second)
		if err != nil {
			panic(err)
		}
	}
}
