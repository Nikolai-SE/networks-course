package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"rdt/internal/client"
	"time"
)

var (
	port    = flag.Int("port", 9999, "server port")
	addr    = flag.String("address", "127.0.0.1", "server address")
	timeout = flag.Int("timeout", 1, "timeout in seconds")
)

func main() {
	flag.Parse()

	c := client.NewClient(fmt.Sprintf("%s:%d", *addr, *port), 0.7)
	rw := bufio.NewReadWriter(bufio.NewReader(os.Stdin), bufio.NewWriter(os.Stdout))
	log.Println(c.Process(rw, time.Duration(*timeout)*time.Second))
}