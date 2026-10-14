package main

import (
	"fmt"
	"os"

	"github.com/Ankumeah/pkgeep/modes"
)

func main() {
	args := os.Args

	if len(args) < 1 {
		fmt.Println("Nice potato pc! (argc <= 1)")
		os.Exit(1)
	} else if len(args) < 2 {
		fmt.Println("No mode provided")
		modes.HelpMode{}.Handle(args)
		os.Exit(1)
	}

	if mode, ok := modes.Modes[args[1]]; !ok {
		fmt.Printf("Undefined mode: %v\n", args[1])
		os.Exit(1)
	} else {
		mode.Handle(args)
	}
}
