# `gen`

> **Work in Progress**: `gen` is not finished and has many rough
> edges.

> Forked from [plop].

An another static site generator.

## Installation

Install the [precompiled binary].

### From source

[Go] 1.14 is required.

```bash
pushd $(mktemp -d); go mod init tmp; go get astrophena.me/gen; popd
```

Or:

```bash
git clone https://github.com/astrophena/gen.git
cd gen
make install
```

`make install` installs `gen`  by default to `$HOME/.local/bin`.

Use `PREFIX` environment variable to change that behavior:

```bash
make install PREFIX="$HOME" # Installs to $HOME/bin.
```

## Getting Started

1. [Install](#installation) gen if you haven't yet.

2. Create a new site:

        $ gen new mysite

3. Change directory to `mysite` and serve the site locally:

        $ cd mysite
        $ gen serve

    Run with `--help` or `-h` for options.

4. Go to `http://localhost:3000`.

## License

© 2020 Ilya Mateyko. All rights reserved.

© 2019 Frédéric Guillot. All rights reserved.

Use of this source code is governed by the MIT license that can be
found in the [LICENSE.md] file.

[plop]: https://github.com/fguillot/plop
[precompiled binary]: https://github.com/astrophena/gen/releases
[Go]: https://golang.org
[LICENSE.md]: LICENSE.md
