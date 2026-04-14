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

type SpawnCondition struct {
	Priority  int32               `nbt:"priority"            json:"priority"`
	Condition *SpawnConditionData `nbt:"condition,omitempty" json:"condition,omitempty"`
}

type SpawnConditionData struct {
	Type string `nbt:"type"                  json:"type"`
	// minecraft:structure
	Structures string `nbt:"structures,omitempty"  json:"structures,omitempty"`
	// minecraft:biome
	Biomes string `nbt:"biomes,omitempty"      json:"biomes,omitempty"`
	// minecraft:moon_brightness
	Range *FloatRange `nbt:"range,omitempty"       json:"range,omitempty"`
}

type FloatRange struct {
	Min *float32 `nbt:"min,omitempty" json:"min,omitempty"`
	Max *float32 `nbt:"max,omitempty" json:"max,omitempty"`
}

type CatVariant struct {
	AssetID         string           `nbt:"asset_id"          json:"asset_id"`
	BabyAssetID     string           `nbt:"baby_asset_id"     json:"baby_asset_id"`
	SpawnConditions []SpawnCondition `nbt:"spawn_conditions"  json:"spawn_conditions"`
}

type CatSoundSet struct {
	AmbientSound      string `nbt:"ambient_sound,omitempty"      json:"ambient_sound,omitempty"`
	BegForFoodSound   string `nbt:"beg_for_food_sound,omitempty" json:"beg_for_food_sound,omitempty"`
	DeathSound        string `nbt:"death_sound,omitempty"        json:"death_sound,omitempty"`
	EatSound          string `nbt:"eat_sound,omitempty"          json:"eat_sound,omitempty"`
	HissSound         string `nbt:"hiss_sound,omitempty"         json:"hiss_sound,omitempty"`
	HurtSound         string `nbt:"hurt_sound,omitempty"         json:"hurt_sound,omitempty"`
	PurrSound         string `nbt:"purr_sound,omitempty"         json:"purr_sound,omitempty"`
	PurreowSound      string `nbt:"purreow_sound,omitempty"      json:"purreow_sound,omitempty"`
	StrayAmbientSound string `nbt:"stray_ambient_sound,omitempty" json:"stray_ambient_sound,omitempty"`
}

type CatSoundVariant struct {
	AdultSounds CatSoundSet `nbt:"adult_sounds" json:"adult_sounds"`
	BabySounds  CatSoundSet `nbt:"baby_sounds"  json:"baby_sounds"`
}

type ChickenVariant struct {
	AssetID         string           `nbt:"asset_id"          json:"asset_id"`
	BabyAssetID     string           `nbt:"baby_asset_id"     json:"baby_asset_id"`
	Model           string           `nbt:"model,omitempty"   json:"model,omitempty"`
	SpawnConditions []SpawnCondition `nbt:"spawn_conditions"  json:"spawn_conditions"`
}

type ChickenSoundSet = PigSoundSet

type ChickenSoundVariant struct {
	AdultSounds ChickenSoundSet `nbt:"adult_sounds" json:"adult_sounds"`
	BabySounds  ChickenSoundSet `nbt:"baby_sounds"  json:"baby_sounds"`
}

type CowVariant struct {
	AssetID         string           `nbt:"asset_id"          json:"asset_id"`
	BabyAssetID     string           `nbt:"baby_asset_id"     json:"baby_asset_id"`
	Model           string           `nbt:"model,omitempty"   json:"model,omitempty"`
	SpawnConditions []SpawnCondition `nbt:"spawn_conditions"  json:"spawn_conditions"`
}

type CowSoundVariant struct {
	AmbientSound string `nbt:"ambient_sound,omitempty" json:"ambient_sound,omitempty"`
	DeathSound   string `nbt:"death_sound,omitempty"   json:"death_sound,omitempty"`
	HurtSound    string `nbt:"hurt_sound,omitempty"    json:"hurt_sound,omitempty"`
	StepSound    string `nbt:"step_sound,omitempty"    json:"step_sound,omitempty"`
}

type FrogVariant struct {
	AssetID         string           `nbt:"asset_id"         json:"asset_id"`
	SpawnConditions []SpawnCondition `nbt:"spawn_conditions" json:"spawn_conditions"`
}

type PigVariant struct {
	AssetID         string           `nbt:"asset_id"          json:"asset_id"`
	BabyAssetID     string           `nbt:"baby_asset_id"     json:"baby_asset_id"`
	Model           string           `nbt:"model,omitempty"   json:"model,omitempty"`
	SpawnConditions []SpawnCondition `nbt:"spawn_conditions"  json:"spawn_conditions"`
}

type PigSoundSet struct {
	AmbientSound string `nbt:"ambient_sound,omitempty" json:"ambient_sound,omitempty"`
	DeathSound   string `nbt:"death_sound,omitempty"   json:"death_sound,omitempty"`
	EatSound     string `nbt:"eat_sound,omitempty"     json:"eat_sound,omitempty"`
	HurtSound    string `nbt:"hurt_sound,omitempty"    json:"hurt_sound,omitempty"`
	StepSound    string `nbt:"step_sound,omitempty"    json:"step_sound,omitempty"`
}

type PigSoundVariant struct {
	AdultSounds PigSoundSet `nbt:"adult_sounds" json:"adult_sounds"`
	BabySounds  PigSoundSet `nbt:"baby_sounds"  json:"baby_sounds"`
}

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
	AmbientSound string `nbt:"ambient_sound,omitempty" json:"ambient_sound,omitempty"`
	DeathSound   string `nbt:"death_sound,omitempty"   json:"death_sound,omitempty"`
	GrowlSound   string `nbt:"growl_sound,omitempty"   json:"growl_sound,omitempty"`
	HurtSound    string `nbt:"hurt_sound,omitempty"    json:"hurt_sound,omitempty"`
	PantSound    string `nbt:"pant_sound,omitempty"    json:"pant_sound,omitempty"`
	StepSound    string `nbt:"step_sound,omitempty"    json:"step_sound,omitempty"`
	WhineSound   string `nbt:"whine_sound,omitempty"   json:"whine_sound,omitempty"`
}

type WolfSoundVariant struct {
	AdultSounds WolfSoundSet `nbt:"adult_sounds" json:"adult_sounds"`
	BabySounds  WolfSoundSet `nbt:"baby_sounds"  json:"baby_sounds"`
}

type ZombieNautilusVariant struct {
	AssetID         string           `nbt:"asset_id"          json:"asset_id"`
	Model           string           `nbt:"model,omitempty"   json:"model,omitempty"`
	SpawnConditions []SpawnCondition `nbt:"spawn_conditions"  json:"spawn_conditions"`
}

type PaintingVariant struct {
	AssetID string `nbt:"asset_id" json:"asset_id"`
	Width   int32  `nbt:"width"    json:"width"`
	Height  int32  `nbt:"height"   json:"height"`
}

type DamageType struct {
	Effects    string  `nbt:"effects,omitempty" json:"effects,omitempty"`
	Exhaustion float32 `nbt:"exhaustion"        json:"exhaustion"`
	MessageID  string  `nbt:"message_id"        json:"message_id"`
	Scaling    string  `nbt:"scaling"           json:"scaling"`
}
