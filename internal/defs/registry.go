package defs

type IntProvider struct {
	Type  string `nbt:"type"  json:"type"`
	Value int32  `nbt:"value" json:"value"`
}

type DimensionTypeEntry struct {
	CoordinateScale             float64       `nbt:"coordinate_scale"               json:"coordinate_scale"`
	HasSkylight                 int8          `nbt:"has_skylight"                   json:"has_skylight"`
	HasCeiling                  int8          `nbt:"has_ceiling"                    json:"has_ceiling"`
	HasEnderDragonFight         int8          `nbt:"has_ender_dragon_fight"         json:"has_ender_dragon_fight"`
	AmbientLight                float32       `nbt:"ambient_light"                  json:"ambient_light"`
	HasFixedTime                int8          `nbt:"has_fixed_time"                 json:"has_fixed_time"`
	FixedTime                   *int64        `nbt:"fixed_time,omitempty"           json:"fixed_time,omitempty"`
	MonsterSpawnBlockLightLimit int32         `nbt:"monster_spawn_block_light_limit" json:"monster_spawn_block_light_limit"`
	MonsterSpawnLightLevel      IntProvider   `nbt:"monster_spawn_light_level"      json:"monster_spawn_light_level"`
	LogicalHeight               int32         `nbt:"logical_height"                 json:"logical_height"`
	MinY                        int32         `nbt:"min_y"                          json:"min_y"`
	Height                      int32         `nbt:"height"                         json:"height"`
	Infiniburn                  string        `nbt:"infiniburn"                     json:"infiniburn"`
	Skybox                      string        `nbt:"skybox"                         json:"skybox"`
	CardinalLight               string        `nbt:"cardinal_light"                 json:"cardinal_light"`
	Attributes                  EmptyCompound `nbt:"attributes"                   json:"attributes"`
	DefaultClock                string        `nbt:"default_clock"                  json:"default_clock"`
	Timelines                   []string      `nbt:"timelines"                      json:"timelines"`
}

type WorldClockEntry struct{}
type EmptyCompound struct{}

type ChatType struct {
	Chat      ChatFormat `nbt:"chat"      json:"chat"`
	Narration ChatFormat `nbt:"narration" json:"narration"`
}

type ChatFormat struct {
	TranslationKey string   `nbt:"translation_key" json:"translation_key"`
	Parameters     []string `nbt:"parameters"      json:"parameters"`
}

type Biome struct {
	HasPrecipitation    int8         `nbt:"has_precipitation"    json:"has_precipitation"`
	Temperature         float32      `nbt:"temperature"          json:"temperature"`
	Downfall            float32      `nbt:"downfall"             json:"downfall"`
	Effects             BiomeEffects `nbt:"effects"              json:"effects"`
	TemperatureModifier string       `nbt:"temperature_modifier" json:"temperature_modifier"`
}

type BiomeEffects struct {
	FogColor           int32  `nbt:"fog_color"            json:"fog_color"`
	SkyColor           int32  `nbt:"sky_color"            json:"sky_color"`
	WaterColor         int32  `nbt:"water_color"          json:"water_color"`
	WaterFogColor      int32  `nbt:"water_fog_color"      json:"water_fog_color"`
	GrassColorModifier string `nbt:"grass_color_modifier" json:"grass_color_modifier"`
}

type SpawnCondition struct{}

type CatVariant struct {
	AssetID         string           `nbt:"asset_id"          json:"asset_id"`
	BabyAssetID     string           `nbt:"baby_asset_id"     json:"baby_asset_id"`
	SpawnConditions []SpawnCondition `nbt:"spawn_conditions"  json:"spawn_conditions"`
}

type CatSoundVariant = WolfSoundVariant // gleiche Struktur

type ChickenVariant struct {
	AssetID         string           `nbt:"asset_id"         json:"asset_id"`
	BabyAssetID     string           `nbt:"baby_asset_id"    json:"baby_asset_id"`
	Model           string           `nbt:"model"            json:"model"`
	SpawnConditions []SpawnCondition `nbt:"spawn_conditions" json:"spawn_conditions"`
}

type ChickenSoundVariant = WolfSoundVariant

type CowVariant struct {
	AssetID         string           `nbt:"asset_id"         json:"asset_id"`
	BabyAssetID     string           `nbt:"baby_asset_id"    json:"baby_asset_id"`
	Model           string           `nbt:"model"            json:"model"`
	SpawnConditions []SpawnCondition `nbt:"spawn_conditions" json:"spawn_conditions"`
}

type CowSoundVariant = WolfSoundVariant

type FrogVariant struct {
	AssetID         string           `nbt:"asset_id"         json:"asset_id"`
	SpawnConditions []SpawnCondition `nbt:"spawn_conditions" json:"spawn_conditions"`
}

type PigVariant struct {
	AssetID         string           `nbt:"asset_id"         json:"asset_id"`
	BabyAssetID     string           `nbt:"baby_asset_id"    json:"baby_asset_id"`
	Model           string           `nbt:"model"            json:"model"`
	SpawnConditions []SpawnCondition `nbt:"spawn_conditions" json:"spawn_conditions"`
}

type PigSoundVariant = WolfSoundVariant

type WolfAssets struct {
	Angry string `nbt:"angry" json:"angry"`
	Wild  string `nbt:"wild"  json:"wild"`
	Tame  string `nbt:"tame"  json:"tame"`
}

type WolfVariant struct {
	Assets          WolfAssets       `nbt:"assets"           json:"assets"`
	BabyAssets      WolfAssets       `nbt:"baby_assets"      json:"baby_assets"`
	SpawnConditions []SpawnCondition `nbt:"spawn_conditions" json:"spawn_conditions"`
}

type WolfSoundSet struct {
	AmbientSound string `nbt:"ambient_sound" json:"ambient_sound"`
	DeathSound   string `nbt:"death_sound"   json:"death_sound"`
	GrowlSound   string `nbt:"growl_sound"   json:"growl_sound"`
	HurtSound    string `nbt:"hurt_sound"    json:"hurt_sound"`
	PantSound    string `nbt:"pant_sound"    json:"pant_sound"`
	StepSound    string `nbt:"step_sound"    json:"step_sound"`
	WhineSound   string `nbt:"whine_sound"   json:"whine_sound"`
}

type WolfSoundVariant struct {
	AdultSounds WolfSoundSet `nbt:"adult_sounds" json:"adult_sounds"`
	BabySounds  WolfSoundSet `nbt:"baby_sounds"  json:"baby_sounds"`
}

type ZombieNautilusVariant struct {
	AssetID         string           `nbt:"asset_id"         json:"asset_id"`
	Model           string           `nbt:"model"            json:"model"`
	SpawnConditions []SpawnCondition `nbt:"spawn_conditions" json:"spawn_conditions"`
}

type PaintingVariant struct {
	AssetID string `nbt:"asset_id" json:"asset_id"`
	Width   int32  `nbt:"width"    json:"width"`
	Height  int32  `nbt:"height"   json:"height"`
}
