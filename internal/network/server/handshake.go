package server

import (
	"bytes"
	"log"

	"github.com/blueice-mc/blueice/internal/network/protocol"
)

func HandleHandshake(client *Client, payload []byte) {
	var packet protocol.PacketHandshakeIn

	if _, err := packet.ReadFrom(bytes.NewReader(payload)); err != nil {
		log.Println(err)
		return
	}

	client.State = int32(packet.Intent)
}
