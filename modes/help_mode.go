package modes

import (
	"fmt"
	"strings"
)

type HelpMode struct{}

func init() { Modes[HelpMode{}.Name()] = HelpMode{} }

func (m HelpMode) Name() string        { return "help" }
func (m HelpMode) Description() string { return "Prints the help text" }
func (m HelpMode) Flags() []flag       { return []flag{} }

func (m HelpMode) Handle(args []string) {
	var text strings.Builder

	text.WriteString("Usage: " + args[0] + " <mode> [<flags>]\n")
	text.WriteString("\n\tTrack all your package across all your package managers\n\n")

	text.WriteString("Modes:\n")
	for _, mode := range Modes {
		text.WriteString("\t" + mode.Name() + ": " + mode.Description() + "\n")

		text.WriteString("\t\tFlags:\n")
		for _, flag := range mode.Flags() {
			text.WriteString("\t\t\t" + flag.name + ", " + flag.shortName + ": " + flag.description + "\n")
			text.WriteString("\t\t\t\tUsage: " + flag.usage + "\n")
		}

		text.WriteByte('\n')
	}

	fmt.Print(text.String())
}
