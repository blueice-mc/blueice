package main

import (
	"BlueIce/internal/config"
	"BlueIce/internal/mojang"
	"BlueIce/internal/network/server"
	"log"
	"os"
	"path/filepath"
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

	minecraftServer := server.NewMinecraftServer(serverConfig, path)
	err := minecraftServer.Start()
	if err != nil {
		log.Fatal("Could not start minecraft server", err)
	}
}
