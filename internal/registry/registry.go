package registry

import (
	"BlueIce/internal/defs"
)

type Registries struct {
	CatSoundVariant       map[string]defs.CatSoundVariant
	CatVariant            map[string]defs.CatVariant
	ChickenSoundVariant   map[string]defs.ChickenSoundVariant
	ChickenVariant        map[string]defs.ChickenVariant
	CowSoundVariant       map[string]defs.CowSoundVariant
	CowVariant            map[string]defs.CowVariant
	PigSoundVariant       map[string]defs.PigSoundVariant
	PigVariant            map[string]defs.PigVariant
	WolfSoundVariant      map[string]defs.WolfSoundVariant
	WolfVariant           map[string]defs.WolfVariant
	FrogVariant           map[string]defs.FrogVariant
	PaintingVariant       map[string]defs.PaintingVariant
	ZombieNautilusVariant map[string]defs.ZombieNautilusVariant
	DamageType            map[string]defs.DamageType
}
