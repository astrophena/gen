# gen

Another static site generator.

Forked from [plop].


## Building

[Go] 1.13 or newer is required.

```sh
$ make
```

## Installation

Install the [pre-compiled binary](https://github.com/astrophena/gen/releases).

### From source

```sh
$ make install PREFIX="$HOME" # installs to $HOME/bin/gen
```

### [Scoop]

```sh
$ scoop bucket add gen https://github.com/astrophena/gen.git
$ scoop install gen
```

## License

Copyright 2020 Ilya Mateyko. All rights reserved.

Copyright 2019 Frédéric Guillot. All rights reserved.

Use of this source code is governed by the MIT license that can be found in the LICENSE file.

[plop]: https://github.com/fguillot/plop
[Go]: https://golang.org/dl/
[Scoop]: https://scoop.sh
