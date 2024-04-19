package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
)

func main() {
	if len(os.Args) != 4 {
		fmt.Println("<server_host> <server_port> <filename>")
		os.Exit(0)
	}

	server_host, server_port, filename := os.Args[1], os.Args[2], os.Args[3]

	servAddr := fmt.Sprintf("%s:%s", server_host, server_port)

	tcpAddr, err := net.ResolveTCPAddr("tcp", servAddr)
	if err != nil {
		println("ResolveTCPAddr failed:", err.Error())
		os.Exit(1)
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)

	if err != nil {
		println("Dial failed:", err.Error())
		os.Exit(1)
	}

	defer conn.Close()

	urlAddr := url.URL{Path: fmt.Sprintf("/%s", filename)}
	urlAddr.Host = server_host

	request := http.Request{
		Proto:      "HTTP/1.1",
		Host:       server_host,
		Method:     http.MethodGet,
		RequestURI: urlAddr.RequestURI(),
		Header:     make(http.Header),
		URL:        &urlAddr,

		Body: nil,
	}

	err = request.Write(conn)
	if err != nil {
		println("Write to server failed:", err.Error())
		os.Exit(1)
	}

	responce, err := http.ReadResponse(bufio.NewReader(conn), &request)
	defer func() { responce.Body.Close() }()

	if err != nil {
		println("Write to server failed:", err.Error())
		os.Exit(1)
	}

	buff, err := io.ReadAll(responce.Body)

	if err != nil {
		println("Read responce failed", err.Error())
		os.Exit(1)
	}

	println(string(buff))
}
