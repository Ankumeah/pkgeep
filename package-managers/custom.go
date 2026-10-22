package packagemanagers

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/Ankumeah/pkgeep/utils"
)

type Custom struct{}

func init() { PackageManagers[Custom{}.Name()] = Custom{} }

func (p Custom) Name() string { return "custom" }

func (p Custom) Setup(config map[string]any) error {
	ctx, cancel := utils.GetTimeoutContext(config)
	defer cancel()

	cmd := ""
	args := []string{}

	cmdArrAny := utils.GetFromConfig[[]any]("setup_cmd", config)
	cmdArr := []string{}
	for _, item := range cmdArrAny {
		itemString, ok := item.(string)
		if !ok {
			fmt.Printf("setup_cmd must be []string")
			os.Exit(1)
		}

		cmdArr = append(cmdArr, itemString)
	}

	if len(cmdArr) > 1 {
		cmd = cmdArr[0]
		args = cmdArr[1:]
	} else if len(cmdArr) > 0 {
		cmd = cmdArr[0]
	}

	runCmd(ctx, cmd, args...)

	return nil
}

func (p Custom) Teardown(config map[string]any) error {
	ctx, cancel := utils.GetTimeoutContext(config)
	defer cancel()

	cmd := ""
	args := []string{}

	cmdArrAny := utils.GetFromConfig[[]any]("teardown_cmd", config)
	cmdArr := []string{}
	for _, item := range cmdArrAny {
		itemString, ok := item.(string)
		if !ok {
			fmt.Printf("teardown_cmd must be []string")
			os.Exit(1)
		}

		cmdArr = append(cmdArr, itemString)
	}

	if len(cmdArr) > 1 {
		cmd = cmdArr[0]
		args = cmdArr[1:]
	} else if len(cmdArr) > 0 {
		cmd = cmdArr[0]
	}

	runCmd(ctx, cmd, args...)

	return nil
}

func (p Custom) GetAllPackages(config map[string]any) ([]Package, error) {
	ctx, cancel := utils.GetTimeoutContext(config)
	defer cancel()

	buf := bytes.NewBuffer(nil)

	c := ""
	args := []string{}

	cmdArrAny := utils.GetFromConfig[[]any]("get_all_packages_cmd", config)
	cmdArr := []string{}
	for _, item := range cmdArrAny {
		itemString, ok := item.(string)
		if !ok {
			fmt.Printf("get_all_packages_cmd must be []string")
			os.Exit(1)
		}

		cmdArr = append(cmdArr, itemString)
	}

	if len(cmdArr) > 1 {
		c = cmdArr[0]
		args = cmdArr[1:]
	} else if len(cmdArr) > 0 {
		c = cmdArr[0]
	}

	cmd := exec.CommandContext(ctx, c, args...)
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
			Name:    utils.GetFromConfig[string]("name", config) + "-" + strings.Fields(pkg)[0],
			Version: strings.Fields(pkg)[1],
			Manager: "custom",
		})
	}

	return packages, nil
}

func (m Custom) Download(packages []Package, config map[string]any) error {
	ctx, cancel := utils.GetTimeoutContext(config)
	defer cancel()

	name := utils.GetFromConfig[string]("name", config)

	cmd := ""
	args := []string{}

	cmdArrAny := utils.GetFromConfig[[]any]("download_cmd", config)
	cmdArr := []string{}
	for _, item := range cmdArrAny {
		itemString, ok := item.(string)
		if !ok {
			fmt.Printf("download_cmd must be []string")
			os.Exit(1)
		}

		cmdArr = append(cmdArr, itemString)
	}

	if len(cmdArr) > 1 {
		cmd = cmdArr[0]
		args = cmdArr[1:]
	} else if len(cmdArr) > 0 {
		cmd = cmdArr[0]
	}

	for _, pkg := range packages {
		if p, ok := strings.CutPrefix(pkg.Name, name+"-"); ok {
			args = append(args, p)
		}
	}

	runCmd(ctx, cmd, args...)

	return nil
}
