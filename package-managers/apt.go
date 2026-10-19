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

const aptTempHookFile = "/etc/apt/apt.conf.d/.pkgeep.tmp"
const aptHookFile = "/etc/apt/apt.conf.d/pkgeep"

var aptHook = `
cat << 'EOF' > %s
DPkg::Post-Invoke { "%s sync --config-file=%s"; };
EOF
`

type Apt struct{}

func init() { PackageManagers[Apt{}.Name()] = Apt{} }

func (p Apt) Name() string { return "apt" }

func (p Apt) Setup(config map[string]any) error {
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
	aptRunCmd(
		ctx,
		su, "sh", "-c",
		fmt.Sprintf(
			aptHook, aptTempHookFile,
			utils.GetFromConfig[string]("cmd", config),
			utils.GetFromConfig[string]("config_file", config),
		),
	)

	aptRunCmd(ctx, su, "mv", aptTempHookFile, aptHookFile)

	return nil
}

func (p Apt) Teardown(config map[string]any) error {
	ctx := context.Background()

	timeout := utils.GetFromConfig[float64]("timeout_seconds", config)
	if timeout < 0 {
		var cancel func()
		ctx, cancel = context.WithTimeout(ctx, time.Second*time.Duration(timeout))
		defer cancel()
	}
	if _, err := os.Stat(aptHookFile); os.IsNotExist(err) {
		return nil
	} else if err != nil {
		return err
	}

	aptRunCmd(ctx, utils.GetFromConfig[string]("su", config), "rm", aptHookFile)

	return nil
}

func (p Apt) GetAllPackages(config map[string]any) ([]Package, error) {
	ctx := context.Background()

	timeout := utils.GetFromConfig[float64]("timeout_seconds", config)
	if timeout < 0 {
		var cancel func()
		ctx, cancel = context.WithTimeout(ctx, time.Second*time.Duration(timeout))
		defer cancel()
	}

	buf := bytes.NewBuffer(nil)

	cmd := exec.CommandContext(ctx, "sh", "-c", `
    for pkg in $(apt-mark showmanual);
      do dpkg-query -W "$pkg";
    done
  `)
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
			Manager: "apt",
		})
	}

	return packages, nil
}

func aptRunCmd(ctx context.Context, name string, args ...string) {
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Printf("Error while running root command: %v\n", err)
		os.Exit(1)
	}
}
