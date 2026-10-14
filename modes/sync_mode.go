package modes

type SyncMode struct {
	configPath string
	verbose    bool
}

func init() { Modes[SyncMode{}.Name()] = SyncMode{} }

func (m SyncMode) Name() string        { return "sync" }
func (m SyncMode) Description() string { return "Syncs all packages to sources" }
func (m SyncMode) Flags() []flag       { return []flag{} }

func (m SyncMode) Handle(args []string) {
	panic("Uniplimented")
}
