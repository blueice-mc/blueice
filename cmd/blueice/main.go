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

	serverConfig, _ := config.InitializeServerConfig(path)

	log.Println(serverConfig)

	minecraftServer := server.NewMinecraftServer(25565)
	err := minecraftServer.Start()

	if err != nil {
		log.Fatal(err)
	}
}
