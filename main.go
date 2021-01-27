// © 2020 Ilya Mateyko. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE.md file.

// gen is an another static site generator.
package main

import (
	"log"
	"os"

	"go.astrophena.name/gen/internal/cli"
)

func main() {
	log.SetFlags(0)

	if err := cli.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
