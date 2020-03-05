# `gen`

[![License](https://img.shields.io/github/license/astrophena/gen)](LICENSE.md)
[![Go](https://img.shields.io/github/go-mod/go-version/astrophena/gen)](https://golang.org)
[![Release](https://img.shields.io/github/v/release/astrophena/gen)](https://github.com/astrophena/gen/releases)
[![Tests](https://github.com/astrophena/gen/workflows/Tests/badge.svg)](https://github.com/astrophena/gen/actions?query=workflow%3ATests)
[![GoReleaser](https://github.com/astrophena/gen/workflows/GoReleaser/badge.svg)](https://github.com/astrophena/gen/actions?query=workflow%3AGoReleaser)

> **Work in Progress**: `gen` is not finished and has many rough
> edges.

An another static site generator.

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

Download the precompiled binary from [releases page].

### From source

1. Install [Go] 1.14 if you haven't yet.

2. Two installation options are supported:

    * Install with `go get`:

            $ pushd $(mktemp -d); go mod init tmp; go get astrophena.me/gen; popd

      `go get` puts binaries by default to `$GOPATH/bin` (e.g.
      `~/go/bin`).

      Use `GOBIN` environment variable to change that behavior.

    * Install with `make`:

            $ git clone https://github.com/astrophena/gen
            $ cd gen
            $ make install

        `make install` installs `gen`  by default to `$HOME/.local/bin`.

        Use `PREFIX` environment variable to change that behavior:

            $ make install PREFIX="$HOME" # Installs to $HOME/bin.

## License

> `gen` is forked from [plop].
>
> ---
>
> © 2020 Ilya Mateyko. All rights reserved.
>
> © 2019 Frédéric Guillot. All rights reserved.
>
> Use of this source code is governed by the MIT license that can be
> found in the [LICENSE.md] file.

[releases page]: https://github.com/astrophena/gen/releases
[Go]: https://golang.org/dl
[plop]: https://github.com/fguillot/plop
[LICENSE.md]: LICENSE.md
