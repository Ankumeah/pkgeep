package packagemanagers

import (
	"fmt"
	"os"
)

func getSu(config map[string]any) string {
	_su, ok := config["su"]
	if !ok {
		fmt.Println("su not provided in config")
		os.Exit(1)
	}

	su, ok := _su.(string)
	if !ok {
		fmt.Println("su in config must be string")
		os.Exit(1)
	}

	if su == "" {
		fmt.Println("su in config cannot be empty. Use a wrapper script for su if you really want to")
		os.Exit(1)
	}

	return su
}

func getCmd(config map[string]any) string {
  _cmd, ok := config["cmd"]
  if !ok {
    fmt.Println("cmd not provided in config")
    os.Exit(1)
  }

  cmd, ok := _cmd.(string)
  if !ok {
    fmt.Println("cmd in config must be string")
    os.Exit(1)
  }

  return cmd
}
