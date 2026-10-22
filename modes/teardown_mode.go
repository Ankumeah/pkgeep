package modes

import (
	"fmt"
	"os"

	"github.com/Ankumeah/pkgeep/package-managers"
	"github.com/Ankumeah/pkgeep/utils"
)

type TeardownMode struct{}

func init() { Modes[TeardownMode{}.Name()] = TeardownMode{} }

func (m TeardownMode) Name() string        { return "teardown" }
func (m TeardownMode) Description() string { return "Teardowns all setup" }
func (m TeardownMode) Flags() []flag       { return []flag{configFileFlag} }

func (m TeardownMode) Handle(args []string) {
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

		fmt.Printf("Tearing down %v hook\n", manager.Name)

		if err := m.Teardown(manager.Config); err != nil {
			fmt.Printf("Error while tearing down %v hook: %v\n", manager.Name, err)
			os.Exit(1)
		}
	}
}
