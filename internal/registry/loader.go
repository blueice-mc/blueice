package registry

import (
	"BlueIce/internal/defs"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

func Load[T any](dataDir string, registry string) (map[string]T, error) {
	entries, err := os.ReadDir(filepath.Join(dataDir, "/lib/data", registry))
	if err != nil {
		return nil, err
	}

	entryMap := make(map[string]T)

	for _, entry := range entries {
		if !strings.HasSuffix(entry.Name(), ".json") {
			continue
		}

		entryName := strings.TrimSuffix(entry.Name(), ".json")

		data, err := os.ReadFile(filepath.Join(dataDir, "/lib/data", registry, entry.Name()))
		if err != nil {
			return nil, err
		}

		var currentEntry T

		if err := json.Unmarshal(data, &currentEntry); err != nil {
			return nil, err
		}

		entryMap[entryName] = currentEntry
	}

	return entryMap, nil
}

func (r *Registries) LoadAll(dataDir string) error {
	catSoundVariants, err := Load[defs.CatSoundVariant](dataDir, "minecraft/cat_sound_variant")
	if err != nil {
		return err
	}
	r.CatSoundVariant = catSoundVariants

	catVariants, err := Load[defs.CatVariant](dataDir, "minecraft/cat_variant")
	if err != nil {
		return err
	}
	r.CatVariant = catVariants

	chickenSoundVariants, err := Load[defs.ChickenSoundVariant](dataDir, "minecraft/chicken_sound_variant")
	if err != nil {
		return err
	}
	r.ChickenSoundVariant = chickenSoundVariants

	chickenVariants, err := Load[defs.ChickenVariant](dataDir, "minecraft/chicken_variant")
	if err != nil {
		return err
	}
	r.ChickenVariant = chickenVariants

	cowSoundVariants, err := Load[defs.CowSoundVariant](dataDir, "minecraft/cow_sound_variant")
	if err != nil {
		return err
	}
	r.CowSoundVariant = cowSoundVariants

	cowVariants, err := Load[defs.CowVariant](dataDir, "minecraft/cow_variant")
	if err != nil {
		return err
	}
	r.CowVariant = cowVariants

	pigSoundVariants, err := Load[defs.PigSoundVariant](dataDir, "minecraft/pig_sound_variant")
	if err != nil {
		return err
	}
	r.PigSoundVariant = pigSoundVariants

	pigVariants, err := Load[defs.PigVariant](dataDir, "minecraft/pig_variant")
	if err != nil {
		return err
	}
	r.PigVariant = pigVariants

	wolfSoundVariants, err := Load[defs.WolfSoundVariant](dataDir, "minecraft/wolf_sound_variant")
	if err != nil {
		return err
	}
	r.WolfSoundVariant = wolfSoundVariants

	wolfVariants, err := Load[defs.WolfVariant](dataDir, "minecraft/wolf_variant")
	if err != nil {
		return err
	}
	r.WolfVariant = wolfVariants

	frogVariant, err := Load[defs.FrogVariant](dataDir, "minecraft/frog_variant")
	if err != nil {
		return err
	}
	r.FrogVariant = frogVariant

	paintingVariant, err := Load[defs.PaintingVariant](dataDir, "minecraft/painting_variant")
	if err != nil {
		return err
	}
	r.PaintingVariant = paintingVariant

	zombieNautilusVariant, err := Load[defs.ZombieNautilusVariant](dataDir, "minecraft/zombie_nautilus_variant")
	if err != nil {
		return err
	}
	r.ZombieNautilusVariant = zombieNautilusVariant

	damageType, err := Load[defs.DamageType](dataDir, "minecraft/damage_type")
	if err != nil {
		return err
	}
	r.DamageType = damageType

	return nil
}
