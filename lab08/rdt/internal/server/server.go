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

func (s Server) Serve(rw *bufio.ReadWriter) error {
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
			rw.Write(packet.Payload)
			responce.AckNo = packet.SeqNo

			s.logger.Println("Send Ack#", responce.AckNo)
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
