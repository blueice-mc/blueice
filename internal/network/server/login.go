package server

import (
	"BlueIce/internal/game/entity"
	"BlueIce/internal/network/protocol"
	"bytes"
	"log"
)

func HandleLoginStart(client *Client, payload []byte) {
	var packet protocol.PacketLoginInStart

	if _, err := packet.ReadFrom(bytes.NewReader(payload)); err != nil {
		log.Println(err)
		return
	}

	log.Printf("Player %s is trying to log in", packet.Name)

	options := protocol.PrefixedArray[protocol.GameProfileOption]{
		Length: 0,
	}

	client.PendingProfile = &protocol.GameProfile{
		Name:    packet.Name,
		UUID:    packet.UUID,
		Options: options,
	}

	var responsePacket protocol.PacketLoginOutSuccess
	responsePacket.Profile = *client.PendingProfile

	if err := client.SendPacket(&responsePacket); err != nil {
		log.Println("Error while sending login response", err)
	}
}

func HandleLoginAcknowledged(client *Client, payload []byte) {
	client.Player = &entity.Player{
		UUID:       client.PendingProfile.UUID,
		PlayerName: client.PendingProfile.Name,
		Connection: client,
	}

	client.Server.mu.Lock()
	client.Server.Players = append(client.Server.Players, client.Player)
	client.Server.mu.Unlock()

	// Switch to configuration state
	client.State = 3

	// Start the configuration for the client
	StartConfiguration(client)
}
