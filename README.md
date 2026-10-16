# pkgeep

pkgeep is a package helper that tracks all your packages across
all your package managers with multiple storage backends

## How it works
You define all the package managers and storage options you want pkgeep to
track and sync to in you config file. pkgeep uses hooks provided
by most package managers to automaticlly maintaina list of all installed
packaes and write them to whatever storage you want

Later you can download all you packages with any one of your storage options

## Out of the box supported package manger
- [pacman](https://pacman.archlinux.page/)

If your package manager is not listed here, keep reading

## Out of the box supported storage backends
- local json file

If your storage option is not listed here, keep reading

## What if my package manager or storage option is not in the supported list?
> [!NOTE]
> This feature is still in dev

pkgeep allows custom entries to make sure you are supported on
whatever bespoke package manger or storage option you use (I dont judge)

## Config
> [!NOTE]
> YAML based configs coming soon!

pkgeep uses a json based config file.
Specs for the aboved mentioned can be found in `./specs/`

## Command usage
```text
Usage: pgkeep <mode> [<flags>]

	Track all your package across all your package managers

Modes:
	help: Prints the help text
		Flags:

	setup: Sets up all hooks and pings sources
		Flags:
			--config-file, -c: Path to the config file
				Usage: (--config-file,-c)=path/to/config/file

	sync: Syncs all packages to sources
		Flags:
			--config-file, -c: Path to the config file
				Usage: (--config-file,-c)=path/to/config/file

	teardown: Teardowns all setup
		Flags:
			--config-file, -c: Path to the config file
				Usage: (--config-file,-c)=path/to/config/file
```

## Contributing
- Follow commit prefix conventions (tho fun commit bodies are appriciated ;) )
- AI assisted is allowed, slop written isen't
