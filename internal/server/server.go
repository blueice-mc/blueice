package server

import (
	"BlueIce/internal/config"
	"fmt"
	"log"
	"net"
	"sync"
)

const MaximumPacketSize = 2097151

type PacketListener func(*Client, []byte)

type PacketKey struct {
	State    int32
	PacketID uint32
}

type MinecraftServer struct {
	Config config.ServerConfig
	Path   string

	mu              sync.RWMutex
	Clients         []*Client
	PacketListeners map[PacketKey][]PacketListener
}

func NewMinecraftServer(serverConfig config.ServerConfig, path string) *MinecraftServer {
	minecraftServer := MinecraftServer{
		Config:          serverConfig,
		Path:            path,
		Clients:         make([]*Client, 0),
		PacketListeners: make(map[PacketKey][]PacketListener),
	}

	minecraftServer.RegisterPacketListener(0, 0x00, HandleHandshake)
	minecraftServer.RegisterPacketListener(1, 0x00, HandleStatusRequest)
	minecraftServer.RegisterPacketListener(1, 0x01, HandlePingRequest)
	minecraftServer.RegisterPacketListener(2, 0x00, HandleLoginStart)

	return &minecraftServer
}

func (server *MinecraftServer) Start() error {
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

func (server *MinecraftServer) RegisterPacketListener(state int32, packetID uint32, listener func(*Client, []byte)) {
	key := PacketKey{state, packetID}

	server.mu.Lock()
	server.PacketListeners[key] = append(server.PacketListeners[key], listener)
	server.mu.Unlock()
}

func (server *MinecraftServer) onClientConnect(conn net.Conn) {
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
