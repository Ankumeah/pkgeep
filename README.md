# pkgeep

pkgeep is a package helper that tracks all your packages across
all your package managers with multiple storage backends

## How it works
You define all the package managers and storage options you want pkgeep to
track and sync to in you config file. pkgeep uses hooks provided
by most package managers to automatically maintaian list of all installed
packaes and write them to whatever storage you want

Later you can download all you packages with any one of your storage options

## Out of the box supported package manger
- [pacman](https://pacman.archlinux.page/) (AUR packages installed through pacman (most AUR helpers like [yay](https://jguer.github.io/yay/) install packages through pacman) are also included)
- [apt](https://wiki.debian.org/AptCLI)
- [dpkg](https://www.dpkg.org/) (It is reccamended to use apt instead, download mode is unsupported)
- [dnf5](https://dnf5.readthedocs.io/en/latest/index.html) (post-transaction-actions plugin is required)

If your package manager is not listed here, keep reading

## Out of the box supported storage backends
- local json file
- [Github Gist](https://gist.github.com/)
- [S3](https://aws.amazon.com/s3) (Any S3 compatable object store)

If your storage option is not listed here, keep reading

## What if my package manager or storage option is not in the supported list?
> [!NOTE]
> This feature is still in dev and not released yet

pkgeep allows custom entries to make sure you are supported on
whatever bespoke package manger or storage option you use (I dont judge)

## Config
pkgeep by default uses a yaml based config file which is
loaded from `~/.config/pkgeep/config.yaml`, it can be specified
otherwise, by specifying you can also use a json config.
Specs for the aboved mentioned can be found at `./specs/main.json`

## Command usage
```text
Usage: ./build/out <mode> [<flags>]

	Track all your package across all your package managers

Modes:
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

	download: Dwonloads all packages from the top source
		Flags:

	help: Prints the help text
		Flags:
```

## Contributing
- Follow commit prefix conventions (tho fun commit bodies are appriciated ;) )
- AI assisted is allowed, slop written isen't
