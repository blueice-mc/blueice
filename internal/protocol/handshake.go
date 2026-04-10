package protocol

import (
	"encoding/binary"
	"io"
)

type HandshakePacketInbound struct {
	ProtocolVersion VarInt
	ServerAddress   String
	ServerPort      uint16

	Intent VarInt
}

func (h *HandshakePacketInbound) ID() VarInt {
	return 0x00
}

func (h *HandshakePacketInbound) ReadFrom(reader io.Reader) (int64, error) {
	var protocolVersion VarInt
	n, err := protocolVersion.ReadFrom(reader)
	size := n
	if err != nil {
		return size, err
	}

	var serverAddress String
	n, err = serverAddress.ReadFrom(reader)
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

	*h = HandshakePacketInbound{
		ProtocolVersion: protocolVersion,
		ServerAddress:   serverAddress,
		ServerPort:      serverPort,
		Intent:          intent,
	}

	return size, nil
}

func (h *HandshakePacketInbound) WriteTo(writer io.Writer) (int64, error) {
	panic("Inbound packet does not support WriteTo")
}
