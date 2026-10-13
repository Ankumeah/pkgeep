package modes

type HelpMode struct{}

func init() { Modes["help"] = HelpMode{} }

func (m HelpMode) Name() string        { return "help" }
func (m HelpMode) Description() string { return "Prints the help text" }
func (m HelpMode) Flags() []flag       { return []flag{} }

func (m HelpMode) Handle(args []string) {
	panic("Uniplimented")
}
