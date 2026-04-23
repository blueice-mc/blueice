package registry

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/blueice-mc/blueice/internal/game/defs"
	"github.com/blueice-mc/blueice/internal/network/protocol"
)

func Load[T any](dataDir string, registryID protocol.Identifier, registry *Registry[T]) error {
	entries, err := os.ReadDir(filepath.Join(dataDir, "lib/data", registryID.Namespace, registryID.Path))
	if err != nil {
		return err
	}

	entryMap := make(map[protocol.Identifier]T)

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		if !strings.HasSuffix(entry.Name(), ".json") {
			continue
		}

		entryName := strings.TrimSuffix(entry.Name(), ".json")

		data, err := os.ReadFile(filepath.Join(dataDir, "lib/data", registryID.Namespace, registryID.Path, entry.Name()))
		if err != nil {
			return err
		}

		var currentEntry T

		if err := json.Unmarshal(data, &currentEntry); err != nil {
			return err
		}

		entryMap[protocol.NewIdentifier(registryID.Namespace, entryName)] = currentEntry
	}

	keys := make([]string, 0, len(entryMap))
	for k := range entryMap {
		keys = append(keys, k.Path)
	}
	sort.Strings(keys)

	idMap := make(map[protocol.Identifier]protocol.VarInt)
	for i, key := range keys {
		idMap[protocol.NewIdentifier(registryID.Namespace, key)] = protocol.VarInt(i)
	}

	registry.Entries = entryMap
	registry.IDs = idMap

	return nil
}

type TagFileContent struct {
	Values []string `json:"values"`
}

func resolveTag(name protocol.Identifier, collectedTags map[protocol.Identifier][]protocol.VarInt, collectedReferences map[protocol.Identifier][]protocol.Identifier, resolved map[protocol.Identifier]bool) {
	if resolved[name] {
		return
	}

	resolved[name] = true

	for _, reference := range collectedReferences[name] {
		resolveTag(reference, collectedTags, collectedReferences, resolved)
		collectedTags[name] = append(collectedTags[name], collectedTags[reference]...)
	}
}

