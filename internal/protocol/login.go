package protocol

import (
	"errors"
	"io"
)

type PacketLoginInStart struct {
	Name String
	UUID [16]byte
}

func (l *PacketLoginInStart) ID() VarInt {
	return 0x00
}

func (l *PacketLoginInStart) ReadFrom(r io.Reader) (int64, error) {
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

type PacketLoginOutDisconnect struct {
	Reason String
}

func (l *PacketLoginOutDisconnect) ID() VarInt {
	return 0x00
}

func (l *PacketLoginOutDisconnect) WriteTo(w io.Writer) (int64, error) {
	return l.Reason.WriteTo(w)
}

type PacketLoginOutSuccess struct {
	Profile GameProfile
}

func (l *PacketLoginOutSuccess) ID() VarInt {
	return 0x02
}

func (l *PacketLoginOutSuccess) WriteTo(w io.Writer) (int64, error) {
	return l.Profile.WriteTo(w)
}
