package server

import (
	"BlueIce/internal/protocol"
	"bytes"
	"log"
)

func HandleLoginStart(client *Client, payload []byte) {
	var packet protocol.PacketLoginInStart

	if _, err := packet.ReadFrom(bytes.NewReader(payload)); err != nil {
		log.Println(err)
		return
	}

	log.Printf("Player %s is trying to log in", string(packet.Name.Content))

	options := protocol.PrefixedArray[protocol.GameProfileOption]{
		Length: 0,
	}

	var responsePacket protocol.PacketLoginOutSuccess
	responsePacket.Profile = protocol.GameProfile{
		Name:    packet.Name,
		UUID:    packet.UUID,
		Options: options,
	}

	if err := client.SendPacket(&responsePacket); err != nil {
		log.Println("Error while sending login response", err)
	}

	// Set client to configuration state
	client.State = 3

	// Start configuration
	StartConfiguration(client)
}
