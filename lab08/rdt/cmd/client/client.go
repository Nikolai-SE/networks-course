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
	port        = flag.Int("port", 9999, "server port")
	addr        = flag.String("address", "127.0.0.1", "server address")
	timeout     = flag.Int("timeout", 1, "timeout in seconds")
	reliability = flag.Float64("reliability", 0.7, "reliability of client")
)

func main() {
	flag.Parse()

	c := client.NewClient(fmt.Sprintf("%s:%d", *addr, *port), *reliability)
	log.Println(c.Process(bufio.NewReader(os.Stdin), bufio.NewWriter(os.Stdout), time.Duration(*timeout)*time.Second))
}
