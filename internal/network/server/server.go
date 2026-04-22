package server

import (
	"fmt"
	"log"
	"net"
	"os"
	"sync"

	"github.com/blueice-mc/blueice/internal/config"
	"github.com/blueice-mc/blueice/internal/events"
	"github.com/blueice-mc/blueice/internal/game"
	"github.com/blueice-mc/blueice/internal/game/registry"
	"github.com/blueice-mc/blueice/internal/mojang"
)

const MaximumPacketSize = 2097151

type PacketListener func(*Client, []byte)

type PacketKey struct {
	State    int32
	PacketID uint32
}

type NetworkServer struct {
	Config     config.ServerConfig
	Path       string
	GameServer *game.Server

	mu              sync.RWMutex
	Clients         []*Client
	PacketListeners map[PacketKey][]PacketListener
	Registries      registry.Registries
}

func NewNetworkServer(serverConfig config.ServerConfig, path string, eventBus *events.EventBus) *NetworkServer {
	networkServer := NetworkServer{
		Config:          serverConfig,
		Path:            path,
		GameServer:      game.NewServer(eventBus),
		Clients:         make([]*Client, 0),
		PacketListeners: make(map[PacketKey][]PacketListener),
	}

	networkServer.RegisterPacketListener(0, 0x00, HandleHandshake)
	networkServer.RegisterPacketListener(1, 0x00, HandleStatusRequest)
	networkServer.RegisterPacketListener(1, 0x01, HandlePingRequest)
	networkServer.RegisterPacketListener(2, 0x00, HandleLoginStart)
	networkServer.RegisterPacketListener(2, 0x03, HandleLoginAcknowledged)
	networkServer.RegisterPacketListener(3, 0x03, HandleConfigurationAcknowledgement)

	return &networkServer
}

func (server *NetworkServer) Start() error {
	if err := os.Mkdir(server.Path+"/lib", 0755); err == nil {
		err := mojang.FetchMinecraftData(server.Path + "/lib")
		if err != nil {
			log.Fatal("Failed to fetch minecraft server data from mojang: ", err)
		}
	} else if !os.IsExist(err) {
		log.Fatal("Failed to create minecraft server lib directory: ", err)
	}

	if err := server.Registries.LoadAll(server.Path); err != nil {
		log.Fatal("Failed to load minecraft registries: ", err)
	}

	log.Println("Starting minecraft server...")
	err := server.GameServer.Start()
	if err != nil {
		log.Fatal("Failed to start minecraft server: ", err)
	}

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

func (server *NetworkServer) RegisterPacketListener(state int32, packetID uint32, listener func(*Client, []byte)) {
	key := PacketKey{state, packetID}

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
