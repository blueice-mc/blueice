package adapters

import (
	"bytes"
	"math"

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

func calculateStatesPalettedContainer(s *api.SerializedChunkSection) (int16, protocol.PalettedContainer) {
	types := make(map[uint32]int32)
	latestType := uint32(0)
	nonEmptyCount := int16(0)

	for _, block := range s.Blocks {
		if block != 0 {
			nonEmptyCount++
		}

		types[block]++
		latestType = block
	}

	typeCount := int32(len(types))
	bitsPerType := uint8(0)

	if typeCount == 1 {
		bitsPerType = 0 // 1 type: 0 bits (single value)
	} else if typeCount <= 16 {
		bitsPerType = 4 // 2-16 types: 4 bits (indirect)
	} else if typeCount <= 256 {
		bitsPerType = uint8(math.Ceil(math.Log2(float64(typeCount)))) // 17-256 types: 4-8 bits (indirect)
	} else {
		bitsPerType = 15 // 257+ types: 15 bits (direct)
	}

	blockContainer := protocol.PalettedContainer{
		BitsPerEntry: bitsPerType,
	}

	var palette []protocol.VarInt
	indexMap := make(map[uint32]int8)

	if bitsPerType == 0 {
		// for direct values just set the single value (no palette required)
		blockContainer.SingleValue = protocol.VarInt(latestType)
		return nonEmptyCount, blockContainer
	} else if bitsPerType <= 8 {
		// palette is required
		i := 0
		for block := range types {
			// add type to the palette
			palette = append(palette, protocol.VarInt(block))
			// add index for type to the map
			indexMap[block] = int8(i)
			i++
		}

		// set the palette
		blockContainer.Palette = protocol.PrefixedArray[protocol.VarInt]{
			Content: palette,
		}
	}

	// length of long array is (4096 states * bits per state) / 64 bits per long (= bits per state * 64)
	data := make([]int64, int(bitsPerType)*64)

	mask := uint64(1<<bitsPerType) - 1

	for i, block := range s.Blocks {
		longIndex := (i * int(bitsPerType)) / 64
		bitIndex := (i * int(bitsPerType)) % 64

		if bitIndex+int(bitsPerType) > 64 {
			bitIndex = 0
			longIndex++
		}

	}
}

func calculateSectionContent(s *api.SerializedChunkSection) []byte {
	nonEmpty, blockContainer := calculateStatesPalettedContainer(s)
}

func calculateChunkContent(s *api.SerializedChunk) []byte {
	var buffer bytes.Buffer

	for _, section := range s.Sections {
		buffer.Write(calculateSectionContent(&section))
	}

	return buffer.Bytes()
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
