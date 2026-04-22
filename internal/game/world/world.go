package world

import (
	"github.com/blueice-mc/blueice/internal/network/protocol"
)

type World struct {
	Name          string
	DimensionType protocol.Identifier
	MinY          int16
	Height        uint16
	SeaLevel      int16
	Chunks        map[ChunkPos]*Chunk
}

type ChunkPos struct {
	X, Z int32
}

func (w *World) GetBlockState(x, z int32, y int16) uint32 {
	pos := ChunkPos{
		X: x >> 4,
		Z: z >> 4,
	}

	chunk, ok := w.Chunks[pos]
	if !ok {
		return 0 // Return air if chunk does not exist
	}

	localX := uint8(((x % 16) + 16) % 16)
	localZ := uint8(((z % 16) + 16) % 16)

	return chunk.GetBlockState(localX<<4+localZ, y)
}

func (w *World) SetBlockState(x, z int32, y int16, state uint32) {
	pos := ChunkPos{
		X: x >> 4,
		Z: z >> 4,
	}

	chunk, ok := w.Chunks[pos]
	if !ok {
		panic("chunk generation not implemented yet")
	}

	localX := uint8(((x % 16) + 16) % 16)
	localZ := uint8(((z % 16) + 16) % 16)

	chunk.SetBlockState(localX<<4+localZ, y, state)
}
