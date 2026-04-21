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
}
