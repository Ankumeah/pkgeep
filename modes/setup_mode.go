package modes

type SetupMode struct {
	configPath string
	verbose    string
}

func init() { Modes[SetupMode{}.Name()] = SetupMode{} }

func (m SetupMode) Name() string        { return "setup" }
func (m SetupMode) Description() string { return "Sets up all hooks and pings sources" }
func (m SetupMode) Flags() []flag       { return []flag{} }

func (m SetupMode) Handle(args []string) {
	panic("Uniplimented")
}
