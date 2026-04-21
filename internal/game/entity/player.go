package entity

import (
	"BlueIce/internal/game/defs"
	"BlueIce/internal/network/protocol"
	"log"
)

type Connection interface {
	SendPacket(packet protocol.ClientboundPacket) error
	GetAddress() string
}

type Player struct {
	Entity
	UUID       [16]byte
	PlayerName string

	Connection Connection
}

func (p *Player) Kick(reason defs.TextComponent) {
	disconnectPacket := &protocol.PacketPlayOutDisconnect{Reason: protocol.NBTValue{Value: reason}}
	if err := p.Connection.SendPacket(disconnectPacket); err != nil {
		log.Println("Error while sending disconnect packet: ", err)
	}
}
