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

	if cmd == "" {
		fmt.Println("cmd in config cannot be empty")
		os.Exit(1)
	}

	return cmd
}

func getConfigFile(config map[string]any) string {
	_configFile, ok := config["config_file"]
	if !ok {
		fmt.Println("config_file not provided in config")
		os.Exit(1)
	}

	config_file, ok := _configFile.(string)
	if !ok {
		fmt.Println("config_file in config must be string")
		os.Exit(1)
	}

	if config_file == "" {
		fmt.Println("config_file in config must not be empty")
		os.Exit(1)
	}

	return config_file
}
