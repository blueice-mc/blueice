package server

import (
	"bytes"
	"log"
	"strconv"

	"github.com/blueice-mc/blueice/internal/network/protocol"
)

func HandleStatusRequest(client *Client, payload []byte) {
	motd := client.Server.Config.Server.MOTD
	maxPlayers := client.Server.Config.Server.MaxPlayers

	responseJson := `{
        "version": {"name": "26.1.2", "protocol": 775},
        "players": {"max": ` + strconv.Itoa(int(maxPlayers)) + `, "online": 0},
        "description": {"text": "` + motd + `"}
    }`

	statusPacket := protocol.PacketStatusOut{
		Status: responseJson,
	}

	if err := client.SendPacket(&statusPacket); err != nil {
		log.Println(err)
	}
}

func HandlePingRequest(client *Client, payload []byte) {
	var pingPacket protocol.PacketStatusInPing

	if _, err := pingPacket.ReadFrom(bytes.NewBuffer(payload)); err != nil {
		log.Println(err)
		return
	}

	pongPacket := protocol.PacketStatusOutPong{
		Timestamp: pingPacket.Timestamp,
	}

	if err := client.SendPacket(&pongPacket); err != nil {
		log.Println(err)
	}
}
