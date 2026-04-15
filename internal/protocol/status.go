package protocol

import (
	"encoding/binary"
	"io"
)

type PacketStatusOut struct {
	Status String
}

func (packet *PacketStatusOut) ID() VarInt {
	return 0x00
}

func (packet *PacketStatusOut) WriteTo(w io.Writer) (int64, error) {
	return packet.Status.WriteTo(w)
}

type PacketStatusPing struct {
	Timestamp int64
}

func (packet *PacketStatusPing) ID() VarInt {
	return 0x01
}

func (packet *PacketStatusPing) WriteTo(w io.Writer) (int64, error) {
	return 8, binary.Write(w, binary.BigEndian, packet.Timestamp)
}

func (packet *PacketStatusPing) ReadFrom(r io.Reader) (int64, error) {
	return 8, binary.Read(r, binary.BigEndian, &packet.Timestamp)
}
