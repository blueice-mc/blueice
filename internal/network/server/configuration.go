package server

import (
	"BlueIce/internal/game/defs"
	"BlueIce/internal/game/registry"
	"BlueIce/internal/network/protocol"
	"bytes"
	"log"
	"sort"
)

func sendRegistryFromMap[T any](client *Client, registryID protocol.Identifier, registry registry.Registry[T]) error {
	keys := make([]string, 0, len(registry.Entries))
	for k := range registry.Entries {
		keys = append(keys, k.Path)
	}
	sort.Strings(keys)

	var registryEntries []protocol.RegistryEntry
	for _, name := range keys {
		id := protocol.NewIdentifier(registryID.Namespace, name)
		data := registry.Entries[id]
		registryEntries = append(registryEntries, protocol.RegistryEntry{
			EntryID: id,
			Data: protocol.PrefixedOptional[protocol.NBTValue]{
				Present: true,
				Content: protocol.NBTValue{Value: &data},
			},
		})
	}

	pkt := &protocol.PacketConfigOutRegistryData{
		RegistryID: registryID,
		Entries:    protocol.PrefixedArray[protocol.RegistryEntry]{Content: registryEntries},
	}

	return client.SendPacket(pkt)
}

func StartConfiguration(client *Client) {
	brand := "BlueIce"
	var buf bytes.Buffer
	length := protocol.VarInt(len(brand))
	if _, err := length.WriteTo(&buf); err != nil {
		log.Println(err)
	}
	buf.WriteString(brand)

	brandPacket := protocol.PacketConfigOutPluginMessage{
		Channel: protocol.Identifier{
			Namespace: "minecraft",
			Path:      "brand",
		},
		Message: buf.Bytes(),
	}

	if err := client.SendPacket(&brandPacket); err != nil {
		log.Println("Error while sending brand: ", err)
		return
	}

	SendRegistryPackets(client)
}

