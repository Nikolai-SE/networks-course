package utils

import (
	"encoding/json"
)

const (
	MaxPayloadLen = 1500
)

type Packet struct {
	SeqNo   int    `json:"seq"`
	AckNo   int    `json:"ack"`
	Payload []byte `json:"payload"`
}

func ReadPacket(data []byte) (packet Packet, err error) {
	err = json.Unmarshal(data, &packet)
	return
}

func WritePacket(packet Packet) (data []byte, err error) {
	data, err = json.Marshal(packet)
	return
}
