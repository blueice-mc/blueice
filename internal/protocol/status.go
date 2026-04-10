package protocol

import (
	"encoding/binary"
	"io"
)

type StatusPacketOutbound struct {
	Status String
}

func (packet StatusPacketOutbound) ID() VarInt {
	return 0x00
}

func (packet StatusPacketOutbound) WriteTo(w io.Writer) (int64, error) {
	return packet.Status.WriteTo(w)
}

func (packet StatusPacketOutbound) ReadFrom(r io.Reader) (int64, error) {
	panic("Outbound packet does not support ReadFrom")
}

type PingPacketInboundOutbound struct {
	Timestamp int64
}

func (packet PingPacketInboundOutbound) ID() VarInt {
	return 0x01
}

func (packet PingPacketInboundOutbound) WriteTo(w io.Writer) (int64, error) {
	return 8, binary.Write(w, binary.BigEndian, packet.Timestamp)
}

func (packet PingPacketInboundOutbound) ReadFrom(r io.Reader) (int64, error) {
	return 8, binary.Read(r, binary.BigEndian, &packet.Timestamp)
}
