package config

import (
	_ "embed"
	"log"
	"os"

	"github.com/BurntSushi/toml"
)

type ServerConfig struct {
	Server struct {
		Port       uint16 `toml:"port"`
		MaxPlayers uint16 `toml:"max_players"`
		MOTD       string `toml:"motd"`
	} `toml:"server"`
}

//go:embed default_config.toml
var defaultConfigFile []byte

func InitializeServerConfig(path string) (ServerConfig, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		log.Println("Config file does not exist, exporting default config...")

		err := os.WriteFile(path, defaultConfigFile, 0644)
		if err != nil {
			log.Fatal("Could not create default config", err)
		}
	} else if err != nil {
		log.Fatal("Could not open config file", err)
	}

	var config ServerConfig
	if _, err := toml.DecodeFile(path, &config); err != nil {
		log.Fatal("Could not decode config file", err)
	}

	return config, nil
}
