package mojang

import (
	"encoding/json"
	"net/http"
)

const manifestURL = "https://launchermeta.mojang.com/mc/game/version_manifest_v2.json"

type VersionManifest struct {
	Versions []struct {
		ID  string `json:"id"`
		URL string `json:"url"`
	} `json:"versions"`
}

type VersionMeta struct {
	Downloads struct {
		Server struct {
			URL  string `json:"url"`
			SHA1 string `json:"sha1"`
		} `json:"server"`
	} `json:"downloads"`
}

func fetchVersionManifest(manifest *VersionManifest, tr *http.Transport) error {
	client := &http.Client{Transport: tr}
	resp, err := client.Get(manifestURL)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(manifest); err != nil {
		return err
	}

	return nil
}

func fetchVersionMetadata(version string, metadata *VersionMeta, tr *http.Transport) error {
	var versionManifest VersionManifest
	if err := fetchVersionManifest(&versionManifest, tr); err != nil {
		return err
	}

	var versionMetadataURL string

	for _, versionmf := range versionManifest.Versions {
		if versionmf.ID == version {
			versionMetadataURL = versionmf.URL
		}
	}

	client := &http.Client{Transport: tr}
	resp, err := client.Get(versionMetadataURL)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(metadata); err != nil {
		return err
	}

	return nil
}
