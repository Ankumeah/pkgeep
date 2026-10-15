package packagemanagers

type Package struct {
	Name    string `json:"name"`
	Manager string `json:"manager"`
}

type PackageManager interface {
	Name() string

	Setup(bin string, config map[string]any) error
	Teardown(config map[string]any) error
	GetAllPackages(config map[string]any) ([]Package, error)
	Add(pkg string, config map[string]any) Package
}
