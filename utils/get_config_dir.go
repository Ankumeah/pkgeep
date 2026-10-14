package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

func ConfigDir() string {
	baseConfigDir, err := os.UserConfigDir()
	if err != nil {
		fmt.Printf(
			"Error while getting config dir "+
				"(what kind of potato are you even "+
				"running this on?): %v\n",
			err,
		)
	}

	return filepath.Join(baseConfigDir, "pkgeep")
}
