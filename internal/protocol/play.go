package protocol

import (
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

func (p *PacketPlayOutLogin) ID() VarInt {
	return 0x31
}

func (p *PacketPlayOutLogin) WriteTo(w io.Writer) (int64, error) {
	total, err := serialize(w, struct {
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
		n, err := serialize(w, p.DeathDimensionName)
		total += n
		if err != nil {
			return total, err
		}
		n, err = serialize(w, p.DeathLocation)
		total += n
		if err != nil {
			return total, err
		}
	}

	n, err := serialize(w, struct {
		PortalCooldown     VarInt
		SeaLevel           VarInt
		EnforcesSecureChat bool
	}{p.PortalCooldown, p.SeaLevel, p.EnforcesSecureChat})
	total += n
	return total, err
}
