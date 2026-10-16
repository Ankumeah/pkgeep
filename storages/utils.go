package storages

import (
	"fmt"
	"os"
)

func getSavePath(config map[string]any) string {
	_savePath, ok := config["save_path"]
	if !ok {
		fmt.Println("save_path not provided in config")
		os.Exit(1)
	}

	savePath, ok := _savePath.(string)
	if !ok {
		fmt.Println("save_path in config must be string")
		os.Exit(1)
	}

	if savePath == "" {
		fmt.Println("save_path in config cannot be empty")
		os.Exit(1)
	}

	return savePath
}
