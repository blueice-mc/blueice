package protocol

import (
	"BlueIce/internal/nbt"
	"io"
)

type PacketConfigOutPluginMessage struct {
	Channel Identifier
	Message []byte
}

func (p *PacketConfigOutPluginMessage) ID() VarInt {
	return 0x01
}

func (p *PacketConfigOutPluginMessage) WriteTo(w io.Writer) (int64, error) {
	size, err := p.Channel.WriteTo(w)
	if err != nil {
		return size, err
	}
	n, err := w.Write(p.Message)
	return size + int64(n), err
}

type RegistryEntry struct {
	EntryID Identifier
	Data    PrefixedOptional[any]
}

type PacketConfigOutRegistryData struct {
	RegistryID Identifier
	Entries    PrefixedArray[RegistryEntry]
}

func (p *PacketConfigOutRegistryData) ID() VarInt {
	return 0x07
}

func (p *PacketConfigOutRegistryData) WriteTo(w io.Writer) (int64, error) {
	size, err := p.RegistryID.WriteTo(w)
	if err != nil {
		return size, err
	}

	p.Entries.Writer = func(w io.Writer, t RegistryEntry) (int64, error) {
		size, err := t.EntryID.WriteTo(w)
		if err != nil {
			return size, err
		}
		t.Data.Writer = func(w io.Writer, t any) (int64, error) {
			return nbt.WriteNBT(w, t)
		}
		n, err := t.Data.WriteTo(w)
		size += int64(n)
		return size, err
	}

	n, err := p.Entries.WriteTo(w)
	size += n
	return size, nil
}

type FinishConfigurationPacketOutbound struct{}

func (p *FinishConfigurationPacketOutbound) ID() VarInt {
	return 0x03
}

func (p *FinishConfigurationPacketOutbound) WriteTo(w io.Writer) (int64, error) {
	return 0, nil
}
