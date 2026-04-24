package protocol

import (
	"encoding/binary"
	"io"
)

type PacketStatusOut struct {
	Status string
}

func (packet *PacketStatusOut) ID() string {
	return "status_response"
}

func (packet *PacketStatusOut) WriteTo(w io.Writer) (int64, error) {
	return WriteString(w, packet.Status)
}

type PacketStatusInPing struct {
	Timestamp int64
}

func (packet *PacketStatusInPing) ReadFrom(r io.Reader) (int64, error) {
	var buffer [8]byte
	total, err := r.Read(buffer[:])
	packet.Timestamp = int64(binary.BigEndian.Uint64(buffer[:]))
	return int64(total), err
}

func (packet *PacketStatusInPing) ID() string {
	return "ping_request"
}

type PacketStatusOutPong struct {
	Timestamp int64
}

func (packet *PacketStatusOutPong) ID() string {
	return "pong_response"
}

func (packet *PacketStatusOutPong) WriteTo(w io.Writer) (int64, error) {
	return WriteInt64(w, packet.Timestamp)
}
