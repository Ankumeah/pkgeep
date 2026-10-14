package packagemanagers

type Package struct {
	Name    string `json:"name"`
	Manager string `json:"manager"`
}

type PackageManager interface {
	Setup(bin string) error
	Teardown() error
	GetAllPackages() ([]Package, error)
	Add(pkg string) Package
}
