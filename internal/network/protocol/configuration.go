package protocol

import (
	"io"
)

type PacketConfigOutPluginMessage struct {
	Channel Identifier
	Message []byte
}

func (p *PacketConfigOutPluginMessage) ID() VarInt {
	return 0x01
}

type RegistryEntry struct {
	EntryID Identifier
	Data    PrefixedOptional[NBTValue]
}

type PacketConfigOutRegistryData struct {
	RegistryID Identifier
	Entries    PrefixedArray[RegistryEntry]
}

func (p *PacketConfigOutRegistryData) ID() VarInt {
	return 0x07
}

type PacketConfigOutFinish struct{}

func (p *PacketConfigOutFinish) ID() VarInt {
	return 0x03
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

func (p *PacketConfigOutUpdateTags) ID() VarInt {
	return 0x0D
}

type PacketConfigInAcknowledged struct{}

func (p *PacketConfigInAcknowledged) ID() VarInt {
	return 0x03
}

func (p *PacketConfigInAcknowledged) ReadFrom(r io.Reader) (int64, error) {
	return 0, nil
}

type PacketConfigOutDisconnect struct {
	Reason NBTValue
}

func (p *PacketConfigOutDisconnect) ID() VarInt {
	return 0x02
}
