package server

import (
	"BlueIce/internal/protocol"
	"bytes"
	"io"
	"log"
	"net"
)

type Client struct {
	Connection net.Conn
	State      int32
	Server     *MinecraftServer
}

func NewClient(conn net.Conn, server *MinecraftServer) *Client {
	return &Client{
		Connection: conn,
		State:      0,
		Server:     server,
	}
}

func (client *Client) Handle() {
	defer client.Connection.Close()

	for {
		var length protocol.VarInt
		if _, err := length.ReadFrom(client.Connection); err != nil {
			log.Println("Read invalid packet length")
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

		log.Printf("Read packet with ID %d", packetID)

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

func (client *Client) SendPacket(packet protocol.Packet) error {
	var buffer bytes.Buffer

	if _, err := packet.ID().WriteTo(&buffer); err != nil {
		return err
	}

	if _, err := packet.WriteTo(&buffer); err != nil {
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
