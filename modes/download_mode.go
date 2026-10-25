package modes

import (
	"fmt"
	"os"

	"github.com/Ankumeah/pkgeep/package-managers"
	"github.com/Ankumeah/pkgeep/storages"
	"github.com/Ankumeah/pkgeep/utils"
)

type DownloadMode struct{}

func init() { Modes[DownloadMode{}.Name()] = DownloadMode{} }

func (m DownloadMode) Name() string        { return "download" }
func (m DownloadMode) Description() string { return "Downloads all packages from the top source" }
func (m DownloadMode) Flags() []flag       { return []flag{configFileFlag} }

func (m DownloadMode) Handle(args []string) {
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

	allPackages := map[string][]packagemanagers.Package{}

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

		for _, pkg := range packages {
			allPackages[pkg.Manager] = append(allPackages[pkg.Manager], pkg)
		}
	}

	for _, managerConfig := range managerConfigs {
		if !managerConfig.Enabled {
			continue
		}

		m, ok := packagemanagers.PackageManagers[managerConfig.Name]
		if !ok {
			fmt.Printf("Unsupported package manager %v\n", managerConfig.Name)
			os.Exit(1)
		}

		fmt.Printf("Downloading for %v\n", managerConfig.Name)

		if err := m.Download(allPackages[managerConfig.Name], managerConfig.Config); err != nil {
			fmt.Printf("Error while downloading for %v: %v\n", managerConfig.Name, err)
			os.Exit(1)
		}
	}
}
