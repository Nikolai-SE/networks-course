package client

import (
	"bufio"
	"log"
	"math/rand"
	"net"
	"rdt/internal/utils"
	"time"
)

type Client struct {
	addr        string
	conn        *net.UDPConn
	Reliability float64
	serverAddr  *net.UDPAddr
}

func NewClient(addr string, reliability float64) Client {
	serverAddr, err := net.ResolveUDPAddr("udp", addr) // "127.0.0.1:9999"
	if err != nil {
		log.Fatalln(err)
	}

	conn, err := net.DialUDP("udp", nil, serverAddr)
	if err != nil {
		log.Fatalln(err)
	}

	return Client{addr: addr, conn: conn, Reliability: reliability, serverAddr: serverAddr}
}

func (c Client) write(data []byte) {
	if rand.Float64() <= c.Reliability {
		_, err := c.conn.Write(data)
		if err != nil {
			log.Println(err)
		}
	} else {
		log.Println("Opps...")
	}
}

func (c Client) Process(r *bufio.Reader, timeout time.Duration) error {
	rxbuf := make([]byte, 65536)
	seqNo := 0
	waiter := func(conn *net.UDPConn, ackNo int) (bool, error) {
		for {
			conn.SetDeadline(time.Now().Add(timeout))
			n, err := conn.Read(rxbuf)
			conn.SetDeadline(time.Time{})

			if err != nil {
				return false, err
			}

			packet, err := utils.ReadPacket(rxbuf[0:n])
			if err != nil {
				log.Println(err)
				continue
			}

			if packet.AckNo == ackNo {
				log.Println("Ack #", packet.AckNo)
				return true, nil
			}
		}
	}

	buffer := make([]byte, 1024*64)
	n := 0
	sendedLen := 0

	packet := utils.Packet{}
	for {
		packet.SeqNo = seqNo
		if sendedLen >= n {
			readedLen, err := r.Read(buffer)
			if err != nil {
				log.Fatalln(err)
			}
			n = readedLen
			sendedLen = 0
		}
		packet.Payload = buffer[sendedLen : sendedLen+min(n-sendedLen, utils.MaxPayloadLen)]
		sendedLen += utils.MaxPayloadLen

		data, err := utils.WritePacket(packet)
		if err != nil {
			log.Fatalln(err)
		}

		for {
			c.write(data)
			recived, err := waiter(c.conn, seqNo)
			if err != nil {
				log.Println(err)
			}
			if recived {
				seqNo++
				break
			}
		}
	}
}
