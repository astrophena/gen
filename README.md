<center><h1>`gen`</h1></center>

> **Work in Progress**: `gen` is not finished and has many rough
> edges.

An another static site generator.

## Installation

Download the precompiled binary from [releases page].

### From source

Install [Go] 1.14 if you haven't yet.

Two installation options are supported:

1. Install with `go get`:

        $ pushd $(mktemp -d); go mod init tmp; go get astrophena.me/gen; popd

2. Install with `make`:

        $ git clone https://github.com/astrophena/gen.git
        $ cd gen
        $ make install

    `make install` installs `gen`  by default to `$HOME/.local/bin`.

    Use `PREFIX` environment variable to change that behavior:

        $ make install PREFIX="$HOME" # Installs to $HOME/bin.

## Getting Started

1. [Install](#installation) gen if you haven't yet.

2. Create a new site:

        $ gen new mysite

3. Change directory to `mysite` and serve the site locally:

        $ cd mysite
        $ gen serve

    Run with `--help` or `-h` for options.

4. Go to `http://localhost:3000`.

## Code Status

[![Tests](https://github.com/astrophena/gen/workflows/Tests/badge.svg)](https://github.com/astrophena/gen/actions?query=workflow%3ATests)
[![GoReleaser](https://github.com/astrophena/gen/workflows/GoReleaser/badge.svg)](https://github.com/astrophena/gen/actions?query=workflow%3AGoReleaser)

## License

> [gen] is forked from [plop].
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
