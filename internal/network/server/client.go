package server

import (
	"bytes"
	"io"
	"log"
	"net"

	"github.com/blueice-mc/blueice/internal/game/entity"
	"github.com/blueice-mc/blueice/internal/network/protocol"
)

type Client struct {
	Connection     net.Conn
	State          protocol.ClientState
	Server         *NetworkServer
	PendingProfile *protocol.GameProfile
	Player         *entity.Player
}

func NewClient(conn net.Conn, server *NetworkServer) *Client {
	return &Client{
		Connection: conn,
		State:      0,
		Server:     server,
	}
}

func (client *Client) GetAddress() string {
	return client.Connection.RemoteAddr().String()
}

func (client *Client) Handle() {
	defer client.Connection.Close()

	for {
		var length protocol.VarInt
		if _, err := length.ReadFrom(client.Connection); err != nil {
			if err != io.EOF {
				log.Println("Read invalid packet length: ", err)
			}
			break
		}

		if length == 0 {
			log.Println("Client disconnected: Received empty packet")
			break
		}

		if length > MaximumPacketSize {
			log.Println("Packet is too large")
			break
		}

		var packetID protocol.VarInt
		n2, err := packetID.ReadFrom(client.Connection)
		if err != nil {
			break
		}

		payloadLength := int(length) - int(n2)

		if payloadLength < 0 {
			log.Println("Negative payload length")
		}

		buffer := make([]byte, payloadLength)
		if _, err := io.ReadFull(client.Connection, buffer); err != nil {
			log.Println("Error while reading packet")
			break
		}

		key := PacketKey{
			State:    client.State,
			PacketID: uint32(packetID),
		}

		client.Server.mu.RLock()
		listeners, ok := client.Server.PacketListeners[key]
		client.Server.mu.RUnlock()

		if ok {
			for _, listener := range listeners {
				listener(client, buffer)
			}
		}
	}
}

func (client *Client) SendPacket(packet protocol.ClientboundPacket) error {
	var buffer bytes.Buffer

	id := protocol.VarInt(protocol.GetPacketID(client.State, protocol.Clientbound, packet.ID()))

	if _, err := id.WriteTo(&buffer); err != nil {
		return err
	}

	if _, err := protocol.WritePacket(&buffer, packet); err != nil {
		return err
	}

	length := protocol.VarInt(buffer.Len())
	if _, err := length.WriteTo(client.Connection); err != nil {
		return err
	}

	if _, err := buffer.WriteTo(client.Connection); err != nil {
		return err
	}

	return nil
}
