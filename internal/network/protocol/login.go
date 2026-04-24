package protocol

import (
	"errors"
	"io"
)

type PacketLoginInStart struct {
	Name string
	UUID [16]byte
}

func (l *PacketLoginInStart) ID() string {
	return "hello"
}

func (l *PacketLoginInStart) ReadFrom(r io.Reader) (int64, error) {
	n, err := ReadString(r, &l.Name)
	if err != nil {
		return n, err
	}

	if len(l.Name) > 16 {
		return n, errors.New("name is too long")
	}

	l.UUID = [16]byte{}
	m, err := io.ReadFull(r, l.UUID[:])

	return n + int64(m), err
}

type PacketLoginOutDisconnect struct {
	Reason string
}

func (l *PacketLoginOutDisconnect) ID() string {
	return "login_disconnect"
}

type PacketLoginOutSuccess struct {
	Profile GameProfile
}

func (l *PacketLoginOutSuccess) ID() string {
	return "login_finished"
}
