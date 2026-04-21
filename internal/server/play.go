package server

import (
	"BlueIce/internal/protocol"
	"log"
)

func StartPlay(client *Client) {
	loginPacket := &protocol.PacketPlayOutLogin{
		EntityID:   1,
		IsHardcore: false,
		DimensionNames: protocol.PrefixedArray[protocol.Identifier]{
			Content: []protocol.Identifier{
				protocol.NewIdentifierFromPath("overworld"),
			},
		},
		MaxPlayers:          20,
		ViewDistance:        10,
		SimulationDistance:  10,
		ReducedDebugInfo:    false,
		EnableRespawnScreen: true,
		DoLimitedCrafting:   false,
		DimensionType:       0,
		DimensionName:       protocol.NewIdentifierFromPath("overworld"),
		HashedSeed:          0,
		GameMode:            1,
		PreviousGameMode:    -1,
		IsDebug:             false,
		IsFlat:              false,
		HasDeathLocation:    false,
		PortalCooldown:      0,
		SeaLevel:            64,
		EnforcesSecureChat:  false,
	}

	if err := client.SendPacket(loginPacket); err != nil {
		log.Println("Error while sending login", err)
	}

	gameEventPacket := &protocol.PacketPlayOutGameEvent{
		Event: 13,
		Value: 0.0,
	}

	if err := client.SendPacket(gameEventPacket); err != nil {
		log.Println("Error while sending game event", err)
	}

	playerPositionPacket := &protocol.PacketPlayOutPlayerPosition{
		TeleportID: 1,
		X:          0.0,
		Y:          64.0,
		Z:          0.0,
		VelocityX:  0.0,
		VelocityY:  0.0,
		VelocityZ:  0.0,
		Yaw:        0.0,
		Pitch:      0.0,
		Flags:      0,
	}

	if err := client.SendPacket(playerPositionPacket); err != nil {
		log.Println("Error while sending player_position", err)
	}
}
