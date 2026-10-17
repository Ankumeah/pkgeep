package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

func LoadConfig(path string) (packageManagers []ConfigEntry, storages []ConfigEntry) {
	if path == "" {
		path = filepath.Join(ConfigDir(), configFile)
	}

	configType := "yaml"
	if strings.HasSuffix(path, ".json") {
		configType = "json"
	}

	configFile, err := os.Open(path)
	if err != nil {
		fmt.Printf("Error while opening config file: %v\n", err)
		os.Exit(1)
	}
	defer configFile.Close()

	configEntries := []ConfigEntry{}

	switch configType {
	case "yaml":
		err = yaml.NewDecoder(configFile).Decode(&configEntries)
	case "json":
		err = json.NewDecoder(configFile).Decode(&configEntries)
	}

	if err != nil {
		fmt.Printf("Error in config: %v\n", err)
		os.Exit(1)
	}

	for _, entry := range configEntries {
		switch entry.Type {
		case "package_manager":
			packageManagers = append(packageManagers, entry)
		case "storage":
			storages = append(storages, entry)
		default:
			fmt.Printf("Invalid config entry %v. Skipping\n", entry.Type)
		}
	}

	return
}
