package storages

import (
	"encoding/json"
	"os"

	"github.com/Ankumeah/pkgeep/package-managers"
	"github.com/Ankumeah/pkgeep/utils"
)

type LocalJSONStorage struct{}

func init() { Stoarges[LocalJSONStorage{}.Name()] = LocalJSONStorage{} }

func (s LocalJSONStorage) Name() string { return "local_json" }

func (s LocalJSONStorage) Save(
	packages []packagemanagers.Package,
	config map[string]any,
) error {
	savePath := utils.GetFromConfig[string]("save_path", config)

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
	).Encode(packages); err != nil {
		return err
	}

	return os.Rename(file.Name(), savePath)
}

func (s LocalJSONStorage) Get(
	config map[string]any,
) ([]packagemanagers.Package, error) {
	savePath := utils.GetFromConfig[string]("save_path", config)

	file, err := os.Open(savePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	packages := []packagemanagers.Package{}
	if err := json.NewDecoder(
		file,
	).Decode(&packages); err != nil {
		return nil, err
	}

	return packages, nil
}