func LoadTags(dataDir string, registryID protocol.Identifier, idMap map[protocol.Identifier]protocol.VarInt) (map[protocol.Identifier][]protocol.VarInt, error) {
	collectedTags := make(map[protocol.Identifier][]protocol.VarInt)
	collectedReferences := make(map[protocol.Identifier][]protocol.Identifier)

	err := filepath.WalkDir(filepath.Join(dataDir, "lib/data", registryID.Namespace, "tags", registryID.Path), func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() || !strings.HasSuffix(path, ".json") {
			return err
		}

		rel, _ := filepath.Rel(filepath.Join(dataDir, "lib/data", registryID.Namespace, "tags", registryID.Path), path)
		tagID := protocol.NewIdentifier(registryID.Namespace, strings.TrimSuffix(rel, ".json"))

		var content TagFileContent

		bytes, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		if err := json.Unmarshal(bytes, &content); err != nil {
			return err
		}

		for _, value := range content.Values {
			if strings.HasPrefix(value, "#") {
				collectedReferences[tagID] = append(collectedReferences[tagID], protocol.NewIdentifierFromString(strings.TrimPrefix(value, "#")))
			} else {
				id, ok := idMap[protocol.NewIdentifierFromString(value)]

				if !ok {
					return fmt.Errorf("invalid tag reference: %s", value)
				}

				collectedTags[tagID] = append(collectedTags[tagID], id)
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	resolved := make(map[protocol.Identifier]bool)
	for tagName := range collectedReferences {
		resolveTag(tagName, collectedTags, collectedReferences, resolved)
	}

	return collectedTags, nil
}

func (r *Registries) LoadAll(dataDir string) error {
	if err := Load[defs.Biome](dataDir, protocol.NewIdentifierFromPath("worldgen/biome"), &r.Biomes); err != nil {
		return err
	}

	if err := Load[defs.CatSoundVariant](dataDir, protocol.NewIdentifierFromPath("cat_sound_variant"), &r.CatSoundVariant); err != nil {
		return err
	}

	if err := Load[defs.CatVariant](dataDir, protocol.NewIdentifierFromPath("cat_variant"), &r.CatVariant); err != nil {
		return err
	}

	if err := Load[defs.ChickenSoundVariant](dataDir, protocol.NewIdentifierFromPath("chicken_sound_variant"), &r.ChickenSoundVariant); err != nil {
		return err
	}

	if err := Load[defs.ChickenVariant](dataDir, protocol.NewIdentifierFromPath("chicken_variant"), &r.ChickenVariant); err != nil {
		return err
	}

	if err := Load[defs.CowSoundVariant](dataDir, protocol.NewIdentifierFromPath("cow_sound_variant"), &r.CowSoundVariant); err != nil {
		return err
	}

	if err := Load[defs.CowVariant](dataDir, protocol.NewIdentifierFromPath("cow_variant"), &r.CowVariant); err != nil {
		return err
	}

	if err := Load[defs.PigSoundVariant](dataDir, protocol.NewIdentifierFromPath("pig_sound_variant"), &r.PigSoundVariant); err != nil {
		return err
	}

	if err := Load[defs.PigVariant](dataDir, protocol.NewIdentifierFromPath("pig_variant"), &r.PigVariant); err != nil {
		return err
	}

	if err := Load[defs.WolfSoundVariant](dataDir, protocol.NewIdentifierFromPath("wolf_sound_variant"), &r.WolfSoundVariant); err != nil {
		return err
	}

	if err := Load[defs.WolfVariant](dataDir, protocol.NewIdentifierFromPath("wolf_variant"), &r.WolfVariant); err != nil {
		return err
	}

	if err := Load[defs.FrogVariant](dataDir, protocol.NewIdentifierFromPath("frog_variant"), &r.FrogVariant); err != nil {
		return err
	}

	if err := Load[defs.PaintingVariant](dataDir, protocol.NewIdentifierFromPath("painting_variant"), &r.PaintingVariant); err != nil {
		return err
	}

	if err := Load[defs.ZombieNautilusVariant](dataDir, protocol.NewIdentifierFromPath("zombie_nautilus_variant"), &r.ZombieNautilusVariant); err != nil {
		return err
	}

	if err := Load[defs.DamageType](dataDir, protocol.NewIdentifierFromPath("damage_type"), &r.DamageType); err != nil {
		return err
	}

	tags, err := LoadTags(dataDir, protocol.NewIdentifierFromPath("damage_type"), r.DamageType.IDs)
	if err != nil {
		return err
	}

	for tag, ids := range tags {
		entry := TagEntry{
			Name: tag,
			IDs:  ids,
		}
		r.DamageType.Tags = append(r.DamageType.Tags, entry)
	}

	if err := Load[defs.TrimMaterial](dataDir, protocol.NewIdentifierFromPath("trim_material"), &r.TrimMaterial); err != nil {
		return err
	}

	if err := Load[defs.JukeboxSong](dataDir, protocol.NewIdentifierFromPath("jukebox_song"), &r.JukeboxSong); err != nil {
		return err
	}

	if err := Load[defs.BannerPattern](dataDir, protocol.NewIdentifierFromPath("banner_pattern"), &r.BannerPattern); err != nil {
		return err
	}

	tags, err = LoadTags(dataDir, protocol.NewIdentifierFromPath("banner_pattern"), r.BannerPattern.IDs)
	if err != nil {
		return err
	}

	for tag, ids := range tags {
		entry := TagEntry{
			Name: tag,
			IDs:  ids,
		}
		r.BannerPattern.Tags = append(r.BannerPattern.Tags, entry)
	}

	if err := Load[defs.Instrument](dataDir, protocol.NewIdentifierFromPath("instrument"), &r.Instrument); err != nil {
		return err
	}

	return nil
}
