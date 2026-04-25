package world

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/blueice-mc/blueice/internal/game/block"
)

type GeneratorConfig struct {
	Type   string
	Preset string
	MinY   int16
	Height uint16
}

type Generator interface {
	Generate(x, z int32) *Chunk
}

type FlatGenerator struct {
	Config       *GeneratorConfig
	ParsedPreset *FlatPreset
}

type FlatLayer struct {
	Block string
	Count uint16
}

type FlatPreset struct {
	Layers  []FlatLayer
	BiomeID uint32
}

func NewFlatGenerator(config *GeneratorConfig, biomes map[string]uint32) (*FlatGenerator, error) {
	if config.Type != "flat" || config.Preset == "" {
		panic("invalid generator config")
	}

	generator := &FlatGenerator{
		Config: config,
	}

	preset, err := parseFlatPreset(config.Preset, biomes)

	if err != nil {
		return nil, err
	}

	generator.ParsedPreset = preset

	return generator, nil
}

func parseFlatPreset(preset string, biomes map[string]uint32) (*FlatPreset, error) {
	parts := strings.Split(preset, ";")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid preset format: expected 'layers;biome'")
	}

	layersStr := parts[0]
	biome := parts[1]

	layerParts := strings.Split(layersStr, ",")
	layers := make([]FlatLayer, 0, len(layerParts))

	for _, layerStr := range layerParts {
		layer, err := parseLayer(layerStr)
		if err != nil {
			return nil, err
		}
		layers = append(layers, layer)
	}

	return &FlatPreset{
		Layers:  layers,
		BiomeID: biomes[biome],
	}, nil
}

func parseLayer(layerStr string) (FlatLayer, error) {
	if strings.Contains(layerStr, "*") {
		parts := strings.Split(layerStr, "*")
		if len(parts) != 2 {
			return FlatLayer{}, fmt.Errorf("invalid layer format: %s", layerStr)
		}

		count, err := strconv.Atoi(parts[0])
		if err != nil {
			return FlatLayer{}, fmt.Errorf("invalid count: %s", parts[0])
		}

		return FlatLayer{
			Block: parts[1],
			Count: uint16(count),
		}, nil
	}

	return FlatLayer{
		Block: layerStr,
		Count: 1,
	}, nil
}

func (g *FlatGenerator) Generate(x, z int32) *Chunk {
	preset := g.ParsedPreset

	chunk := Chunk{
		Position: ChunkPos{
			X: x,
			Z: z,
		},
		Height:   g.Config.Height,
		MinY:     g.Config.MinY,
		Sections: make([]Section, g.Config.Height/16),
	}

	for i := range chunk.Sections {
		chunk.Sections[i] = Section{}

		for biomeIdx := 0; biomeIdx < 64; biomeIdx++ {
			chunk.mu.Lock()
			chunk.Sections[i].Biomes[biomeIdx] = preset.BiomeID
			chunk.mu.Unlock()
		}
	}

	currentY := g.Config.MinY

	for _, layer := range preset.Layers {
		stateId := block.BlockStates[layer.Block]

		for xz := uint16(0); xz < 0x100; xz++ {
			for y := currentY; y < currentY+int16(layer.Count); y++ {
				chunk.SetBlockState(uint8(xz), y, stateId)
			}
		}

		currentY += int16(layer.Count)
	}

	return &chunk
}
