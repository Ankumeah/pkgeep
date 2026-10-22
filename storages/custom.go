package storages

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"

	"github.com/Ankumeah/pkgeep/package-managers"
	"github.com/Ankumeah/pkgeep/utils"
)

type CustomStorage struct{}

func init() { Stoarges[CustomStorage{}.Name()] = CustomStorage{} }

func (s CustomStorage) Name() string { return "custom" }

func (s CustomStorage) Save(
	packages []packagemanagers.Package,
	config map[string]any,
) error {
	ctx, cancel := utils.GetTimeoutContext(config)
	defer cancel()

	cmd := ""
	args := []string{}

	cmdArrAny := utils.GetFromConfig[[]any]("save_cmd", config)
	cmdArr := []string{}
	for _, item := range cmdArrAny {
		itemString, ok := item.(string)
		if !ok {
			fmt.Printf("save_cmd must be []string")
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

  data, err := json.Marshal(packages)
  if err != nil {
    return err
  }

  args = append(args, string(data))

  c := exec.CommandContext(ctx, cmd, args...)
  c.Stdout = os.Stdout
  c.Stderr = os.Stderr
  c.Stdin = os.Stdin

  return c.Run()
}

func (s CustomStorage) Get(
	config map[string]any,
) ([]packagemanagers.Package, error) {
	ctx, cancel := utils.GetTimeoutContext(config)
	defer cancel()

	cmd := ""
	args := []string{}

	cmdArrAny := utils.GetFromConfig[[]any]("get_cmd", config)
	cmdArr := []string{}
	for _, item := range cmdArrAny {
		itemString, ok := item.(string)
		if !ok {
			fmt.Printf("get_cmd must be []string")
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

  buf := bytes.NewBuffer(nil)

  c := exec.CommandContext(ctx, cmd, args...)
  c.Stdout = buf
  c.Stderr = os.Stderr

  if err := c.Run(); err != nil {
    return nil, err
  }

  packages := []packagemanagers.Package{}
  if err := json.NewDecoder(buf).Decode(&packages); err != nil {
    return nil, err
  }

  return packages, nil
}
