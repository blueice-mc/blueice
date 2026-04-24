package protocol

import (
	"io"
)

type PacketConfigOutPluginMessage struct {
	Channel Identifier
	Message []byte
}

func (p *PacketConfigOutPluginMessage) ID() string {
	return "custom_payload"
}

type RegistryEntry struct {
	EntryID Identifier
	Data    PrefixedOptional[NBTValue]
}

type PacketConfigOutRegistryData struct {
	RegistryID Identifier
	Entries    PrefixedArray[RegistryEntry]
}

func (p *PacketConfigOutRegistryData) ID() string {
	return "registry_data"
}

type PacketConfigOutFinish struct{}

func (p *PacketConfigOutFinish) ID() string {
	return "finish_configuration"
}

func (p *PacketConfigOutFinish) WriteTo(w io.Writer) (int64, error) {
	return 0, nil
}

type Tag struct {
	TagName Identifier
	Entries PrefixedArray[VarInt]
}

type RegistryTags struct {
	Registry Identifier
	Tags     PrefixedArray[Tag]
}

type PacketConfigOutUpdateTags struct {
	TaggedRegistries PrefixedArray[RegistryTags]
}

func (p *PacketConfigOutUpdateTags) ID() string {
	return "update_tags"
}

type PacketConfigInAcknowledged struct{}

func (p *PacketConfigInAcknowledged) ID() string {
	return "finish_configuration"
}

func (p *PacketConfigInAcknowledged) ReadFrom(r io.Reader) (int64, error) {
	return 0, nil
}

type PacketConfigOutDisconnect struct {
	Reason NBTValue
}

func (p *PacketConfigOutDisconnect) ID() string {
	return "disconnect"
}
