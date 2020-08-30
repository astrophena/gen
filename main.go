// © 2020 Ilya Mateyko. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE.md file.

// gen is an another static site generator.
package main // import "github.com/astrophena/gen"

import (
	"fmt"
	"os"

	"github.com/astrophena/gen/internal/cli"

	"github.com/logrusorgru/aurora"
)

func main() {
	if err := cli.Run(os.Args); err != nil {
		fmt.Println(aurora.Red(err.Error()))
		os.Exit(1)
	}
}
