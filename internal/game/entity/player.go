package entity

import (
	"log"

	"github.com/blueice-mc/blueice/internal/api"
	"github.com/blueice-mc/blueice/internal/game/defs"
	"github.com/blueice-mc/blueice/internal/network/protocol"
)

type Connection interface {
	SendPacket(packet protocol.ClientboundPacket) error
	GetAddress() string
}

type Player struct {
	Entity
	PlayerProfile

	X, Y, Z    float64
	Yaw, Pitch float32

	Gamemode uint8
	Health   float32
	Food     int32

	Connection Connection
}

func (p *Player) Serialize() api.SerializedPlayer {
	return api.SerializedPlayer{
		UUID:     p.UUID,
		Name:     p.Name,
		X:        p.X,
		Y:        p.Y,
		Z:        p.Z,
		Yaw:      p.Yaw,
		Pitch:    p.Pitch,
		Gamemode: p.Gamemode,
		Health:   p.Health,
	}
}

func (p *Player) Deserialize(serialized api.SerializedPlayer) {
	p.UUID = serialized.UUID
	p.Name = serialized.Name
	p.X = serialized.X
	p.Y = serialized.Y
	p.Z = serialized.Z
	p.Yaw = serialized.Yaw
	p.Pitch = serialized.Pitch
	p.Gamemode = serialized.Gamemode
	p.Health = serialized.Health
}

func (p *Player) Kick(reason defs.TextComponent) {
	disconnectPacket := &protocol.PacketPlayOutDisconnect{Reason: protocol.NBTValue{Value: reason}}
	if err := p.Connection.SendPacket(disconnectPacket); err != nil {
		log.Println("Error while sending disconnect packet: ", err)
	}
}

type PlayerProfile struct {
	UUID [16]byte
	Name string
}

func (p *PlayerProfile) Serialize() api.SerializedGameProfile {
	return api.SerializedGameProfile{
		UUID: p.UUID,
		Name: p.Name,
	}
}
