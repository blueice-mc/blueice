package adapters

import (
	"github.com/blueice-mc/blueice/internal/api"
	"github.com/blueice-mc/blueice/internal/events"
	"github.com/blueice-mc/blueice/internal/game/world"
	"github.com/blueice-mc/blueice/internal/network/protocol"
)

type ChunkAdapter struct {
	EventBus *events.EventBus
}

func NewChunkAdapter(eventBus *events.EventBus) *ChunkAdapter {
	return &ChunkAdapter{
		EventBus: eventBus,
	}
}

func constructPacket(s *api.SerializedChunk) protocol.PacketPlayOutLevelChunkWithLight {
	return protocol.PacketPlayOutLevelChunkWithLight{
		ChunkX: s.X,
		ChunkZ: s.Z,
		Heightmaps: protocol.PrefixedArray[protocol.Heightmap]{
			Content: []protocol.Heightmap{
				{
					Type:        1,
					Data:        s.HeightMap,
					WorldHeight: s.Height,
				},
			},
		},
	}
}

func (a *ChunkAdapter) Convert(chunk *world.Chunk) protocol.PacketPlayOutLevelChunkWithLight {
	serialized := chunk.Serialize()

	event := events.Event{
		Type:    events.PlayerLoadChunk,
		Payload: serialized,
	}

	event, err := a.EventBus.Emit(event)
	if err != nil {
		panic(err)
	}

	return constructPacket(event.Payload.(*api.SerializedChunk))
}
