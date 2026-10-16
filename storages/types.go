package storages

import (
	"github.com/Ankumeah/pkgeep/package-managers"
)

type Storage interface {
	Name() string

	Save(managers []packagemanagers.Package, config map[string]any) error
	Get(config map[string]any) ([]packagemanagers.Package, error)
}
