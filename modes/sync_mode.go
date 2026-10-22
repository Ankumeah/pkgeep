package modes

import (
	"fmt"
	"os"

	"github.com/Ankumeah/pkgeep/package-managers"
	"github.com/Ankumeah/pkgeep/storages"
	"github.com/Ankumeah/pkgeep/utils"
)

type SyncMode struct{}

func init() { Modes[SyncMode{}.Name()] = SyncMode{} }

func (m SyncMode) Name() string        { return "sync" }
func (m SyncMode) Description() string { return "Syncs all packages to sources" }
func (m SyncMode) Flags() []flag       { return []flag{configFileFlag} }

func (m SyncMode) Handle(args []string) {
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

	allPackages := []packagemanagers.Package{}
	for _, managerConfig := range managerConfigs {
		if !managerConfig.Enabled {
			continue
		}

		manager, ok := packagemanagers.PackageManagers[managerConfig.Name]
		if !ok {
			fmt.Printf("Unsupported package manager %v\n", managerConfig.Name)
			os.Exit(1)
		}

		fmt.Printf("Syncing %v\n", managerConfig.Name)

		packages, err := manager.GetAllPackages(managerConfig.Config)
		if err != nil {
			fmt.Printf("Error while getting packages: %v\n", err)
			os.Exit(1)
		}
		allPackages = append(allPackages, packages...)
	}

	for _, storageConfig := range storageConfigs {
		if !storageConfig.Enabled {
			continue
		}

		storage, ok := storages.Stoarges[storageConfig.Name]
		if !ok {
			fmt.Printf("Unsupported storage %v\n", storageConfig.Name)
			os.Exit(1)
		}

		fmt.Printf("Syncing to %v\n", storageConfig.Name)

		if err := storage.Save(allPackages, storageConfig.Config); err != nil {
			fmt.Printf("Error while syncing to source: %v\n", err)
			fmt.Println("Not all sources have been synced. It is reccamended to sync again")
			os.Exit(1)
		}
	}
}
