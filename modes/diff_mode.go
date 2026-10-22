package modes

import (
	"fmt"
	"os"
	"strconv"

	"github.com/Ankumeah/pkgeep/package-managers"
	"github.com/Ankumeah/pkgeep/storages"
	"github.com/Ankumeah/pkgeep/utils"
)

type DiffMode struct{}

func init() { Modes[DiffMode{}.Name()] = DiffMode{} }

func (m DiffMode) Name() string { return "diff" }
func (m DiffMode) Description() string {
	return "Diff of your installed packages and packages in the top source"
}
func (m DiffMode) Flags() []flag { return []flag{configFileFlag} }

func (m DiffMode) Handle(args []string) {
	configFile := ""

	for _, arg := range args[2:] {
		if value, isConfigFlag := isFlag(arg, configFileFlag); isConfigFlag {
			configFile = value
		} else {
			fmt.Printf("Unknown flag %v. Aborting\n", arg)
			os.Exit(0)
		}
	}

	managerConfigs, storageConfigs := utils.LoadConfig(configFile)

	storedPackages := map[string]map[packagemanagers.Package]struct{}{}

	for _, storageConfig := range storageConfigs {
		if !storageConfig.Enabled {
			continue
		}

		storage, ok := storages.Stoarges[storageConfig.Name]
		if !ok {
			fmt.Printf("Unsupported storage %v\n", storageConfig.Name)
			os.Exit(1)
		}

		packages, err := storage.Get(storageConfig.Config)
		if err != nil {
			fmt.Printf("Error while getting packages from %v: %v", storageConfig.Name, err)
			os.Exit(1)
		}

		if len(packages) > 0 {
			storedPackages[packages[0].Manager] = map[packagemanagers.Package]struct{}{}
		}
		for _, pkg := range packages {
			storedPackages[pkg.Manager][pkg] = struct{}{}
		}
	}

	notInstalled := []packagemanagers.Package{}
	notStored := []packagemanagers.Package{}

	maxNameLen := 0
	maxVersionLen := 0

	for _, managerConfig := range managerConfigs {
		if !managerConfig.Enabled {
			continue
		}

		m, ok := packagemanagers.PackageManagers[managerConfig.Name]
		if !ok {
			fmt.Printf("Unsupported package manager %v\n", managerConfig.Name)
			os.Exit(1)
		}

		fmt.Printf("Diffing for %v\n", managerConfig.Name)

		installedPackages, err := m.GetAllPackages(managerConfig.Config)
		if err != nil {
			fmt.Printf("Error while getting packages from %v: %v\n", managerConfig.Name, err)
			os.Exit(1)
		}

		for _, pkg := range installedPackages {
			maxNameLen = max(maxNameLen, len(pkg.Name))
			maxVersionLen = max(maxVersionLen, len(pkg.Version))

			if _, ok := storedPackages[pkg.Manager][pkg]; !ok {
				notStored = append(notStored, pkg)
			}
			delete(storedPackages[pkg.Manager], pkg)
		}
	}

	for _, manager := range storedPackages {
		for pkg, _ := range manager {
			maxNameLen = max(maxNameLen, len(pkg.Name))
			maxVersionLen = max(maxVersionLen, len(pkg.Version))

			notInstalled = append(notInstalled, pkg)
		}
	}

	fmt.Println("Not Installed:")
	maxDigits := len(strconv.Itoa(len(notInstalled)))
	for i, pkg := range notInstalled {
		fmt.Printf(
			"  %0*d. %-*s %-*s [%s]\n",
			maxDigits, i, maxNameLen, pkg.Name, maxVersionLen, pkg.Version, pkg.Manager,
		)
	}

	fmt.Println("Not Stored:")
	maxDigits = len(strconv.Itoa(len(notStored)))
	for i, pkg := range notStored {
		fmt.Printf(
			"  %0*d. %-*s %-*s [%s]\n",
			maxDigits, i, maxNameLen, pkg.Name, maxVersionLen, pkg.Version, pkg.Manager,
		)
	}
}
