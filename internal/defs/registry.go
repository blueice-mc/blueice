package defs

type IntProvider struct {
	Type  string `nbt:"type"`
	Value int8   `nbt:"value"`
}

type DimensionTypeEntry struct {
	CoordinateScale     float64 `nbt:"coordinate_scale"`
	HasSkylight         int8    `nbt:"has_skylight"`
	HasCeiling          int8    `nbt:"has_ceiling"`
	HasEnderDragonFight int8    `nbt:"has_ender_dragon_fight"`
	AmbientLight        float32 `nbt:"ambient_light"`

	HasFixedTime int8   `nbt:"has_fixed_time"`
	FixedTime    *int64 `nbt:"fixed_time,omitempty"`

	MonsterSpawnBlockLightLimit int32       `nbt:"monster_spawn_block_light_limit"`
	MonsterSpawnLightLevel      IntProvider `nbt:"monster_spawn_light_level"`

	LogicalHeight int32 `nbt:"logical_height"`
	MinY          int32 `nbt:"min_y"`
	Height        int32 `nbt:"height"`

	Infiniburn string `nbt:"infiniburn"`

	Skybox        string `nbt:"skybox"`         // "overworld", "none", "end"
	CardinalLight string `nbt:"cardinal_light"` // "default", "nether"

	Attributes EmptyCompound `nbt:"attributes"`

	DefaultClock string `nbt:"default_clock"`

	Timelines []string `nbt:"timelines"`
}

type EmptyCompound struct{}

type ChatType struct {
	Chat      ChatFormat `nbt:"chat"`
	Narration ChatFormat `nbt:"narration"`
}

type ChatFormat struct {
	TranslationKey string   `nbt:"translation_key"`
	Parameters     []string `nbt:"parameters"`
}

type Biome struct {
	HasPrecipitation    int8         `nbt:"has_precipitation"`
	Temperature         float32      `nbt:"temperature"`
	Downfall            float32      `nbt:"downfall"`
	Effects             BiomeEffects `nbt:"effects"`
	TemperatureModifier string       `nbt:"temperature_modifier"`
}

type BiomeEffects struct {
	FogColor           int32  `nbt:"fog_color"`
	SkyColor           int32  `nbt:"sky_color"`
	WaterColor         int32  `nbt:"water_color"`
	WaterFogColor      int32  `nbt:"water_fog_color"`
	GrassColorModifier string `nbt:"grass_color_modifier"`
}
