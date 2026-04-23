package events

import (
	"log"
	"sync"
)

type EventBus struct {
	handlers map[EventType]map[EventPriority][]func(Event) (Event, error)
	mu       sync.RWMutex
}

type EventType string

const (
	// Player events
	PlayerLogin     EventType = "player.login"
	PlayerJoin      EventType = "player.join"
	PlayerQuit      EventType = "player.quit"
	PlayerLoadChunk EventType = "player.load_chunk"

	// Server events
	ServerTick              EventType = "server.tick"
	ServerLifecycleStarting EventType = "server.lifecycle_starting"
	ServerLifecycleStarted  EventType = "server.lifecycle_started"
	ServerLifecycleStopping EventType = "server.lifecycle_stopping"
	ServerLifecycleStopped  EventType = "server.lifecycle_stopped"
	ServerGenerateChunk     EventType = "server.generate_chunk"
)

type EventPriority uint8

const (
	PriorityLowest EventPriority = iota
	PriorityLow
	PriorityNormal
	PriorityHigh
	PriorityHighest
)

type Event struct {
	Type    EventType
	Payload interface{}
}

func NewBus() *EventBus {
	return &EventBus{
		handlers: make(map[EventType]map[EventPriority][]func(Event) (Event, error)),
	}
}

func (e *EventBus) Emit(event Event) (Event, error) {
	e.mu.RLock()
	handlers, ok := e.handlers[event.Type]
	e.mu.RUnlock()
	currentEvent := event
	if ok {
		for priority := PriorityLowest; priority <= PriorityHighest; priority++ {
			for _, handler := range handlers[priority] {
				evt, err := handler(currentEvent)
				if err != nil {
					log.Printf("Error while handling %s: %v", event.Type, err)
					continue
				}
				currentEvent = evt
			}
		}
	}
	return currentEvent, nil
}

func (e *EventBus) Subscribe(eventType EventType, eventPriority EventPriority, handler func(Event) (Event, error)) {
	e.mu.Lock()
	handlers, ok := e.handlers[eventType]
	if !ok {
		handlers = make(map[EventPriority][]func(Event) (Event, error))
		e.handlers[eventType] = handlers
	}
	e.handlers[eventType][eventPriority] = append(handlers[eventPriority], handler)
	e.mu.Unlock()
}

func (e *EventBus) HasSubscribers(eventType EventType) bool {
	e.mu.RLock()
	defer e.mu.RUnlock()
	handlers, ok := e.handlers[eventType]
	return ok && len(handlers) > 0
}
