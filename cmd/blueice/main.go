package main

import (
	"BlueIce/internal/config"
	"BlueIce/internal/mojang"
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

	if err := os.Mkdir(path+"/lib", 0755); err != nil {
		log.Fatal("Could not create lib directory: ", err)
	}

	if err := mojang.FetchMinecraftData(path + "/lib"); err != nil {
		log.Fatal("Could not fetch minecraft server data: ", err)
	}

	minecraftServer := server.NewMinecraftServer(serverConfig)
	err := minecraftServer.Start()
	if err != nil {
		log.Fatal("Could not start minecraft server", err)
	}
}
