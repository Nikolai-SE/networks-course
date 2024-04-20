package main

import (
	"bufio"
	"flag"
	"log"
	"os"
	Server "rdt/internal/server"
)

var (
	port        = flag.Int("port", 9999, "server port")
	reliability = flag.Float64("reliability", 1.0, "reliability of server")
)

func main() {
	flag.Parse()
	l := log.New(os.Stderr, "Server log: ", log.Ltime)
	s := Server.NewServer(*port, l, *reliability)
	rw := bufio.NewReadWriter(bufio.NewReader(os.Stdin), bufio.NewWriter(os.Stdout))
	s.Serve(rw)
}
