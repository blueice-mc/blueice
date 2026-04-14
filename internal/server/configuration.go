package server

import (
	"BlueIce/internal/defs"
	"BlueIce/internal/protocol"
	"bytes"
	"log"
)

func StartConfiguration(client *Client) {
	brand := "BlueIce"
	var buf bytes.Buffer
	if _, err := protocol.VarInt(len(brand)).WriteTo(&buf); err != nil {
		log.Println(err)
	}
	buf.WriteString(brand)

	brandPacket := protocol.PluginClientMessagePacketOutbound{
		Channel: protocol.Identifier{
			Namespace: "minecraft",
			Path:      "brand",
		},
		Message: buf.Bytes(),
	}

	if err := client.SendPacket(&brandPacket); err != nil {
		log.Println("Error while sending brand", err)
		return
	}

	SendRegistryPackets(client)
}

func SendRegistryPackets(client *Client) {

	// Create overworld packet

	overworldData := &defs.DimensionTypeEntry{
		CoordinateScale:     1.0, // [Double] -> float64
		HasSkylight:         1,   // [Boolean] -> int8
		HasCeiling:          0,   // [Boolean] -> int8
		HasEnderDragonFight: 0,   // [Boolean] -> int8
		AmbientLight:        0.0, // [Float] -> float32

		// Zeit-Logik
		HasFixedTime: 0,   // [Boolean] -> int8 (false)
		FixedTime:    nil, // (optional) -> bleibt weg, da HasFixedTime 0 ist

		MonsterSpawnBlockLightLimit: 0, // [Int] -> int32
		MonsterSpawnLightLevel: defs.IntProvider{
			Type:  "minecraft:constant",
			Value: 0,
		},

		LogicalHeight: 384, // [Int] -> int32
		MinY:          -64, // [Int] -> int32
		Height:        384, // [Int] -> int32

		Infiniburn: "#minecraft:infiniburn_overworld", // [String] mit # laut deiner Quelle

		// Die neuen Felder 1.21.2+
		Skybox:        "minecraft:overworld", // [String]
		CardinalLight: "minecraft:default",   // [String]

		// Attributes muss ein Compound sein (in Go eine Map oder leeres Struct)
		Attributes: defs.EmptyCompound{},

		DefaultClock: "minecraft:overworld", // [String]

		// Timelines ist eine Liste (NBT List)
		Timelines: []string{},
	}

	overworldPacket := &protocol.RegistryDataPacketOutbound{
		RegistryID: protocol.Identifier{Namespace: "minecraft", Path: "dimension_type"},
		Entries: protocol.PrefixedArray[protocol.RegistryEntry]{
			Content: []protocol.RegistryEntry{
				{
					EntryID: protocol.Identifier{Namespace: "minecraft", Path: "overworld"},
					Data: protocol.PrefixedOptional[any]{
						Present: true,
						Content: overworldData,
					},
				},
			},
		},
	}

	if err := client.SendPacket(overworldPacket); err != nil {
		log.Println("Error while sending dimension_type", err)
	}

	chatData := &defs.ChatType{
		Chat: defs.ChatFormat{
			TranslationKey: "chat.type.text",
			Parameters:     []string{"sender", "content"},
		},
		Narration: defs.ChatFormat{
			TranslationKey: "chat.type.text.narrate",
			Parameters:     []string{"sender", "content"},
		},
	}

	chatPacket := &protocol.RegistryDataPacketOutbound{
		RegistryID: protocol.Identifier{Namespace: "minecraft", Path: "chat_type"},
		Entries: protocol.PrefixedArray[protocol.RegistryEntry]{
			Content: []protocol.RegistryEntry{
				{
					EntryID: protocol.Identifier{Namespace: "minecraft", Path: "chat"},
					Data: protocol.PrefixedOptional[any]{
						Present: true,
						Content: chatData,
					},
				},
			},
		},
	}

	if err := client.SendPacket(chatPacket); err != nil {
		log.Println("Error while sending chat_type", err)
	}

	biomeData := &defs.Biome{
		HasPrecipitation:    1,
		Temperature:         0.5,
		Downfall:            0.5,
		TemperatureModifier: "none",
		Effects: defs.BiomeEffects{
			FogColor:           12638463,
			SkyColor:           7907327,
			WaterColor:         4159280,
			WaterFogColor:      329011,
			GrassColorModifier: "none",
		},
	}

	biomePacket := &protocol.RegistryDataPacketOutbound{
		RegistryID: protocol.Identifier{Namespace: "minecraft", Path: "worldgen/biome"},
		Entries: protocol.PrefixedArray[protocol.RegistryEntry]{
			Content: []protocol.RegistryEntry{
				{
					EntryID: protocol.Identifier{Namespace: "minecraft", Path: "plains"},
					Data: protocol.PrefixedOptional[any]{
						Present: true,
						Content: biomeData,
					},
				},
			},
		},
	}

	if err := client.SendPacket(biomePacket); err != nil {
		log.Println("Error while sending biome", err)
	}

	FinishConfiguration(client)
}

func FinishConfiguration(client *Client) {
	var finishPacket protocol.FinishConfigurationPacketOutbound
	if err := client.SendPacket(&finishPacket); err != nil {
		log.Println("Error while sending finish_configuration", err)
	}
}
