package server

import (
	"BlueIce/internal/defs"
	"BlueIce/internal/protocol"
	"bytes"
	"log"
	"sort"
)

func sendRegistryFromMap[T any](client *Client, registryPath string, entries map[string]T) error {
	keys := make([]string, 0, len(entries))
	for k := range entries {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var registryEntries []protocol.RegistryEntry
	for _, name := range keys {
		data := entries[name]
		registryEntries = append(registryEntries, protocol.RegistryEntry{
			EntryID: protocol.Identifier{Namespace: "minecraft", Path: name},
			Data: protocol.PrefixedOptional[any]{
				Present: true,
				Content: &data,
			},
		})
	}

	pkt := &protocol.RegistryDataPacketOutbound{
		RegistryID: protocol.Identifier{Namespace: "minecraft", Path: registryPath},
		Entries:    protocol.PrefixedArray[protocol.RegistryEntry]{Content: registryEntries},
	}

	return client.SendPacket(pkt)
}

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

	clockPacket := protocol.RegistryDataPacketOutbound{
		RegistryID: protocol.Identifier{Namespace: "minecraft", Path: "world_clock"},
		Entries: protocol.PrefixedArray[protocol.RegistryEntry]{
			Content: []protocol.RegistryEntry{
				{
					EntryID: protocol.Identifier{Namespace: "minecraft", Path: "overworld"},
					Data: protocol.PrefixedOptional[any]{
						Present: true,
						Content: defs.WorldClockEntry{},
					},
				},
			},
		},
	}

	if err := client.SendPacket(&clockPacket); err != nil {
		log.Println("Error while sending clock", err)
		return
	}

	overworldData := &defs.DimensionTypeEntry{
		CoordinateScale:     1.0,
		HasSkylight:         1,
		HasCeiling:          0,
		HasEnderDragonFight: 0,
		AmbientLight:        0.0,

		HasFixedTime: 0,
		FixedTime:    nil,

		MonsterSpawnBlockLightLimit: 0,
		MonsterSpawnLightLevel: defs.IntProvider{
			Type:  "minecraft:constant",
			Value: 0,
		},

		LogicalHeight: 384,
		MinY:          -64,
		Height:        384,

		Infiniburn: "#minecraft:infiniburn_overworld",

		Skybox:        "overworld",
		CardinalLight: "default",

		Attributes: defs.EmptyCompound{},

		DefaultClock: "minecraft:overworld",

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

	sendRegistryFromMap(client, "cat_sound_variant", client.Server.Registries.CatSoundVariant)
	sendRegistryFromMap(client, "cat_variant", client.Server.Registries.CatVariant)
	sendRegistryFromMap(client, "chicken_sound_variant", client.Server.Registries.ChickenSoundVariant)
	sendRegistryFromMap(client, "chicken_variant", client.Server.Registries.ChickenVariant)
	sendRegistryFromMap(client, "cow_sound_variant", client.Server.Registries.CowSoundVariant)
	sendRegistryFromMap(client, "cow_variant", client.Server.Registries.CowVariant)
	sendRegistryFromMap(client, "pig_sound_variant", client.Server.Registries.PigSoundVariant)
	sendRegistryFromMap(client, "pig_variant", client.Server.Registries.PigVariant)
	sendRegistryFromMap(client, "wolf_sound_variant", client.Server.Registries.WolfSoundVariant)
	sendRegistryFromMap(client, "wolf_variant", client.Server.Registries.WolfVariant)
	sendRegistryFromMap(client, "frog_variant", client.Server.Registries.FrogVariant)
	sendRegistryFromMap(client, "painting_variant", client.Server.Registries.PaintingVariant)
	sendRegistryFromMap(client, "zombie_nautilus_variant", client.Server.Registries.ZombieNautilusVariant)
	sendRegistryFromMap(client, "damage_type", client.Server.Registries.DamageType)

	FinishConfiguration(client)
}

func FinishConfiguration(client *Client) {
	var finishPacket protocol.FinishConfigurationPacketOutbound
	if err := client.SendPacket(&finishPacket); err != nil {
		log.Println("Error while sending finish_configuration", err)
	}
}
