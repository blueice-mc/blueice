package server

import (
	"BlueIce/internal/protocol"
	"bytes"
	"log"
)

func HandleStatusRequest(client *Client, payload []byte) {
	responseJson := protocol.NewString(`{
        "version": {"name": "26.1.2", "protocol": 775},
        "players": {"max": 100, "online": 5},
        "description": {"text": "§7Another §bBlue§9Ice §7server"}
    }`)

	statusPacket := protocol.StatusPacketOutbound{
		Status: responseJson,
	}

	if err := client.SendPacket(statusPacket); err != nil {
		log.Println(err)
	}
}

func HandlePingRequest(client *Client, payload []byte) {
	var pingPacket protocol.PingPacketInboundOutbound

	if _, err := pingPacket.ReadFrom(bytes.NewBuffer(payload)); err != nil {
		log.Println(err)
		return
	}

	if err := client.SendPacket(pingPacket); err != nil {
		log.Println(err)
	}
}
