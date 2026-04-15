package protocol

import "io"

type Packet interface {
	ID() VarInt
}

type ServerboundPacket interface {
	Packet
	ReadFrom(r io.Reader) (int64, error)
}

type ClientboundPacket interface {
	Packet
	WriteTo(w io.Writer) (int64, error)
}
