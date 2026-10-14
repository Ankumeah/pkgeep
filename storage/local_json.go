package storage

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/Ankumeah/pkgeep/utils"
)

const localJSONStorageFile = "local_json_storage.json"
const localJSONStorageTempFilePattern = "pkgeep-local_json_storage-*"

type LocalJSONStorage struct{}

func (s LocalJSONStorage) Save(
	managers []PackageManager,
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

func Get() ([]PackageManager, error) {
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

	managers := []PackageManager{}
	if err := json.NewDecoder(
		file,
	).Decode(&managers); err != nil {
		return nil, err
	}

	return managers, nil
}
