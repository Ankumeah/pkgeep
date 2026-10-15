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
- None as of now (still in dev)

If your package manager is not listed here, keep reading

## Out of the box supported storage backends
- None as of now (still in dev)

If your storage option is not listed here, keep reading

## What if my package manager or storage option is not in the supported list?
pkgeep allows custom entries to make sure you are supported on
whatever bespoke package manger or storage option you use (I dont judge)

## Command usage
In dev


## Example config

> [!NOTE]
> Full spec + YAML configs coming soon!

```json
[
  {
    "name": "pacman",
    "type": "package_manager",
    "enabled": true,
    "config": {
      "su": "sudo"
    }
  },
  {
    "name": "local_json",
    "type": "storage",
    "enabled": true,
    "config": {}
  }
]
```

## Contributing
- Follow commit prefix conventions (tho fun commit bodies are appriciated ;) )
- AI assisted is allowed, slop written isen't
