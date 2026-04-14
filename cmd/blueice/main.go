package main

import (
	"BlueIce/internal/config"
	"BlueIce/internal/server"
	"log"
	"os"
	"path/filepath"
)

func main() {
	log.Println("Starting server")

	ex, _ := os.Executable()
	path := filepath.Dir(ex)

	serverConfig := config.ServerConfig{}

	if err := config.InitializeServerConfig(path+"/config.toml", &serverConfig); err != nil {
		log.Fatal("Could not read or create server config", err)
	}

	minecraftServer := server.NewMinecraftServer(serverConfig, path)
	err := minecraftServer.Start()
	if err != nil {
		log.Fatal("Could not start minecraft server", err)
	}
}
