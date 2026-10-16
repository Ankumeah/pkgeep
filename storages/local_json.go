package storages

import (
	"encoding/json"
	"os"

	"github.com/Ankumeah/pkgeep/package-managers"
)

type LocalJSONStorage struct{}

func init() { Stoarges[LocalJSONStorage{}.Name()] = LocalJSONStorage{} }

func (s LocalJSONStorage) Name() string { return "local_json" }

func (s LocalJSONStorage) Save(
	managers []packagemanagers.Package,
	config map[string]any,
) error {
	savePath := getSavePath(config)

	file, err := os.OpenFile(
		savePath+".tmp",
		os.O_EXCL|os.O_CREATE|os.O_RDWR,
		0o0644,
	)
	if err != nil {
		return err
	}
	defer func() {
		file.Close()
		os.Remove(file.Name())
	}()

	if err := json.NewEncoder(
		file,
	).Encode(managers); err != nil {
		return err
	}

	return os.Rename(file.Name(), savePath)
}

func (s LocalJSONStorage) Get(config map[string]any) ([]packagemanagers.Package, error) {
	savePath := getSavePath(config)

	file, err := os.Open(savePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	managers := []packagemanagers.Package{}
	if err := json.NewDecoder(
		file,
	).Decode(&managers); err != nil {
		return nil, err
	}

	return managers, nil
}
