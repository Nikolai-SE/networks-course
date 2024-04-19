package main

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"strings"
	"time"
)

func serve(port string) {
	local_addr, err := net.ResolveUDPAddr("udp", port)
	if err != nil {
		log.Fatalln("Error: ", err)
	}

	conn, err := net.ListenUDP("udp", local_addr)
	if err != nil {
		log.Fatalln("Error: ", err)
	}

	buf := make([]byte, 65536)

	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	defer conn.Close()
	for {
		n, addr, err := conn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println("Error: ", err)
			continue
		}
		recv := string(buf[0:n])
		fmt.Printf("Received: '%s' from %v\n", recv, addr)

		if rnd.Float32() < 0.8 {
			upper := strings.ToUpper(recv)
			conn.WriteToUDP([]byte(upper), addr)
		}
	}
}

func main() {
	serve(":9999")
}
