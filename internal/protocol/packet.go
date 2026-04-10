package protocol

import "io"

type Packet interface {
	ID() VarInt
	WriteTo(w io.Writer) (int64, error)
	ReadFrom(r io.Reader) (int64, error)
}
