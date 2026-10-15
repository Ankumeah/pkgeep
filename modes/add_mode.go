package modes

type AddMode struct{}

func init() { Modes[AddMode{}.Name()] = AddMode{} }

func (m AddMode) Name() string        { return "add" }
func (m AddMode) Description() string { return "Adds a package to sources" }
func (m AddMode) Flags() []flag       { return []flag{} }

func (m AddMode) Handle(args []string) {
	panic("Uniplimented")
}
