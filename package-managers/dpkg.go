package packagemanagers

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/Ankumeah/pkgeep/utils"
)

const dpkgHookDir = "/etc/dpkg/dpkg.cfg.d/"
const dpkgTempHookFile = dpkgHookDir + ".pkgeep.tmp"
const dpkgHookFile = dpkgHookDir + "pkgeep"

var dpkgHook = `
cat << 'EOF' > %s
post-invoke=%s sync --config-file=%s
EOF
`

type Dpkg struct{}

func init() { PackageManagers[Dpkg{}.Name()] = Dpkg{} }

func (p Dpkg) Name() string { return "dpkg" }

func (p Dpkg) Setup(config map[string]any) error {
	ctx, cancel := utils.GetTimeoutContext(config)
	defer cancel()

	su := utils.GetFromConfig[string]("su", config)

  runCmd(ctx, su, "mkdir", "-p", dpkgHookDir)

	// Not using /tmp to avoid cross fs errors
	// when renaming while ensuring atomicity
	runCmd(
		ctx,
		su, "sh", "-c",
		fmt.Sprintf(
			dpkgHook, dpkgTempHookFile,
			utils.GetFromConfig[string]("cmd", config),
			utils.GetFromConfig[string]("config_file", config),
		),
	)

	runCmd(ctx, su, "mv", dpkgTempHookFile, dpkgHookFile)

	return nil
}

func (p Dpkg) Teardown(config map[string]any) error {
	ctx, cancel := utils.GetTimeoutContext(config)
	defer cancel()

	if _, err := os.Stat(dpkgHookFile); os.IsNotExist(err) {
		return nil
	} else if err != nil {
		return err
	}

	runCmd(ctx, utils.GetFromConfig[string]("su", config), "rm", dpkgHookFile)

	return nil
}

func (p Dpkg) GetAllPackages(config map[string]any) ([]Package, error) {
	ctx, cancel := utils.GetTimeoutContext(config)
	defer cancel()

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

func (m Dpkg) Download(packages []Package, config map[string]any) error {
	return errors.New("Cannot download with dpkg")
}
