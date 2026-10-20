package packagemanagers

type Package struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Manager string `json:"manager"`
}

type PackageManager interface {
	Name() string

	Setup(config map[string]any) error
	Teardown(config map[string]any) error
	GetAllPackages(config map[string]any) ([]Package, error)
	Download(packages []Package, config map[string]any) error
}
