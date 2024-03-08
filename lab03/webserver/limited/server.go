package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"path"
	"strconv"
)

func handle(tcp *net.TCPConn) (err error) {
	defer tcp.Close()
	request, err := http.ReadRequest(bufio.NewReader(tcp))

	if err != nil {
		return
	}

	defer request.Body.Close()

	responce := http.Response{
		Proto:      request.Proto,
		Request:    request,
		ProtoMajor: request.ProtoMajor,
		ProtoMinor: request.ProtoMinor,
		Header:     make(http.Header),
	}

	file, err := os.ReadFile(path.Join(".", request.RequestURI))
	if err != nil {
		return
	}

	if os.IsNotExist(err) {
		responce.StatusCode = http.StatusNotFound
	} else if err != nil {
		responce.StatusCode = http.StatusBadRequest
	} else {

		responce.Body = io.NopCloser(bytes.NewBuffer(file))
		responce.ContentLength = int64(len(file))
		responce.StatusCode = http.StatusOK
	}

	responce.Status = http.StatusText(responce.StatusCode)
	return responce.Write(tcp)
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("<port> <goroutines-limit>")
		os.Exit(0)
	}

	port, err := strconv.Atoi(os.Args[1])

	if err != nil {
		fmt.Printf("'%s' is not port (error: %s)", os.Args[1], err.Error())
		os.Exit(1)
	}

	goroutinesLimit, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Printf("Parse goroutines limit failed. (error: %s)", err.Error())
		os.Exit(1)
	}

	var address net.TCPAddr
	address.IP = net.ParseIP("localhost")
	address.Port = port

	listener, err := net.ListenTCP("tcp", &address)

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	defer listener.Close()

	guard := make(chan struct{}, goroutinesLimit)

	for {
		tcp, err := listener.AcceptTCP()
		if err == nil {
			guard <- struct{}{} // would block if guard channel is already filled
			go func(tcp *net.TCPConn) {
				handle(tcp)
				<-guard
			}(tcp)
		}
	}
}
