package storage

import (
	"github.com/Ankumeah/pkgeep/package-managers"
)

type Storage interface {
	Name() string

	Save(managers []packagemanagers.PackageManager) error
	Get() ([]packagemanagers.PackageManager, error)
}
