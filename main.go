// © 2020 Ilya Mateyko. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE.md file.

// gen is an another static site generator.
package main // import "go.astrophena.me/gen"

import (
	"fmt"
	"os"

	"go.astrophena.me/gen/internal/cli"
)

func main() {
	if err := cli.Run(os.Args); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
