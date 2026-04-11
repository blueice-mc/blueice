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

func InitializeServerConfig(path string, config *ServerConfig) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		log.Println("Config file does not exist, exporting default config...")

		err := os.WriteFile(path, defaultConfigFile, 0644)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	if _, err := toml.DecodeFile(path, config); err != nil {
		return err
	}

	return nil
}
