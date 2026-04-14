package mojang

import (
	"BlueIce/internal/version"
	"archive/zip"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func unzipInnerJar(path string, version string) error {
	parent := filepath.Dir(path)

	r, err := zip.OpenReader(path)
	if err != nil {
		return err
	}

	defer r.Close()

	for _, f := range r.File {
		if strings.Contains(f.Name, "server-"+version+".jar") {
			rc, err := f.Open()
			if err != nil {
				return err
			}

			copyToPath := filepath.Join(parent, filepath.Base(f.Name))

			log.Printf("Extracting %s", f.Name)

			file, err := os.OpenFile(copyToPath, os.O_CREATE|os.O_RDWR, 0600)
			if err != nil {
				return err
			}

			_, err = io.Copy(file, rc)
			if err != nil {
				if err = file.Close(); err != nil {
					return err
				}
				return err
			}

			if err = file.Close(); err != nil {
				return err
			}
		}
	}

	return nil
}

func unzipJar(path string) error {
	parent := filepath.Dir(path)

	r, err := zip.OpenReader(path)
	if err != nil {
		return err
	}

	defer r.Close()

	for _, f := range r.File {
		if strings.HasPrefix(f.Name, "data/minecraft") && !f.FileInfo().IsDir() {
			rc, err := f.Open()
			if err != nil {
				return err
			}

			copyToPath := filepath.Join(parent, f.Name)

			log.Printf("Inflating %s", copyToPath)

			if err = os.MkdirAll(filepath.Dir(copyToPath), 0755); err != nil {
				return err
			}

			file, err := os.OpenFile(copyToPath, os.O_CREATE|os.O_RDWR, 0600)
			if err != nil {
				return err
			}

			_, err = io.Copy(file, rc)
			if err != nil {
				if err = file.Close(); err != nil {
					return err
				}
				return err
			}

			if err = file.Close(); err != nil {
				return err
			}
		}
	}

	return nil
}

func FetchMinecraftData(path string) error {
	if err := DownloadServerJAR(version.GameVersion, path+"/server.jar"); err != nil {
		return err
	}

	if err := unzipInnerJar(path+"/server.jar", version.GameVersion); err != nil {
		return err
	}
	
	if err := unzipJar(path + "/server-" + version.GameVersion + ".jar"); err != nil {
		return err
	}

	if err := os.Remove(path + "/server.jar"); err != nil {
		return err
	}

	if err := os.Remove(path + "/server-" + version.GameVersion + ".jar"); err != nil {
		return err
	}

	return nil
}
