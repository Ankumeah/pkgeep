package modes

import (
	"fmt"
	"os"

	packagemanagers "github.com/Ankumeah/pkgeep/package-managers"
	"github.com/Ankumeah/pkgeep/utils"
)

type SetupMode struct{}

func init() { Modes[SetupMode{}.Name()] = SetupMode{} }

func (m SetupMode) Name() string        { return "setup" }
func (m SetupMode) Description() string { return "Sets up all hooks and pings sources" }
func (m SetupMode) Flags() []flag       { return []flag{configFileFlag} }

func (m SetupMode) Handle(args []string) {
	configFile := ""

	for _, arg := range args[2:] {
		if value, isConfigFlag := isFlag(arg, configFileFlag); isConfigFlag {
			configFile = value
		} else {
			fmt.Printf("Unknown flag %v. Aborting\n", arg)
			os.Exit(0)
		}
	}

	managers, _ := utils.LoadConfig(configFile)

	for _, manager := range managers {
		if !manager.Enabled {
			continue
		}

		m, ok := packagemanagers.PackageManagers[manager.Name]
		if !ok {
			fmt.Printf("Unsupported package manager %v\n", manager.Name)
			os.Exit(1)
		}

		fmt.Printf("Setting up %v hook\n", manager.Name)

		if err := m.Setup(manager.Config); err != nil {
			fmt.Printf("Error while setting up %v hook: %v\n", manager.Name, err)
			os.Exit(1)
		}
	}
}
