package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/blueice-mc/blueice/internal/config"
	"github.com/blueice-mc/blueice/internal/events"
	"github.com/blueice-mc/blueice/internal/game"
	"github.com/blueice-mc/blueice/internal/mojang"
	"github.com/blueice-mc/blueice/internal/network/protocol"
	"github.com/blueice-mc/blueice/internal/network/server"
)

func main() {
	log.Println("Starting server")

	ex, _ := os.Executable()
	path := filepath.Dir(ex)

	if accepted, err := mojang.EulaAccepted(filepath.Join(path, "eula.txt")); err != nil || !accepted {
		log.Println("EULA not accepted. Please accept the EULA before starting the server.")
		if err != nil {
			log.Fatal(err)
		}
		os.Exit(1)
	}

	serverConfig := config.ServerConfig{}

	if err := config.InitializeServerConfig(filepath.Join(path, "config.toml"), &serverConfig); err != nil {
		log.Fatal("Could not read or create server config", err)
	}

	if err := os.Mkdir(filepath.Join(path, "lib"), 0755); err == nil {
		err := mojang.FetchMinecraftData(filepath.Join(path, "lib"))
		if err != nil {
			log.Fatal("Failed to fetch minecraft server data from mojang: ", err)
		}
	} else if !os.IsExist(err) {
		log.Fatal("Failed to create minecraft server lib directory: ", err)
	}

	if err := protocol.InitializePacketRegistry(filepath.Join(path, "lib/protocol.toml")); err != nil {
		log.Fatal("Could not initialize packet registry: ", err)
	}

	eventBus := events.NewBus()

	gameServer := game.NewServer(eventBus)
	if err := gameServer.Start(); err != nil {
		log.Fatal("Could not start minecraft server: ", err)
	}

	// start the game tick loop in a goroutine
	go gameServer.Run()

	// start the tcp server in the current thread
	networkServer := server.NewNetworkServer(serverConfig, path, gameServer)
	if err := networkServer.Start(); err != nil {
		log.Fatal("Could not start TCP server: ", err)
	}
}
