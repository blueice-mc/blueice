package protocol

import (
	"BlueIce/internal/nbt"
	"io"
)

type PluginClientMessagePacketOutbound struct {
	Channel Identifier
	Message []byte
}

func (p *PluginClientMessagePacketOutbound) ID() VarInt {
	return 0x01
}

func (p *PluginClientMessagePacketOutbound) WriteTo(w io.Writer) (int64, error) {
	size, err := p.Channel.WriteTo(w)
	if err != nil {
		return size, err
	}
	n, err := w.Write(p.Message)
	return size + int64(n), err
}

func (p *PluginClientMessagePacketOutbound) ReadFrom(r io.Reader) (int64, error) {
	panic("Outbound packet does not support ReadFrom")
}

type RegistryEntry struct {
	EntryID Identifier
	Data    PrefixedOptional[any]
}

type RegistryDataPacketOutbound struct {
	RegistryID Identifier
	Entries    PrefixedArray[RegistryEntry]
}

func (p *RegistryDataPacketOutbound) ID() VarInt {
	return 0x07
}

func (p *RegistryDataPacketOutbound) WriteTo(w io.Writer) (int64, error) {
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

func (p *RegistryDataPacketOutbound) ReadFrom(r io.Reader) (int64, error) {
	panic("Outbound packet does not support ReadFrom")
}

type FinishConfigurationPacketOutbound struct{}

func (p *FinishConfigurationPacketOutbound) ID() VarInt {
	return 0x03
}

func (p *FinishConfigurationPacketOutbound) WriteTo(w io.Writer) (int64, error) {
	return 0, nil
}

func (p *FinishConfigurationPacketOutbound) ReadFrom(r io.Reader) (int64, error) {
	panic("Outbound packet does not support ReadFrom")
}
