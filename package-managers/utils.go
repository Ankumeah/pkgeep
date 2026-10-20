package packagemanagers

import (
	"context"
	"fmt"
	"os"
	"os/exec"
)

func runCmd(ctx context.Context, name string, args ...string) {
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Printf("Error while running root command: %v\n", err)
		os.Exit(1)
	}
}
