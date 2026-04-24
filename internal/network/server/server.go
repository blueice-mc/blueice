package server

import (
	"fmt"
	"log"
	"net"
	"sync"

	"github.com/blueice-mc/blueice/internal/config"
	"github.com/blueice-mc/blueice/internal/game"
	"github.com/blueice-mc/blueice/internal/network/protocol"
)

const MaximumPacketSize = 2097151

type PacketListener func(*Client, []byte)

type PacketKey struct {
	State    protocol.ClientState
	PacketID uint32
}

type NetworkServer struct {
	Config     config.ServerConfig
	Path       string
	GameServer *game.Server

	mu              sync.RWMutex
	Clients         []*Client
	PacketListeners map[PacketKey][]PacketListener
}

func NewNetworkServer(serverConfig config.ServerConfig, path string, gameServer *game.Server) *NetworkServer {
	networkServer := NetworkServer{
		Config:          serverConfig,
		Path:            path,
		GameServer:      gameServer,
		Clients:         make([]*Client, 0),
		PacketListeners: make(map[PacketKey][]PacketListener),
	}

	networkServer.RegisterPacketListener(protocol.Handshake, "intention", HandleHandshake)
	networkServer.RegisterPacketListener(protocol.Status, "status_request", HandleStatusRequest)
	networkServer.RegisterPacketListener(protocol.Status, "ping_request", HandlePingRequest)
	networkServer.RegisterPacketListener(protocol.Login, "hello", HandleLoginStart)
	networkServer.RegisterPacketListener(protocol.Login, "login_acknowledged", HandleLoginAcknowledged)
	networkServer.RegisterPacketListener(protocol.Configuration, "finish_configuration", HandleConfigurationAcknowledgement)

	return &networkServer
}

func (server *NetworkServer) Start() error {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", server.Config.Server.Port))
	log.Println("Listening on", listener.Addr())

	if err != nil {
		return err
	}

	for {
		conn, err := listener.Accept()

		log.Printf("Accepting connection from %s", conn.RemoteAddr())

		if err != nil {
			err := conn.Close()
			if err != nil {
			}
		}

		go server.onClientConnect(conn)
	}
}

func (server *NetworkServer) RegisterPacketListener(state protocol.ClientState, packetName string, listener func(*Client, []byte)) {
	id := protocol.GetPacketID(state, protocol.Serverbound, packetName)
	key := PacketKey{state, uint32(id)}

	server.mu.Lock()
	server.PacketListeners[key] = append(server.PacketListeners[key], listener)
	server.mu.Unlock()
}

func (server *NetworkServer) onClientConnect(conn net.Conn) {
	client := NewClient(conn, server)

	server.mu.Lock()
	server.Clients = append(server.Clients, client)
	server.mu.Unlock()

	client.Handle()

	server.mu.Lock()
	defer server.mu.Unlock()

	for i, c := range server.Clients {
		if c == client {
			server.Clients = append(server.Clients[:i], server.Clients[i+1:]...)
			log.Printf("Client disconnected.")
			break
		}
	}
}
