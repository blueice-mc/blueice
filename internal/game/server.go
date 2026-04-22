package game

import (
	"log"
	"sync"
	"sync/atomic"

	"github.com/blueice-mc/blueice/internal/api"
	"github.com/blueice-mc/blueice/internal/events"
	"github.com/blueice-mc/blueice/internal/game/entity"
	"github.com/blueice-mc/blueice/internal/game/world"
)

type Server struct {
	players  map[[16]byte]*entity.Player // map of all online players by UUID
	worlds   map[[16]byte]*world.World   // map of all loaded worlds by UUID
	eventBus *events.EventBus
	entityID atomic.Int32
	mu       sync.RWMutex
}

func NewServer(eventBus *events.EventBus) *Server {
	server := &Server{
		players:  make(map[[16]byte]*entity.Player),
		worlds:   make(map[[16]byte]*world.World),
		eventBus: eventBus,
	}

	// begin with entity ID 1
	server.entityID.Store(1)

	return server
}

// Start starts the game server. It fires the corresponding lifecycle events.
func (s *Server) Start() error {
	_, err := s.eventBus.Emit(events.Event{
		Type:    events.ServerLifecycleStarting,
		Payload: nil,
	})

	if err != nil {
		return err
	}

	// loading logic belongs here

	_, err = s.eventBus.Emit(events.Event{
		Type: events.ServerLifecycleStarted,
	})

	return err
}

// Stop stops the game server. It fires the corresponding lifecycle events.
func (s *Server) Stop() error {
	_, err := s.eventBus.Emit(events.Event{
		Type:    events.ServerLifecycleStopping,
		Payload: nil,
	})

	// unloading logic belongs here

	if err != nil {
		return err
	}

	_, err = s.eventBus.Emit(events.Event{
		Type:    events.ServerLifecycleStopped,
		Payload: nil,
	})

	return err
}

// PlayerLogin triggers a LoginEvent and returns if the login was rejected and a reason if the login was rejected.
func (s *Server) PlayerLogin(profile *entity.PlayerProfile) (bool, string) {
	loginEvent := api.SerializedLoginEvent{
		UUID:          profile.UUID,
		Name:          profile.Name,
		Cancelled:     false,
		CancelMessage: "",
	}

	modifiedEvent, err := s.eventBus.Emit(events.Event{
		Type:    events.PlayerLogin,
		Payload: loginEvent,
	})

	if err != nil {
		log.Println("Error while emitting PlayerLogin event", err)
		return true, "Internal server error."
	}

	return modifiedEvent.Payload.(api.SerializedLoginEvent).Cancelled,
		modifiedEvent.Payload.(api.SerializedLoginEvent).CancelMessage
}

// AddPlayer adds a player to the server and fires a player join event.
func (s *Server) AddPlayer(player *entity.Player) error {
	s.mu.Lock()
	s.players[player.UUID] = player
	s.mu.Unlock()

	_, err := s.eventBus.Emit(events.Event{
		Type:    events.PlayerJoin,
		Payload: player.Serialize(),
	})

	return err
}
