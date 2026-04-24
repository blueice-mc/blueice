package protocol

import (
	"encoding/binary"
	"io"
)

type PacketHandshakeIn struct {
	ProtocolVersion VarInt
	ServerAddress   string
	ServerPort      uint16

	Intent VarInt
}

func (h *PacketHandshakeIn) ID() string {
	return "intention"
}

func (h *PacketHandshakeIn) ReadFrom(reader io.Reader) (int64, error) {
	var protocolVersion VarInt
	n, err := protocolVersion.ReadFrom(reader)
	size := n
	if err != nil {
		return size, err
	}

	var serverAddress string
	n, err = ReadString(reader, &serverAddress)
	size += n
	if err != nil {
		return size, err
	}

	var serverPort uint16
	err = binary.Read(reader, binary.BigEndian, &serverPort)
	size += 2
	if err != nil {
		return size, err
	}

	var intent VarInt
	n, err = intent.ReadFrom(reader)
	size += n
	if err != nil {
		return size, err
	}

	*h = PacketHandshakeIn{
		ProtocolVersion: protocolVersion,
		ServerAddress:   serverAddress,
		ServerPort:      serverPort,
		Intent:          intent,
	}

	return size, nil
}
