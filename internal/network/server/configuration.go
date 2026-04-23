package server

import (
	"bytes"
	"log"
	"sort"

	"github.com/blueice-mc/blueice/internal/game/defs"
	"github.com/blueice-mc/blueice/internal/game/entity"
	"github.com/blueice-mc/blueice/internal/game/registry"
	"github.com/blueice-mc/blueice/internal/network/protocol"
	"github.com/blueice-mc/blueice/internal/version"
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

	cancelled, reason := client.Server.GameServer.PlayerLogin(&entity.PlayerProfile{
		UUID: client.PendingProfile.UUID,
		Name: client.PendingProfile.Name,
	})

	if cancelled {
		var responsePacket protocol.PacketConfigOutDisconnect
		responsePacket.Reason = protocol.NBTValue{Value: defs.TextComponent{
			Text: reason,
		}}

		if err := client.SendPacket(&responsePacket); err != nil {
			log.Println("Error while sending login response", err)
		}

		return
	}

	brand := "BlueIce " + version.ServerVersion
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
	sendRegistryFromMap(client, protocol.NewIdentifierFromPath("chat_type"), client.Server.GameServer.Registries.ChatType)
	sendRegistryFromMap(client, protocol.NewIdentifierFromPath("world_clock"), client.Server.GameServer.Registries.WorldClock)
	sendRegistryFromMap(client, protocol.NewIdentifierFromPath("dimension_type"), client.Server.GameServer.Registries.DimensionType)
	sendRegistryFromMap(client, protocol.NewIdentifierFromPath("timeline"), client.Server.GameServer.Registries.Timeline)
	sendRegistryFromMap(client, protocol.NewIdentifierFromPath("worldgen/biome"), client.Server.GameServer.Registries.Biomes)
	sendRegistryFromMap(client, protocol.NewIdentifierFromPath("cat_sound_variant"), client.Server.GameServer.Registries.CatSoundVariant)
	sendRegistryFromMap(client, protocol.NewIdentifierFromPath("cat_variant"), client.Server.GameServer.Registries.CatVariant)
	sendRegistryFromMap(client, protocol.NewIdentifierFromPath("chicken_sound_variant"), client.Server.GameServer.Registries.ChickenSoundVariant)
	sendRegistryFromMap(client, protocol.NewIdentifierFromPath("chicken_variant"), client.Server.GameServer.Registries.ChickenVariant)
	sendRegistryFromMap(client, protocol.NewIdentifierFromPath("cow_sound_variant"), client.Server.GameServer.Registries.CowSoundVariant)
	sendRegistryFromMap(client, protocol.NewIdentifierFromPath("cow_variant"), client.Server.GameServer.Registries.CowVariant)
	sendRegistryFromMap(client, protocol.NewIdentifierFromPath("pig_sound_variant"), client.Server.GameServer.Registries.PigSoundVariant)
	sendRegistryFromMap(client, protocol.NewIdentifierFromPath("pig_variant"), client.Server.GameServer.Registries.PigVariant)
	sendRegistryFromMap(client, protocol.NewIdentifierFromPath("wolf_sound_variant"), client.Server.GameServer.Registries.WolfSoundVariant)
	sendRegistryFromMap(client, protocol.NewIdentifierFromPath("wolf_variant"), client.Server.GameServer.Registries.WolfVariant)
	sendRegistryFromMap(client, protocol.NewIdentifierFromPath("frog_variant"), client.Server.GameServer.Registries.FrogVariant)
	sendRegistryFromMap(client, protocol.NewIdentifierFromPath("painting_variant"), client.Server.GameServer.Registries.PaintingVariant)
	sendRegistryFromMap(client, protocol.NewIdentifierFromPath("zombie_nautilus_variant"), client.Server.GameServer.Registries.ZombieNautilusVariant)
	sendRegistryFromMap(client, protocol.NewIdentifierFromPath("damage_type"), client.Server.GameServer.Registries.DamageType)
	sendRegistryFromMap(client, protocol.NewIdentifierFromPath("trim_material"), client.Server.GameServer.Registries.TrimMaterial)
	sendRegistryFromMap(client, protocol.NewIdentifierFromPath("jukebox_song"), client.Server.GameServer.Registries.JukeboxSong)
	sendRegistryFromMap(client, protocol.NewIdentifierFromPath("banner_pattern"), client.Server.GameServer.Registries.BannerPattern)
	sendRegistryFromMap(client, protocol.NewIdentifierFromPath("instrument"), client.Server.GameServer.Registries.Instrument)

	SendTagUpdate(client)
}

func SendTagUpdate(client *Client) {
	var tagUpdatePacket protocol.PacketConfigOutUpdateTags

	type registryTagSource struct {
		registryID protocol.Identifier
		tags       []registry.TagEntry
	}

	sources := []registryTagSource{
		{protocol.Identifier{Namespace: "minecraft", Path: "timeline"}, client.Server.GameServer.Registries.Timeline.Tags},
		{protocol.Identifier{Namespace: "minecraft", Path: "damage_type"}, client.Server.GameServer.Registries.DamageType.Tags},
		{protocol.Identifier{Namespace: "minecraft", Path: "banner_pattern"}, client.Server.GameServer.Registries.BannerPattern.Tags},
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
