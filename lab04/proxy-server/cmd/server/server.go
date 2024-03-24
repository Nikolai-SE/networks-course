package main

import (
	"flag"

	"nikolai/proxy-server/internal/proxy"
)

var (
	address   = flag.String("address", "localhost", "proxy address")
	port      = flag.Int("port", 8888, "proxy port")
	blacklist = flag.String("blacklist", "blacklist.json", "blacklist path")
)

func main() {
	proxy.RunProxy(*address, *port, *blacklist)
}
