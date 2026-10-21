package packagemanagers

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/Ankumeah/pkgeep/utils"
)

const dnf5HookDir = "/etc/dnf/libdnf5-plugins/actions.d/"
const dnf5TempHookFile = dnf5HookDir + ".pkgeep.tmp"
const dnf5HookFile = dnf5HookDir + "pkgeep.actions"

var dnf5Hook = `
cat << 'EOF' > %s
post_transaction::::%s sync --config-file=%s
EOF
`

type Dnf5 struct{}

func init() { PackageManagers[Dnf5{}.Name()] = Dnf5{} }

func (p Dnf5) Name() string { return "dnf5" }

func (p Dnf5) Setup(config map[string]any) error {
	ctx, cancel := utils.GetTimeoutContext(config)
	defer cancel()

	su := utils.GetFromConfig[string]("su", config)

  runCmd(ctx, su, "mkdir", "-p", dnf5HookDir)

	// Not using /tmp to avoid cross fs errors
	// when renaming while ensuring atomicity
	runCmd(
		ctx,
		su, "sh", "-c",
		fmt.Sprintf(
			dnf5Hook, dnf5TempHookFile,
			utils.GetFromConfig[string]("cmd", config),
			utils.GetFromConfig[string]("config_file", config),
		),
	)

	runCmd(ctx, su, "mv", dnf5TempHookFile, dnf5HookFile)

	return nil
}

func (p Dnf5) Teardown(config map[string]any) error {
	ctx, cancel := utils.GetTimeoutContext(config)
	defer cancel()

	if _, err := os.Stat(dnf5HookFile); os.IsNotExist(err) {
		return nil
	} else if err != nil {
		return err
	}

	runCmd(ctx, utils.GetFromConfig[string]("su", config), "rm", dnf5HookFile)

	return nil
}

func (p Dnf5) GetAllPackages(config map[string]any) ([]Package, error) {
	ctx, cancel := utils.GetTimeoutContext(config)
	defer cancel()

	buf := bytes.NewBuffer(nil)

  // We cant repoquery while parent dnf keeps its lock
  // so we sync with the --skip-file-locks flag
	cmd := exec.CommandContext(ctx, "dnf5", "--skip-file-locks", "repoquery", "--userinstalled", "--queryformat", "%{name} %{version}\\n")
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
			Manager: "dnf5",
		})
	}

	return packages, nil
}

func (m Dnf5) Download(packages []Package, config map[string]any) error {
	ctx, cancel := utils.GetTimeoutContext(config)
	defer cancel()

	args := []string{"dnf5", "install"}
	for _, pkg := range packages {
		args = append(args, pkg.Name)
	}

  runCmd(ctx, utils.GetFromConfig[string]("su", config), args...)

  return nil
}
