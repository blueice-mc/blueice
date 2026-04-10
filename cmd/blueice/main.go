package main

import (
	"BlueIce/internal/server"
	"log"
)

func main() {
	log.Println("Starting server")

	minecraftServer := server.NewMinecraftServer(25565)
	err := minecraftServer.Start()

	if err != nil {
		log.Fatal(err)
	}
}
