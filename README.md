<div align="center">
  <h1>gen</h1>
</div>

`gen` is an another static site generator.

## Getting Started

1. [Install](#installation) gen if you haven't yet.

2. Create a new site:

        $ gen new mysite

3. Change directory to `mysite`, build and serve the site locally:

        $ cd mysite
        $ gen build
        $ gen serve

    Run with `--help` or `-h` for options.

4. Go to `http://localhost:3000`.

**ProTip!** You can use [entr] to automatically rebuild the site when changing files:

        $ while true; do find . -type f -not -path '*/\.git/*' | entr -d gen build; done

## Installation

### From binary

Download the precompiled binary from [releases page].

### From source

1. Install the latest version of [Go] if you haven't yet.

2. Install with `go install`:

        $ go install go.astrophena.name/gen@latest

   `go install` puts binaries by default to `$GOPATH/bin` (e.g.
   `~/go/bin`).

   Use `GOBIN` environment variable to change this behavior.

## License

[MIT] © Ilya Mateyko

[entr]: http://eradman.com/entrproject/
[releases page]: https://github.com/astrophena/gen/releases
[Go]: https://golang.org/dl
[MIT]: LICENSE.md
