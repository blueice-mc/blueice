package world

import (
	"sync"

	"github.com/blueice-mc/blueice/internal/api"
)

type Section struct {
	BlockStates [4096]uint32 // 1x1x1 block state ids
	Biomes      [64]uint32   // 4x4x4 biome ids
}

func (s *Section) GetBlockState(xz uint8, y int16) uint32 {
	return s.BlockStates[int16(xz>>4)+y<<8+int16(xz&0xF)<<4]
}

func (s *Section) GetBiomeAtBlock(xz uint8, y int16) uint32 {
	return s.Biomes[int16(xz>>4)>>2+(y>>2)<<4+(int16(xz&0xF)>>2)<<2]
}

func (s *Section) SetBlockState(xz uint8, y int16, state uint32) {
	s.BlockStates[int16(xz>>4)+y<<8+int16(xz&0xF)<<4] = state
}

func (s *Section) SetBiomeAtBlock(xz uint8, y int16, biome uint32) {
	s.Biomes[int16(xz>>4)>>2+(y>>2)<<4+(int16(xz&0xF)>>2)<<2] = biome
}

type Chunk struct {
	Position ChunkPos
	Sections []Section
	MinY     int16
	Height   uint16

	mu sync.RWMutex
}

func (c *Chunk) GetBlockState(xz uint8, y int16) uint32 {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.Sections[(y-c.MinY)/16].GetBlockState(xz, (y%16+16)%16)
}

func (c *Chunk) GetBiomeAtBlock(xz uint8, y int16) uint32 {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.Sections[(y-c.MinY)/16].GetBiomeAtBlock(xz, (y%16+16)%16)
}

func (c *Chunk) SetBlockState(xz uint8, y int16, state uint32) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.Sections[(y-c.MinY)/16].SetBlockState(xz, (y%16+16)%16, state)
}

func (c *Chunk) SetBiomeAtBlock(xz uint8, y int16, biome uint32) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.Sections[(y-c.MinY)/16].SetBiomeAtBlock(xz, (y%16+16)%16, biome)
}

func (c *Chunk) Serialize() api.SerializedChunk {
	c.mu.RLock()
	defer c.mu.RUnlock()

	serialized := api.SerializedChunk{
		X:        c.Position.X,
		Z:        c.Position.Z,
		MinY:     c.MinY,
		Height:   c.Height,
		Sections: make([]api.SerializedChunkSection, len(c.Sections)),
	}

	for i, section := range c.Sections {
		serialized.Sections[i] = api.SerializedChunkSection{
			Y:      int16(i*16) + c.MinY,
			Blocks: section.BlockStates,
		}
	}

	return serialized
}

func (c *Chunk) Deserialize(serialized api.SerializedChunk) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.Position = ChunkPos{X: serialized.X, Z: serialized.Z}
	c.MinY = serialized.MinY
	c.Height = serialized.Height
	c.Sections = make([]Section, len(serialized.Sections))
	for i, section := range serialized.Sections {
		c.Sections[i] = Section{
			BlockStates: section.Blocks,
		}
	}
}

func (c *Chunk) CalculateHeightmap() [256]int16 {
	heightmap := [256]int16{}

	for xz := uint8(0); uint16(xz) < 0x100; xz++ {
		height := c.MinY - 1

		for y := c.MinY + int16(c.Height) - 1; y >= c.MinY; y-- {
			if c.GetBlockState(xz, y) != 0 {
				height = y
				break
			}
		}

		heightmap[xz] = height
	}

	return heightmap
}
