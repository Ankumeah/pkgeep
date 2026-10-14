package storage

import (
	"github.com/Ankumeah/pkgeep/package-managers"
)

type PackageManager struct {
	Name     string                    `json:"name"`
	Enabled  bool                      `json:"enabled"`
	Packages []packagemanagers.Package `json:"packages"`
	Config   map[string]any            `json:"config"`
}

type Storage interface {
	Save(managers []PackageManager) error
	Get() ([]PackageManager, error)
}
