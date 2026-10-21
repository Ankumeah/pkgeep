package packagemanagers

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/Ankumeah/pkgeep/utils"
)

const aptHookDir = "/etc/apt/apt.conf.d/"
const aptTempHookFile = aptHookDir + ".pkgeep.tmp"
const aptHookFile = aptHookDir + "pkgeep"

var aptHook = `
cat << 'EOF' > %s
DPkg::Post-Invoke { "%s sync --config-file=%s"; };
EOF
`

type Apt struct{}

func init() { PackageManagers[Apt{}.Name()] = Apt{} }

func (p Apt) Name() string { return "apt" }

func (p Apt) Setup(config map[string]any) error {
	ctx, cancel := utils.GetTimeoutContext(config)
	defer cancel()

	su := utils.GetFromConfig[string]("su", config)

  runCmd(ctx, su, "mkdir", "-p", aptHookDir)

	// Not using /tmp to avoid cross fs errors
	// when renaming while ensuring atomicity
	runCmd(
		ctx,
		su, "sh", "-c",
		fmt.Sprintf(
			aptHook, aptTempHookFile,
			utils.GetFromConfig[string]("cmd", config),
			utils.GetFromConfig[string]("config_file", config),
		),
	)

	runCmd(ctx, su, "mv", aptTempHookFile, aptHookFile)

	return nil
}

func (p Apt) Teardown(config map[string]any) error {
	ctx, cancel := utils.GetTimeoutContext(config)
	defer cancel()

	if _, err := os.Stat(aptHookFile); os.IsNotExist(err) {
		return nil
	} else if err != nil {
		return err
	}

	runCmd(ctx, utils.GetFromConfig[string]("su", config), "rm", aptHookFile)

	return nil
}

func (p Apt) GetAllPackages(config map[string]any) ([]Package, error) {
	ctx, cancel := utils.GetTimeoutContext(config)
	defer cancel()

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

func (m Apt) Download(packages []Package, config map[string]any) error {
	ctx, cancel := utils.GetTimeoutContext(config)
	defer cancel()

	runCmd(ctx, utils.GetFromConfig[string]("su", config), "apt", "update")
	runCmd(ctx, utils.GetFromConfig[string]("su", config), "apt", "upgrade")

	args := []string{"apt", "install"}
	for _, pkg := range packages {
		args = append(args, pkg.Name)
	}

	runCmd(ctx, utils.GetFromConfig[string]("su", config), args...)

	return nil
}