func SendRegistryPackets(client *Client) {

	// Create overworld packet

	clockPacket := protocol.PacketConfigOutRegistryData{
		RegistryID: protocol.Identifier{Namespace: "minecraft", Path: "world_clock"},
		Entries: protocol.PrefixedArray[protocol.RegistryEntry]{
			Content: []protocol.RegistryEntry{
				{
					EntryID: protocol.Identifier{Namespace: "minecraft", Path: "overworld"},
					Data: protocol.PrefixedOptional[protocol.NBTValue]{
						Present: true,
						Content: protocol.NBTValue{Value: defs.WorldClockEntry{}},
					},
				},
			},
		},
	}

	if err := client.SendPacket(&clockPacket); err != nil {
		log.Println("Error while sending clock: ", err)
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

	overworldPacket := &protocol.PacketConfigOutRegistryData{
		RegistryID: protocol.Identifier{Namespace: "minecraft", Path: "dimension_type"},
		Entries: protocol.PrefixedArray[protocol.RegistryEntry]{
			Content: []protocol.RegistryEntry{
				{
					EntryID: protocol.Identifier{Namespace: "minecraft", Path: "overworld"},
					Data: protocol.PrefixedOptional[protocol.NBTValue]{
						Present: true,
						Content: protocol.NBTValue{Value: overworldData},
					},
				},
			},
		},
	}

	if err := client.SendPacket(overworldPacket); err != nil {
		log.Println("Error while sending dimension_type: ", err)
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

	chatPacket := &protocol.PacketConfigOutRegistryData{
		RegistryID: protocol.Identifier{Namespace: "minecraft", Path: "chat_type"},
		Entries: protocol.PrefixedArray[protocol.RegistryEntry]{
			Content: []protocol.RegistryEntry{
				{
					EntryID: protocol.Identifier{Namespace: "minecraft", Path: "chat"},
					Data: protocol.PrefixedOptional[protocol.NBTValue]{
						Present: true,
						Content: protocol.NBTValue{Value: chatData},
					},
				},
			},
		},
	}

	if err := client.SendPacket(chatPacket); err != nil {
		log.Println("Error while sending chat_type: ", err)
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

	biomePacket := &protocol.PacketConfigOutRegistryData{
		RegistryID: protocol.Identifier{Namespace: "minecraft", Path: "worldgen/biome"},
		Entries: protocol.PrefixedArray[protocol.RegistryEntry]{
			Content: []protocol.RegistryEntry{
				{
					EntryID: protocol.Identifier{Namespace: "minecraft", Path: "plains"},
					Data: protocol.PrefixedOptional[protocol.NBTValue]{
						Present: true,
						Content: protocol.NBTValue{Value: biomeData},
					},
				},
			},
		},
	}

	if err := client.SendPacket(biomePacket); err != nil {
		log.Println("Error while sending biome: ", err)
	}

	sendRegistryFromMap(client, protocol.NewIdentifierFromPath("cat_sound_variant"), client.Server.Registries.CatSoundVariant)
	sendRegistryFromMap(client, protocol.NewIdentifierFromPath("cat_variant"), client.Server.Registries.CatVariant)
	sendRegistryFromMap(client, protocol.NewIdentifierFromPath("chicken_sound_variant"), client.Server.Registries.ChickenSoundVariant)
	sendRegistryFromMap(client, protocol.NewIdentifierFromPath("chicken_variant"), client.Server.Registries.ChickenVariant)
	sendRegistryFromMap(client, protocol.NewIdentifierFromPath("cow_sound_variant"), client.Server.Registries.CowSoundVariant)
	sendRegistryFromMap(client, protocol.NewIdentifierFromPath("cow_variant"), client.Server.Registries.CowVariant)
	sendRegistryFromMap(client, protocol.NewIdentifierFromPath("pig_sound_variant"), client.Server.Registries.PigSoundVariant)
	sendRegistryFromMap(client, protocol.NewIdentifierFromPath("pig_variant"), client.Server.Registries.PigVariant)
	sendRegistryFromMap(client, protocol.NewIdentifierFromPath("wolf_sound_variant"), client.Server.Registries.WolfSoundVariant)
	sendRegistryFromMap(client, protocol.NewIdentifierFromPath("wolf_variant"), client.Server.Registries.WolfVariant)
	sendRegistryFromMap(client, protocol.NewIdentifierFromPath("frog_variant"), client.Server.Registries.FrogVariant)
	sendRegistryFromMap(client, protocol.NewIdentifierFromPath("painting_variant"), client.Server.Registries.PaintingVariant)
	sendRegistryFromMap(client, protocol.NewIdentifierFromPath("zombie_nautilus_variant"), client.Server.Registries.ZombieNautilusVariant)
	sendRegistryFromMap(client, protocol.NewIdentifierFromPath("damage_type"), client.Server.Registries.DamageType)
	sendRegistryFromMap(client, protocol.NewIdentifierFromPath("trim_material"), client.Server.Registries.TrimMaterial)
	sendRegistryFromMap(client, protocol.NewIdentifierFromPath("jukebox_song"), client.Server.Registries.JukeboxSong)
	sendRegistryFromMap(client, protocol.NewIdentifierFromPath("banner_pattern"), client.Server.Registries.BannerPattern)
	sendRegistryFromMap(client, protocol.NewIdentifierFromPath("instrument"), client.Server.Registries.Instrument)

	SendTagUpdate(client)
}

func SendTagUpdate(client *Client) {
	var tagUpdatePacket protocol.PacketConfigOutUpdateTags

	type registryTagSource struct {
		registryID protocol.Identifier
		tags       []registry.TagEntry
	}

	sources := []registryTagSource{
		{protocol.Identifier{Namespace: "minecraft", Path: "damage_type"}, client.Server.Registries.DamageType.Tags},
		{protocol.Identifier{Namespace: "minecraft", Path: "banner_pattern"}, client.Server.Registries.BannerPattern.Tags},
	}

	for _, source := range sources {
		registryTags := protocol.RegistryTags{
			Registry: source.registryID,
		}

		for _, tagEntry := range source.tags {
			tag := protocol.Tag{
				TagName: tagEntry.Name,
				Entries: protocol.PrefixedArray[protocol.VarInt]{
					Content: tagEntry.IDs,
				},
			}
			registryTags.Tags.Content = append(registryTags.Tags.Content, tag)
		}

		tagUpdatePacket.TaggedRegistries.Content = append(tagUpdatePacket.TaggedRegistries.Content, registryTags)
	}

	if err := client.SendPacket(&tagUpdatePacket); err != nil {
		log.Println("Error while sending tag_update: ", err)
	}

	FinishConfiguration(client)
}

func FinishConfiguration(client *Client) {
	var finishPacket protocol.PacketConfigOutFinish
	if err := client.SendPacket(&finishPacket); err != nil {
		log.Println("Error while sending finish_configuration: ", err)
	}
}

func HandleConfigurationAcknowledgement(client *Client, payload []byte) {
	client.State = 4
	StartPlay(client)
}
