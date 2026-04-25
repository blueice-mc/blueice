package world

import (
	"log"
	"sync"

	"github.com/blueice-mc/blueice/internal/api"
	"github.com/blueice-mc/blueice/internal/events"
	"github.com/blueice-mc/blueice/internal/network/protocol"
)

type World struct {
	Name          string
	DimensionType protocol.Identifier
	MinY          int16
	Height        uint16
	SeaLevel      int16
	Chunks        map[ChunkPos]*Chunk

	Generator Generator

	genQueue   chan ChunkGenerationRequest
	genWorkers int
	mu         sync.RWMutex
	eventBus   *events.EventBus
}

type ChunkPos struct {
	X, Z int32
}

type ChunkGenerationRequest struct {
	ChunkPos
	Callback chan *Chunk
}

func NewWorld(name string, generator Generator, eventBus *events.EventBus) *World {
	world := &World{
		Name:          name,
		DimensionType: protocol.Identifier{Namespace: "minecraft", Path: "overworld"},
		MinY:          -64,
		Height:        384,
		SeaLevel:      63,
		Chunks:        make(map[ChunkPos]*Chunk),
		Generator:     generator,
		genQueue:      make(chan ChunkGenerationRequest, 100),
		genWorkers:    4,
		eventBus:      eventBus,
	}

	for i := 0; i < world.genWorkers; i++ {
		go world.generationWorker()
	}

	return world
}

func (w *World) generationWorker() {
	for req := range w.genQueue {
		chunk, err := w.GenerateChunk(req.ChunkPos)

		if err != nil {
			log.Println("Failed to generate chunk: ", err, " at ", req.ChunkPos.X, req.ChunkPos.Z)
		}

		w.mu.Lock()
		w.Chunks[req.ChunkPos] = chunk
		w.mu.Unlock()

		if req.Callback != nil {
			req.Callback <- chunk
		}
	}
}

func (w *World) RequestChunkGeneration(x, z int32, callback chan *Chunk) {
	w.mu.RLock()
	// check if the chunk is in memory
	chunk, ok := w.Chunks[ChunkPos{X: x, Z: z}]
	w.mu.RUnlock()

	if ok {
		callback <- chunk
		// if it is, just do nothing
		return
	}

	// if it is not, add it to the queue
	w.genQueue <- ChunkGenerationRequest{
		ChunkPos: ChunkPos{
			X: x,
			Z: z,
		},
		Callback: callback,
	}
}

func (w *World) GetChunk(x, z int32) *Chunk {
	callback := make(chan *Chunk)
	w.RequestChunkGeneration(x, z, callback)
	return <-callback
}

func (w *World) GenerateChunk(pos ChunkPos) (*Chunk, error) {
	serializedChunk := api.SerializedChunk{
		X:        pos.X,
		Z:        pos.Z,
		Height:   w.Height,
		Sections: make([]api.SerializedChunkSection, w.Height/16),
	}

	event := events.Event{
		Type: events.ServerGenerateChunk,
		Payload: api.SerializedChunkGenerationEvent{
			Chunk:     &serializedChunk,
			Generated: false,
		},
	}

	modifiedEvent, err := w.eventBus.Emit(event)
	if err != nil {
		return nil, err
	}

	payload := modifiedEvent.Payload.(api.SerializedChunkGenerationEvent)

	var chunk *Chunk

	if !payload.Generated {
		println("starting generation of chunk at ", pos.X, pos.Z, "")
		chunk = w.Generator.Generate(pos.X, pos.Z)
	} else {
		chunk = &Chunk{}
		chunk.Deserialize(*payload.Chunk)
	}

	return chunk, nil
}

func (w *World) GetBlockState(x, z int32, y int16) uint32 {
	chunk := w.GetChunk(x>>4, z>>4)

	localX := uint8(((x % 16) + 16) % 16)
	localZ := uint8(((z % 16) + 16) % 16)

	return chunk.GetBlockState(localX<<4+localZ, y)
}

func (w *World) SetBlockState(x, z int32, y int16, state uint32) {
	chunk := w.GetChunk(x>>4, z>>4)

	localX := uint8(((x % 16) + 16) % 16)
	localZ := uint8(((z % 16) + 16) % 16)

	chunk.SetBlockState(localX<<4+localZ, y, state)
}
