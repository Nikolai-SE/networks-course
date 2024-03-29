package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"os/exec"
	"strings"
)

func handle(tcp *net.TCPConn) (err error) {
	rw := bufio.ReadWriter{Reader: bufio.NewReader(tcp), Writer: bufio.NewWriter(tcp)}
	defer tcp.Close()
	defer rw.Writer.Flush()

	for {
		line, err := rw.Reader.ReadString('\n')
		if err != nil {
			return err
		}

		spl := strings.Split(strings.TrimRight(line, "\r\n"), " ")

		switch spl[0] {
		case "q", "quit", "exit":
			return nil
		default:
		}

		cmd := exec.Command(spl[0], spl[1:]...)
		out, err := cmd.CombinedOutput()
		if err != nil {
			return err
		}

		rw.Writer.Write(out)
		rw.Writer.Flush()
	}
}

var (
	addr = flag.String("address", "localhost", "server address")
	port = flag.Int("port", 1350, "port for remote procedure call")
)

func main() {
	flag.Parse()

	var address net.TCPAddr
	address.IP = net.ParseIP(*addr)
	address.Port = *port

	listener, err := net.ListenTCP("tcp", &address)

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	defer listener.Close()

	for {
		tcp, err := listener.AcceptTCP()
		if err == nil {
			go handle(tcp)
		}
	}
}
