package protocol

type PacketPlayOutLogin struct {
	EntityID            VarInt
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
	GameMode            uint8
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
