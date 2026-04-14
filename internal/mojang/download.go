package mojang

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func validateSHA1Sum(path string, checksum []byte) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}

	defer file.Close()

	hash := sha1.New()
	if _, err := io.Copy(hash, file); err != nil {
		return err
	}

	calculatedChecksum := hash.Sum(nil)
	if !bytes.Equal(calculatedChecksum, checksum) {
		return errors.New("checksum mismatch")
	}

	return nil
}

func DownloadServerJAR(version string, path string) error {
	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
	}

	var meta VersionMeta
	if err := fetchVersionMetadata(version, &meta, tr); err != nil {
		return err
	}

	client := &http.Client{Transport: tr}

	log.Println("Starting download of mojang server JAR...")
	log.Printf("Downloading from: %s", meta.Downloads.Server.URL)

	resp, err := client.Get(meta.Downloads.Server.URL)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return err
	}

	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return err
	}

	err = file.Close()
	if err != nil {
		return err
	}

	log.Println("Finished download! Validating...")

	checksum, err := hex.DecodeString(meta.Downloads.Server.SHA1)
	if err != nil {
		return err
	}

	err = validateSHA1Sum(path, checksum)
	if err != nil {
		err2 := os.Remove(path)
		if err2 != nil {
			log.Fatalf("Could not validate SHA1 checksum and delete the corrupted Jarfile. "+
				"Please remove the following file manually: %s", path)
		}

		return err
	}

	log.Println("Validation complete!")

	return nil
}
