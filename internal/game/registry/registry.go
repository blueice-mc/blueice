package registry

import (
	"github.com/blueice-mc/blueice/internal/game/defs"
	"github.com/blueice-mc/blueice/internal/network/protocol"
)

type TagEntry struct {
	Name protocol.Identifier
	IDs  []protocol.VarInt
}

type Registry[T any] struct {
	Entries map[protocol.Identifier]T
	IDs     map[protocol.Identifier]protocol.VarInt
	Tags    []TagEntry
}

type Registries struct {
	CatSoundVariant       Registry[defs.CatSoundVariant]
	CatVariant            Registry[defs.CatVariant]
	ChickenSoundVariant   Registry[defs.ChickenSoundVariant]
	ChickenVariant        Registry[defs.ChickenVariant]
	CowSoundVariant       Registry[defs.CowSoundVariant]
	CowVariant            Registry[defs.CowVariant]
	PigSoundVariant       Registry[defs.PigSoundVariant]
	PigVariant            Registry[defs.PigVariant]
	WolfSoundVariant      Registry[defs.WolfSoundVariant]
	WolfVariant           Registry[defs.WolfVariant]
	FrogVariant           Registry[defs.FrogVariant]
	PaintingVariant       Registry[defs.PaintingVariant]
	ZombieNautilusVariant Registry[defs.ZombieNautilusVariant]
	DamageType            Registry[defs.DamageType]
	TrimMaterial          Registry[defs.TrimMaterial]
	JukeboxSong           Registry[defs.JukeboxSong]
	BannerPattern         Registry[defs.BannerPattern]
	Instrument            Registry[defs.Instrument]
}
