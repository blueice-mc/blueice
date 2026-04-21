package protocol

import (
	"io"
)

type Packet interface {
	ID() VarInt
}

type ServerboundPacket interface {
	Packet
}

type ClientboundPacket interface {
	Packet
}

func WritePacket(w io.Writer, packet ClientboundPacket) (int64, error) {
	return serialize(w, packet)
}
