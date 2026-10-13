package modes

type Mode interface {
	Name() string
	Description() string
	Flags() []flag

	Handle(args []string)
}

type flag interface {
	name() string
	shortName() string
	description() string
}
