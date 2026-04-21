package server

import (
	"bytes"
	"log"

	"github.com/blueice-mc/blueice/internal/network/protocol"
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

	SendEmptyChunks(client)
}

func NewEmptyChunk(chunkX, chunkZ int32) *protocol.PacketPlayOutLevelChunkWithLight {
	var emptySkyLight protocol.BitSet
	emptySkyLight.SetRange(0, 25)

	var emptyBlockLight protocol.BitSet
	emptyBlockLight.SetRange(0, 25)

	worldSurface := &protocol.Heightmap{Type: 1, WorldHeight: 384}
	motionBlocking := &protocol.Heightmap{Type: 4, WorldHeight: 384}

	var sectionBuf bytes.Buffer

	for i := 0; i < 24; i++ {
		// nonEmptyBlockCount
		protocol.WriteInt16(&sectionBuf, 0)
		// fluidCount
		protocol.WriteInt16(&sectionBuf, 0)
		// Block PalettedContainer
		protocol.WriteUint8(&sectionBuf, 0) // bits per entry = 0
		varInt := protocol.VarInt(0)
		varInt.WriteTo(&sectionBuf) // Luft ID
		varInt.WriteTo(&sectionBuf) // 0 storage longs
		// Biome PalettedContainer
		protocol.WriteUint8(&sectionBuf, 0) // bits per entry = 0
		varInt.WriteTo(&sectionBuf)         // plains ID
		varInt.WriteTo(&sectionBuf)         // 0 storage longs
	}

	return &protocol.PacketPlayOutLevelChunkWithLight{
		ChunkX: chunkX,
		ChunkZ: chunkZ,
		Heightmaps: protocol.PrefixedArray[protocol.Heightmap]{
			Content: []protocol.Heightmap{*worldSurface, *motionBlocking},
		},
		Data:          protocol.PrefixedArray[uint8]{Content: sectionBuf.Bytes()},
		BlockEntities: protocol.PrefixedArray[protocol.BlockEntity]{Content: []protocol.BlockEntity{}},
		LightData: protocol.LightData{
			SkyLightMask:        protocol.BitSet{},
			BlockLightMask:      protocol.BitSet{},
			EmptySkyLightMask:   emptySkyLight,
			EmptyBlockLightMask: emptyBlockLight,
			SkyLightArray:       protocol.PrefixedArray[protocol.LightArray]{Content: []protocol.LightArray{}},
			BlockLightArray:     protocol.PrefixedArray[protocol.LightArray]{Content: []protocol.LightArray{}},
		},
	}
}

func SendEmptyChunks(client *Client) {
	for x := int32(-2); x <= 2; x++ {
		for z := int32(-2); z <= 2; z++ {
			chunk := NewEmptyChunk(x, z)
			client.SendPacket(chunk)
		}
	}
}
