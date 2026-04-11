package protocol

import (
	"errors"
	"io"
)

type LoginStartPacketInbound struct {
	Name String
	UUID [16]byte
}

func (l LoginStartPacketInbound) ID() VarInt {
	return 0x00
}

func (l LoginStartPacketInbound) WriteTo(w io.Writer) (int64, error) {
	panic("Inbound packet does not support WriteTo")
}

func (l *LoginStartPacketInbound) ReadFrom(r io.Reader) (int64, error) {
	n, err := l.Name.ReadFrom(r)
	if err != nil {
		return n, err
	}

	if l.Name.Length > 16 {
		return n, errors.New("name is too long")
	}

	l.UUID = [16]byte{}
	m, err := io.ReadFull(r, l.UUID[:])

	return n + int64(m), err
}

type LoginDisconnectPacketOutbound struct {
	Reason String
}

func (l LoginDisconnectPacketOutbound) ID() VarInt {
	return 0x00
}

func (l LoginDisconnectPacketOutbound) WriteTo(w io.Writer) (int64, error) {
	return l.Reason.WriteTo(w)
}

func (l *LoginDisconnectPacketOutbound) ReadFrom(r io.Reader) (int64, error) {
	panic("Outbound packet does not support ReadFrom")
}

type LoginSuccessPacketOutbound struct {
	Profile GameProfile
}

func (l LoginSuccessPacketOutbound) ID() VarInt {
	return 0x02
}

func (l *LoginSuccessPacketOutbound) WriteTo(w io.Writer) (int64, error) {
	return l.Profile.WriteTo(w)
}

func (l *LoginSuccessPacketOutbound) ReadFrom(r io.Reader) (int64, error) {
	panic("Outbound packet does not support ReadFrom")
}
