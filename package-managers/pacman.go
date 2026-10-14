package packagemanagers

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

const pacmanHookFile = "/usr/share/libalpm/hooks/pkgeep.hook"
const pacmanTempFilePattern = "pkgeek-pacman-*"

const pacmanHook = `
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
When = PostTransaction\n
`

type Pacman struct{}

func (p Pacman) Setup(bin string) error {
	tmp, err := os.CreateTemp("", pacmanTempFilePattern)
	if err != nil {
		return err
	}
	defer func() {
		os.Remove(tmp.Name())
		tmp.Close()
	}()

	hook := fmt.Sprintf(pacmanHook, bin)

	if _, err := tmp.WriteString(hook); err != nil {
		return err
	}

	if err := os.Rename(tmp.Name(), pacmanHookFile); err != nil {
		return err
	}

	return nil
}

func (p Pacman) Teardown() error {
	return os.Remove(pacmanHookFile)
}

func (p Pacman) GetAllPackages() ([]Package, error) {
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

func (p Pacman) Add(pkg string) Package {
	return Package{
		Name:    pkg,
		Manager: "pacman",
	}
}
