package mojang

import (
	"bytes"
	"fmt"
	"os"
	"strings"
)

func GenerateEulaFile(path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	lines := []string{
		"#You need to agree to the Minecraft EULA in order to run the server.",
		"#By changing the setting below to true, you are indicating your agreement to the minecraft EULA (https://aka.ms/MinecraftEULA).",
		"eula=false",
	}

	for _, line := range lines {
		if _, err := file.WriteString(line + "\n"); err != nil {
			return err
		}
	}

	return err
}

func EulaAccepted(path string) (bool, error) {
	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, GenerateEulaFile(path)
		}

		return false, err
	}
	defer file.Close()

	var buf bytes.Buffer

	if _, err := buf.ReadFrom(file); err != nil {
		return false, err
	}

	eula := false

	for line := range strings.Lines(buf.String()) {
		if strings.HasPrefix(line, "#") {
			continue
		}

		contents := strings.Split(line, "=")
		if len(contents) != 2 {
			return false, fmt.Errorf("invalid EULA file: %s", path)
		}

		if strings.TrimSpace(contents[0]) == "eula" {
			eula = strings.TrimSpace(contents[1]) == "true"
			if !eula {
				return false, nil
			}
			eula = true
		}
	}

	return eula, nil
}
