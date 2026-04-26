package protocol

import (
	"encoding/binary"
	"io"
)

type PacketPlayOutLogin struct {
	EntityID            int32
	IsHardcore          bool
	DimensionNames      PrefixedArray[Identifier]
	MaxPlayers          VarInt
	ViewDistance        VarInt
	SimulationDistance  VarInt
	ReducedDebugInfo    bool
	EnableRespawnScreen bool
	DoLimitedCrafting   bool
	DimensionType       VarInt
	DimensionName       Identifier
	HashedSeed          int64
	GameMode            int8
	PreviousGameMode    int8
	IsDebug             bool
	IsFlat              bool
	HasDeathLocation    bool
	DeathDimensionName  *Identifier
	DeathLocation       *Position
	PortalCooldown      VarInt
	SeaLevel            VarInt
	EnforcesSecureChat  bool
}

func (p *PacketPlayOutLogin) ID() string {
	return "login"
}

func (p *PacketPlayOutLogin) WriteTo(w io.Writer) (int64, error) {
	total, err := Serialize(w, struct {
		EntityID            int32
		IsHardcore          bool
		DimensionNames      PrefixedArray[Identifier]
		MaxPlayers          VarInt
		ViewDistance        VarInt
		SimulationDistance  VarInt
		ReducedDebugInfo    bool
		EnableRespawnScreen bool
		DoLimitedCrafting   bool
		DimensionType       VarInt
		DimensionName       Identifier
		HashedSeed          int64
		GameMode            int8
		PreviousGameMode    int8
		IsDebug             bool
		IsFlat              bool
		HasDeathLocation    bool
	}{
		p.EntityID, p.IsHardcore, p.DimensionNames, p.MaxPlayers,
		p.ViewDistance, p.SimulationDistance, p.ReducedDebugInfo,
		p.EnableRespawnScreen, p.DoLimitedCrafting, p.DimensionType,
		p.DimensionName, p.HashedSeed, p.GameMode, p.PreviousGameMode,
		p.IsDebug, p.IsFlat, p.HasDeathLocation,
	})
	if err != nil {
		return total, err
	}

	if p.HasDeathLocation {
		n, err := Serialize(w, p.DeathDimensionName)
		total += n
		if err != nil {
			return total, err
		}
		n, err = Serialize(w, p.DeathLocation)
		total += n
		if err != nil {
			return total, err
		}
	}

	n, err := Serialize(w, struct {
		PortalCooldown     VarInt
		SeaLevel           VarInt
		EnforcesSecureChat bool
	}{p.PortalCooldown, p.SeaLevel, p.EnforcesSecureChat})
	total += n
	return total, err
}

type PacketPlayOutGameEvent struct {
	Event uint8
	Value float32
}

func (p *PacketPlayOutGameEvent) ID() string {
	return "game_event"
}

type Section struct {
	NonEmptyBlockCount VarInt
	BlockStates        PalettedContainer
	Biomes             PalettedContainer
}

type ContainerType int8

const (
	BlockStates ContainerType = iota
	Biomes
)

type PalettedContainer struct {
	BitsPerEntry uint8
	Palette      PrefixedArray[VarInt]
	Storage      []int64

	SingleValue   VarInt
	ContainerType ContainerType
}

func (p *PalettedContainer) WriteTo(w io.Writer) (int64, error) {
	// single valued for 0 bits per entry
	if p.BitsPerEntry == 0 {
		n, err := Serialize(w, struct {
			BitsPerEntry uint8
			Value        VarInt
		}{
			BitsPerEntry: 0,
			Value:        p.SingleValue,
		})
		return n, err
	}

	total := int64(0)
	n, err := w.Write([]byte{p.BitsPerEntry})
	total += int64(n)
	if err != nil {
		return total, err
	}

	// indirect for 4-8 bits per block state or 1-3 bits per biome
	if (4 <= p.BitsPerEntry && p.BitsPerEntry <= 8 && p.ContainerType == BlockStates) || (1 <= p.BitsPerEntry && p.BitsPerEntry <= 3 && p.ContainerType == Biomes) {
		n, err := p.Palette.WriteTo(w)
		total += n
		if err != nil {
			return total, err
		}
	}

	// direct for everything else

	buffer := make([]byte, len(p.Storage)*8)
	for i, v := range p.Storage {
		binary.BigEndian.PutUint64(buffer[i*8:], uint64(v))
	}
	n, err = w.Write(buffer)
	total += int64(n)

	return total, err
}

type BlockEntity struct {
	PackedXZ uint8
	Y        int16
	Type     VarInt
	Data     NBTValue
}

type LightData struct {
	SkyLightMask        BitSet
	BlockLightMask      BitSet
	EmptySkyLightMask   BitSet
	EmptyBlockLightMask BitSet
	SkyLightArray       PrefixedArray[LightArray]
	BlockLightArray     PrefixedArray[LightArray]
}

type PacketPlayOutLevelChunkWithLight struct {
	ChunkX        int32
	ChunkZ        int32
	Heightmaps    PrefixedArray[Heightmap]
	Data          PrefixedArray[uint8]
	BlockEntities PrefixedArray[BlockEntity]
	LightData     LightData
}

func (p *PacketPlayOutLevelChunkWithLight) ID() string {
	return "level_chunk_with_light"
}

type PacketPlayOutPlayerPosition struct {
	TeleportID VarInt
	X          float64
	Y          float64
	Z          float64
	VelocityX  float64
	VelocityY  float64
	VelocityZ  float64
	Yaw        float32
	Pitch      float32
	Flags      int32
}

func (p *PacketPlayOutPlayerPosition) ID() string {
	return "player_position"
}

type PacketPlayOutDisconnect struct {
	Reason NBTValue // text component encoded as NBT
}

func (p *PacketPlayOutDisconnect) ID() string {
	return "disconnect"
}
