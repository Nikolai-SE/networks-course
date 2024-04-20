package server

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"net"
	"rdt/internal/utils"
	"time"
)

type Server struct {
	Port        int
	Reliability float64
	logger      *log.Logger
	rnd         *rand.Rand
}

func NewServer(port int, logger *log.Logger, reliability float64) Server {
	return Server{
		Port:        port,
		Reliability: reliability,
		rnd:         rand.New(rand.NewSource(time.Now().UnixNano())),
		logger:      logger,
	}
}

func (s Server) Serve(w *bufio.Writer, r *bufio.Reader) error {
	localAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%d", s.Port))
	if err != nil {
		log.Println("Error: ", err)
		return err
	}

	conn, err := net.ListenUDP("udp", localAddr)
	if err != nil {
		log.Println("Error: ", err)
		return err
	}

	buf := make([]byte, 65536)

	defer conn.Close()

	clientSeqNo := make(map[string]int)
	responce := utils.Packet{}

	lengthToSend := 0
	sendedLen := 0
	buffer := make([]byte, 1024*64)

	nonBlockRead := make(chan []byte)
	go func(ch chan []byte) {
		if r == nil {
			close(ch)
			return
		}
		defer close(ch)

		localBuffer := make([]byte, 1024*64)
		for {
			n, err := r.Read(localBuffer)
			if err != nil {
				return
			}
			ch <- localBuffer[0:n]
		}
	}(nonBlockRead)

	for {
		n, addr, err := conn.ReadFromUDP(buf)
		if err != nil {
			s.logger.Fatalln(err)
			continue
		}

		saddr := addr.String()

		packet, perr := utils.ReadPacket(buf[0:n])
		if perr != nil {
			s.logger.Fatalln(perr)
		}

		if _, ok := clientSeqNo[saddr]; !ok {
			clientSeqNo[saddr] = packet.SeqNo
		}

		s.logger.Printf("Received: '%v' from %v\n", packet, addr)

		if packet.SeqNo == clientSeqNo[saddr] {
			w.Write(packet.Payload)
			responce.AckNo = packet.SeqNo

			/* Dual transmission. */
			if sendedLen >= lengthToSend {
				// Read from r
				select {
				case in, ok := <-nonBlockRead:
					if ok {
						buffer = in
						lengthToSend = len(in)
						sendedLen = 0
					}
				default:
				}
			}
			if sendedLen < lengthToSend {
				responce.Payload = buffer[sendedLen : sendedLen+min(lengthToSend-sendedLen, utils.MaxPayloadLen)]
				sendedLen += utils.MaxPayloadLen
			} else {
				responce.Payload = responce.Payload[0:0]
			}
			s.logger.Println("Send Ack#", responce.AckNo)
			w.Flush()
			clientSeqNo[saddr]++
		}
		err = s.send(conn, addr, responce)
		if err != nil {
			s.logger.Fatalln(err)
		}
	}
}

func (s Server) send(conn *net.UDPConn, addr *net.UDPAddr, packet utils.Packet) error {
	data, err := utils.WritePacket(packet)
	if err != nil {
		s.logger.Fatalln(err)
	}

	if s.rnd.Float64() <= s.Reliability {
		_, err := conn.WriteToUDP(data, addr)
		return err
	} else {
		log.Println("Server has lost the packet :(")
	}
	return nil
}
