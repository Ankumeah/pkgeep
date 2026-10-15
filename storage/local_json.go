package storage

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/Ankumeah/pkgeep/package-managers"
	"github.com/Ankumeah/pkgeep/utils"
)

const localJSONStorageFile = "local_json_storage.json"
const localJSONStorageTempFilePattern = "pkgeep-local_json_storage-*"

type LocalJSONStorage struct{}

func init() { Stoarges[LocalJSONStorage{}.Name()] = LocalJSONStorage{} }

func (s LocalJSONStorage) Name() string { return "local_json" }

func (s LocalJSONStorage) Save(
	managers []packagemanagers.PackageManager,
) error {
	tmp, err := os.CreateTemp(
		"", localJSONStorageTempFilePattern,
	)
	if err != nil {
		return err
	}
	defer func() {
		os.Remove(tmp.Name())
		tmp.Close()
	}()

	if err := json.NewEncoder(
		tmp,
	).Encode(managers); err != nil {
		return err
	}

	return os.Rename(
		tmp.Name(),
		filepath.Join(
			utils.ConfigDir(),
			localJSONStorageFile,
		),
	)
}

func (s LocalJSONStorage) Get() ([]packagemanagers.PackageManager, error) {
	file, err := os.Open(
		filepath.Join(
			utils.ConfigDir(),
			localJSONStorageFile,
		),
	)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	managers := []packagemanagers.PackageManager{}
	if err := json.NewDecoder(
		file,
	).Decode(&managers); err != nil {
		return nil, err
	}

	return managers, nil
}
