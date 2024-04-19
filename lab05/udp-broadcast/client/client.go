package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
)

var (
	port = flag.Int("port", 8829, "broadcast port")
)

func main() {
	flag.Parse()

	// Resolve the string address to a UDP address
	udpAddr, err := net.ResolveUDPAddr("udp4", fmt.Sprintf(":%d", *port))

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Dial to the address with UDP
	conn, err := net.DialUDP("udp4", nil, udpAddr)

	if err != nil {
		log.Fatal(err)
	}

	reader := bufio.NewReader(conn)
	for {
		// Read from the connection untill a new line is send
		data, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(data)
	}
}
