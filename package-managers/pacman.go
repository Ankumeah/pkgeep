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
Exec = %s add
NeedsTargets
When = PostTransaction
EOF
`

type Pacman struct{}

func init() { PackageManagers[Pacman{}.Name()] = Pacman{} }

func (p Pacman) Name() string { return "pacman" }

func (p Pacman) Setup(bin string, config map[string]any) error {
	su := pacmanGetSu(config)

	// Not using /tmp to avoid cross fs errors
	// when renaming while ensuring atomicity
	cmd := exec.Command(
		su, "sh", "-c",
		fmt.Sprintf(pacmanHook, pacmanHookFile+".tmp", bin),
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

	cmd := exec.Command(
		pacmanGetSu(config),
		"rm", pacmanHookFile,
	)
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
		packages = append(packages, Package{
			Name:    strings.Fields(pkg)[0],
			Manager: "pacman",
		})
	}

	return packages, nil
}

func (p Pacman) Add(pkg string, config map[string]any) Package {
	return Package{
		Name:    pkg,
		Manager: "pacman",
	}
}

func pacmanGetSu(config map[string]any) string {
	_su, ok := config["su"]
	if !ok {
		fmt.Println("pacman setup needs root access, specify su within your config")
		os.Exit(1)
	}

	su, ok := _su.(string)
	if !ok {
		fmt.Println("su for pacman must of type string")
		os.Exit(1)
	}

	if su == "" {
		fmt.Println("su for pacman cannot be empty. Use a wrapper script for su if you really want to")
		os.Exit(1)
	}

	return su
}
