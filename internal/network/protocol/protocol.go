package protocol

import (
	"embed"
	"errors"
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type StateRegistry struct {
	Clientbound map[string]int32 `TOML:"clientbound"`
	Serverbound map[string]int32 `TOML:"serverbound"`
}

var registry map[string]StateRegistry

//go:embed protocol.toml
var protocolFile embed.FS

type ClientState int8

const (
	Handshake ClientState = iota
	Status
	Login
	Configuration
	Play
)

type Direction int8

const (
	Clientbound Direction = iota
	Serverbound
)

func InitializePacketRegistry(path string) error {
	err := os.CopyFS(filepath.Dir(path), protocolFile) // copy the protocol.toml to the specified path

	if err == nil || errors.Is(err, fs.ErrExist) { // if file already exists or no error, decode the file
		if _, err := toml.DecodeFile(path, &registry); err != nil {
			return err
		}

		return nil
	}

	log.Println("Failed to copy protocol.toml:", err)

	// return error otherwise
	return err
}

func GetPacketID(state ClientState, direction Direction, identifier string) int32 {
	if registry == nil {
		panic("Packet registry not initialized")
	}

	stateString := []string{"handshake", "status", "login", "configuration", "play"}[state]

	stateRegistry, ok := registry[stateString]
	if !ok {
		log.Printf("unknown state: %s\n\n", stateString)
		// unknown state, return 0
		return 0
	}

	if direction == Clientbound {
		id, ok := stateRegistry.Clientbound[identifier]
		if !ok {
			log.Printf("unknown packet id: %s\n", identifier)
		}
		return id
	}

	if direction == Serverbound {
		id, ok := stateRegistry.Serverbound[identifier]
		if !ok {
			log.Printf("unknown packet id: %s\n", identifier)
		}
		return id
	}

	// unknown direction, return 0
	return 0
}
