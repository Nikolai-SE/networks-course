package main

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
)

func main() {
	server_addr, err := net.ResolveUDPAddr("udp", "127.0.0.1:9999")
	if err != nil {
		log.Fatalln("Error: ", err)
	}

	local_addr, err := net.ResolveUDPAddr("udp", "127.0.0.1:0")
	if err != nil {
		log.Fatalln("Error: ", err)
	}

	conn, err := net.DialUDP("udp", local_addr, server_addr)
	if err != nil {
		log.Fatalln("Error: ", err)
	}

	var wg sync.WaitGroup

	waiter := func() {
		defer wg.Done()
		rxbuf := make([]byte, 65536)
		conn.SetDeadline(time.Now().Add(time.Second))
		n, _, err := conn.ReadFromUDP(rxbuf)

		if err != nil {
			log.Println("Error: ", err)
			return
		}

		endT := time.Now()

		recv := string(rxbuf[0:n])

		us, err := strconv.ParseInt(strings.Split(recv, " ")[2], 10, 64)
		if err != nil {
			log.Println("Error: ", err)
		} else {
			fmt.Printf("Received '%s'  (RTT: %d us)\n", recv, endT.UnixMicro()-us)
		}
	}

	defer conn.Close()
	for i := range 10 {
		wg.Add(1)

		beginT := time.Now()
		msg := fmt.Sprintf("Ping %d %d", 1+i, beginT.UnixMicro())
		txbuf := []byte(msg)
		_, err := conn.Write(txbuf)
		if err != nil {
			fmt.Println(msg, err)
		}
		go waiter()
	}

	wg.Wait()
}
