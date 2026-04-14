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

	variantRegistries := []struct {
		path    string
		entryID string
		data    any
	}{
		{"cat_variant", "tabby", defs.CatVariant{
			AssetID:         "minecraft:entity/cat/tabby",
			BabyAssetID:     "minecraft:entity/cat/tabby",
			SpawnConditions: []defs.SpawnCondition{},
		}},
		{"chicken_variant", "default", defs.ChickenVariant{
			AssetID:         "minecraft:entity/chicken",
			BabyAssetID:     "minecraft:entity/chicken",
			Model:           "normal",
			SpawnConditions: []defs.SpawnCondition{},
		}},
		{"cow_variant", "default", defs.CowVariant{
			AssetID:         "minecraft:entity/cow",
			BabyAssetID:     "minecraft:entity/cow",
			Model:           "normal",
			SpawnConditions: []defs.SpawnCondition{},
		}},
		{"pig_variant", "default", defs.PigVariant{
			AssetID:         "minecraft:entity/pig",
			BabyAssetID:     "minecraft:entity/pig",
			Model:           "normal",
			SpawnConditions: []defs.SpawnCondition{},
		}},
		{"wolf_variant", "pale", defs.WolfVariant{
			Assets: defs.WolfAssets{
				Angry: "minecraft:entity/wolf/wolf_angry",
				Wild:  "minecraft:entity/wolf/wolf",
				Tame:  "minecraft:entity/wolf/wolf_tame",
			},
			BabyAssets: defs.WolfAssets{
				Angry: "minecraft:entity/wolf/wolf_angry",
				Wild:  "minecraft:entity/wolf/wolf",
				Tame:  "minecraft:entity/wolf/wolf_tame",
			},
			SpawnConditions: []defs.SpawnCondition{},
		}},
		{"wolf_sound_variant", "default", defs.WolfSoundVariant{
			AdultSounds: defs.WolfSoundSet{
				AmbientSound: "minecraft:entity.wolf.ambient",
				DeathSound:   "minecraft:entity.wolf.death",
				GrowlSound:   "minecraft:entity.wolf.growl",
				HurtSound:    "minecraft:entity.wolf.hurt",
				PantSound:    "minecraft:entity.wolf.pant",
				StepSound:    "minecraft:entity.wolf.step",
				WhineSound:   "minecraft:entity.wolf.whine",
			},
			BabySounds: defs.WolfSoundSet{
				AmbientSound: "minecraft:entity.baby_wolf.ambient",
				DeathSound:   "minecraft:entity.baby_wolf.death",
				GrowlSound:   "minecraft:entity.baby_wolf.growl",
				HurtSound:    "minecraft:entity.baby_wolf.hurt",
				PantSound:    "minecraft:entity.baby_wolf.pant",
				StepSound:    "minecraft:entity.baby_wolf.step",
				WhineSound:   "minecraft:entity.baby_wolf.whine",
			},
		}},
	}

	for _, r := range variantRegistries {
		pkt := &protocol.RegistryDataPacketOutbound{
			RegistryID: protocol.Identifier{Namespace: "minecraft", Path: r.path},
			Entries: protocol.PrefixedArray[protocol.RegistryEntry]{
				Content: []protocol.RegistryEntry{
					{
						EntryID: protocol.Identifier{Namespace: "minecraft", Path: r.entryID},
						Data: protocol.PrefixedOptional[any]{
							Present: true,
							Content: r.data,
						},
					},
				},
			},
		}
		if err := client.SendPacket(pkt); err != nil {
			log.Printf("Error sending registry %s: %v", r.path, err)
			return
		}
	}

	FinishConfiguration(client)
}

func FinishConfiguration(client *Client) {
	var finishPacket protocol.FinishConfigurationPacketOutbound
	if err := client.SendPacket(&finishPacket); err != nil {
		log.Println("Error while sending finish_configuration", err)
	}
}
