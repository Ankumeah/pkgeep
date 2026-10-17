package packagemanagers

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/Ankumeah/pkgeep/utils"
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
	ctx := context.Background()

	timeout := utils.GetFromConfig[float64]("timeout_seconds", config)
	if timeout < 0 {
		var cancel func()
		ctx, cancel = context.WithTimeout(ctx, time.Second*time.Duration(timeout))
		defer cancel()
	}

	su := utils.GetFromConfig[string]("su", config)

	// Not using /tmp to avoid cross fs errors
	// when renaming while ensuring atomicity
	pacmanRunCmd(
		ctx,
		su, "sh", "-c",
		fmt.Sprintf(
			pacmanHook, pacmanHookFile+".tmp",
			utils.GetFromConfig[string]("cmd", config),
			utils.GetFromConfig[string]("config_file", config),
		),
	)

	pacmanRunCmd(ctx, su, "mv", pacmanHookFile+".tmp", pacmanHookFile)

	return nil
}

func (p Pacman) Teardown(config map[string]any) error {
	ctx := context.Background()

	timeout := utils.GetFromConfig[float64]("timeout_seconds", config)
	if timeout < 0 {
		var cancel func()
		ctx, cancel = context.WithTimeout(ctx, time.Second*time.Duration(timeout))
		defer cancel()
	}
	if _, err := os.Stat(pacmanHookFile); os.IsNotExist(err) {
		return nil
	} else if err != nil {
		return err
	}

	pacmanRunCmd(ctx, utils.GetFromConfig[string]("su", config), "rm", pacmanHookFile)

	return nil
}

func (p Pacman) GetAllPackages(config map[string]any) ([]Package, error) {
	ctx := context.Background()

	timeout := utils.GetFromConfig[float64]("timeout_seconds", config)
	if timeout < 0 {
		var cancel func()
		ctx, cancel = context.WithTimeout(ctx, time.Second*time.Duration(timeout))
		defer cancel()
	}

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

func pacmanRunCmd(ctx context.Context, name string, args ...string) {
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Printf("Error while running root command: %v\n", err)
		os.Exit(1)
	}
}
