package server

import (
	"BlueIce/internal/protocol"
	"bytes"
	"log"
)

func HandleHandshake(client *Client, payload []byte) {
	var packet protocol.HandshakePacketInbound

	if _, err := packet.ReadFrom(bytes.NewReader(payload)); err != nil {
		log.Println(err)
		return
	}

	log.Println("Received handshake packet")
	client.State = 1
}
