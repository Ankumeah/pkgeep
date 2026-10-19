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

const dpkgTempHookFile = "/etc/dpkg/dpkg.cfg.d/.pkgeep.tmp"
const dpkgHookFile = "/etc/dpkg/dpkg.cfg.d/pkgeep"

var dpkgHook = `
cat << 'EOF' > %s
post-invoke=%s sync --config-file=%s
EOF
`

type Dpkg struct{}

func init() { PackageManagers[Dpkg{}.Name()] = Dpkg{} }

func (p Dpkg) Name() string { return "dpkg" }

func (p Dpkg) Setup(config map[string]any) error {
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
	dpkgRunCmd(
		ctx,
		su, "sh", "-c",
		fmt.Sprintf(
			dpkgHook, dpkgTempHookFile,
			utils.GetFromConfig[string]("cmd", config),
			utils.GetFromConfig[string]("config_file", config),
		),
	)

	dpkgRunCmd(ctx, su, "mv", dpkgTempHookFile, dpkgHookFile)

	return nil
}

func (p Dpkg) Teardown(config map[string]any) error {
	ctx := context.Background()

	timeout := utils.GetFromConfig[float64]("timeout_seconds", config)
	if timeout < 0 {
		var cancel func()
		ctx, cancel = context.WithTimeout(ctx, time.Second*time.Duration(timeout))
		defer cancel()
	}
	if _, err := os.Stat(dpkgHookFile); os.IsNotExist(err) {
		return nil
	} else if err != nil {
		return err
	}

	dpkgRunCmd(ctx, utils.GetFromConfig[string]("su", config), "rm", dpkgHookFile)

	return nil
}

func (p Dpkg) GetAllPackages(config map[string]any) ([]Package, error) {
	ctx := context.Background()

	timeout := utils.GetFromConfig[float64]("timeout_seconds", config)
	if timeout < 0 {
		var cancel func()
		ctx, cancel = context.WithTimeout(ctx, time.Second*time.Duration(timeout))
		defer cancel()
	}

	buf := bytes.NewBuffer(nil)

	cmd := exec.CommandContext(ctx, "dpkg-query", "-W")
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
			Manager: "dpkg",
		})
	}

	return packages, nil
}

func dpkgRunCmd(ctx context.Context, name string, args ...string) {
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Printf("Error while running root command: %v\n", err)
		os.Exit(1)
	}
}
