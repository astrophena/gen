# `gen`

> **Work in Progress**: `gen` is not finished and has many, many
> rough edges. I don't know when `gen` will be finished. Maybe never.

`gen` is an another static site generator.

## Getting Started

1. [Install](#installation) gen if you haven't yet.

2. Create a new site:

        $ gen new mysite

3. Change directory to `mysite`, build and serve the site locally:

        $ cd mysite
        $ gen build
        $ gen server

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

[MIT] © Ilya Mateyko

[releases page]: https://github.com/astrophena/gen/releases
[Go]: https://golang.org/dl
[MIT]: LICENSE.md
