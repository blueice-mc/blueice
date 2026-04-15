package server

import (
	"BlueIce/internal/protocol"
	"bytes"
	"log"
	"strconv"
)

func HandleStatusRequest(client *Client, payload []byte) {
	motd := client.Server.Config.Server.MOTD
	maxPlayers := client.Server.Config.Server.MaxPlayers

	responseJson := protocol.NewString(`{
        "version": {"name": "26.1.2", "protocol": 775},
        "players": {"max": ` + strconv.Itoa(int(maxPlayers)) + `, "online": 0},
        "description": {"text": "` + motd + `"}
    }`)

	statusPacket := protocol.PacketStatusOut{
		Status: responseJson,
	}

	if err := client.SendPacket(&statusPacket); err != nil {
		log.Println(err)
	}
}

func HandlePingRequest(client *Client, payload []byte) {
	var pingPacket protocol.PacketStatusPing

	if _, err := pingPacket.ReadFrom(bytes.NewBuffer(payload)); err != nil {
		log.Println(err)
		return
	}

	if err := client.SendPacket(&pingPacket); err != nil {
		log.Println(err)
	}
}
