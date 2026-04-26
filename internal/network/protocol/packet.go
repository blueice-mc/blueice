package protocol

import (
	"io"
)

type Packet interface {
	ID() string
}

type ServerboundPacket interface {
	Packet
}

type ClientboundPacket interface {
	Packet
}

func WritePacket(w io.Writer, packet ClientboundPacket) (int64, error) {
	return Serialize(w, packet)
}
