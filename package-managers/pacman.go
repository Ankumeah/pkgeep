package packagemanagers

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/Ankumeah/pkgeep/utils"
)

const pacmanHookDir = "/usr/share/libalpm/hooks/"
const pacmanTempHookFile = pacmanHookDir + "pkgeep.hook"
const pacmanHookFile = pacmanHookDir + "pkgeep.hook"

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
	ctx, cancel := utils.GetTimeoutContext(config)
	defer cancel()

	su := utils.GetFromConfig[string]("su", config)

	runCmd(ctx, su, "mkdir", "-p", pacmanHookDir)

	// Not using /tmp to avoid cross fs errors
	// when renaming while ensuring atomicity
	runCmd(
		ctx,
		su, "sh", "-c",
		fmt.Sprintf(
			pacmanHook, pacmanTempHookFile,
			utils.GetFromConfig[string]("cmd", config),
			utils.GetFromConfig[string]("config_file", config),
		),
	)

	runCmd(ctx, su, "mv", pacmanTempHookFile, pacmanHookFile)

	return nil
}

func (p Pacman) Teardown(config map[string]any) error {
	ctx, cancel := utils.GetTimeoutContext(config)
	defer cancel()

	if _, err := os.Stat(pacmanHookFile); os.IsNotExist(err) {
		return nil
	} else if err != nil {
		return err
	}

	runCmd(ctx, utils.GetFromConfig[string]("su", config), "rm", pacmanHookFile)

	return nil
}

func (p Pacman) GetAllPackages(config map[string]any) ([]Package, error) {
	ctx, cancel := utils.GetTimeoutContext(config)
	defer cancel()

	buf := bytes.NewBuffer(nil)

	cmd := exec.CommandContext(ctx, "pacman", "-Qe")
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
			Version: strings.Fields(pkg)[1],
			Manager: "pacman",
		})
	}

	return packages, nil
}

func (m Pacman) Download(packages []Package, config map[string]any) error {
	ctx, cancel := utils.GetTimeoutContext(config)
	defer cancel()

	args := []string{"-Syu"}
	for _, pkg := range packages {
		args = append(args, pkg.Name)
	}

	if utils.GetFromConfig[bool]("aur", config) {
		runCmd(
			ctx,
			utils.GetFromConfig[string]("aur_helper", config),
			args...,
		)
	} else {
		args = append([]string{"pacman"}, args...)
		runCmd(
			ctx,
			utils.GetFromConfig[string]("su", config),
			args...,
		)
	}

	return nil
}
