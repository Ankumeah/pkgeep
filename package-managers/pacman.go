package packagemanagers

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

const pacmanHookFile = "/usr/share/libalpm/hooks/pkgeep.hook"

var pacmanHook = `
cat << 'EOF' > %s
[Trigger]
Operation = Install
Operation = Upgrade
Operation = Remove
Type = Package
Target = *

[Action]
Description = Tracking packages for pkgeek
Exec = %s sync --config-file=%s
NeedsTargets
When = PostTransaction
EOF
`

type Pacman struct{}

func init() { PackageManagers[Pacman{}.Name()] = Pacman{} }

func (p Pacman) Name() string { return "pacman" }

func (p Pacman) Setup(config map[string]any) error {
	su := getSu(config)

	// Not using /tmp to avoid cross fs errors
	// when renaming while ensuring atomicity
	cmd := exec.Command(
		su, "sh", "-c",
		fmt.Sprintf(
			pacmanHook, pacmanHookFile+".tmp",
			getCmd(config), getConfigFile(config),
		),
	)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout

	if err := cmd.Run(); err != nil {
		fmt.Printf("Error while writing pacman hook: %v\n", err)
		os.Exit(1)
	}
	defer os.Remove(pacmanHookFile + ".tmp")

	cmd = exec.Command(su, "mv", pacmanHookFile+".tmp", pacmanHookFile)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Printf("Error while writing pacman hook: %v\n", err)
		os.Exit(1)
	}

	return nil
}

func (p Pacman) Teardown(config map[string]any) error {
	if _, err := os.Stat(pacmanHookFile); os.IsNotExist(err) {
		return nil
	} else if err != nil {
		return err
	}

	cmd := exec.Command(getSu(config), "rm", pacmanHookFile)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func (p Pacman) GetAllPackages(config map[string]any) ([]Package, error) {
	buf := bytes.NewBuffer(nil)

	cmd := exec.Command("pacman", "-Qe")
	cmd.Stdout = buf

	if err := cmd.Run(); err != nil {
		return nil, err
	}

	packages := []Package{}

	for pkg := range strings.SplitSeq(buf.String(), "\n") {
		if pkg == "" {
			continue
		}

		packages = append(packages, Package{
			Name:    strings.Fields(pkg)[0],
			Manager: "pacman",
		})
	}

	return packages, nil
}
