# `gen`

> **Work in Progress**: `gen` is not finished and has many rough
> edges.

An another static site generator.

```
~ $ gen
NAME:
   gen - An another static site generator.

USAGE:
   gen [global options] command [command options] [arguments...]

VERSION:
   0.2.4

COMMANDS:
   build, b   Performs a one off site build
   create, c  Creates a new site in the provided directory
   serve, s   Builds site and serves it locally
   remove, r  Removes all generated files
   help, h    Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --source DIR, -s DIR, --src DIR       read files from DIR (default: ".")
   --destination DIR, -d DIR, --dst DIR  write files to DIR (default: "site")
   --help, -h                            show help (default: false)
   --version, -v                         print the version (default: false)
~ $
```

## Getting Started

1. [Install](#installation) gen if you haven't yet.

2. Create a new site:

        $ gen create mysite

3. Change directory to `mysite` and serve the site locally:

        $ cd mysite
        $ gen serve

    Run with `--help` or `-h` for options.

4. Go to `http://localhost:3000`.

## Installation

### From binary

Download the precompiled binary from [releases page].

### From source

1. Install [Go] 1.14 if you haven't yet.

2. Two installation options are supported:

    * Install with `go get`:

            $ pushd $(mktemp -d); go mod init tmp; go get go.astrophena.me/gen; popd

      `go get` puts binaries by default to `$GOPATH/bin` (e.g.
      `~/go/bin`).

      Use `GOBIN` environment variable to change this behavior.

    * Install with `make`:

            $ git clone https://github.com/astrophena/gen
            $ cd gen
            $ make install

        `make install` installs `gen`  by default to `$HOME/.local/bin`.

        Use `PREFIX` environment variable to change this behavior:

            $ make install PREFIX="$HOME" # Installs to $HOME/bin.

## License

[MIT].

[releases page]: https://github.com/astrophena/gen/releases
[Go]: https://golang.org/dl
[MIT]: LICENSE.md
