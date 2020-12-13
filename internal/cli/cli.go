// © 2020 Ilya Mateyko. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE.md file.

// Package cli implements the command line interface of gen.
package cli // import "go.astrophena.name/gen/internal/cli"

import (
	"errors"

	"go.astrophena.name/gen/internal/scaffold"
	"go.astrophena.name/gen/internal/site"
	"go.astrophena.name/gen/internal/version"

	"github.com/urfave/cli/v2"
)

// Run invokes the command line interface of gen.
func Run(args []string) (err error) {
	return app().Run(args)
}

func app() *cli.App {
	return &cli.App{
		Name:                 "gen",
		Usage:                "An another static site generator.",
		Version:              version.Version,
		EnableBashCompletion: true,
		HideHelpCommand:      true,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "source",
				Aliases: []string{"s", "src"},
				Usage:   "read files from `DIR`",
				Value:   ".",
			},
			&cli.StringFlag{
				Name:    "destination",
				Aliases: []string{"d", "dst"},
				Usage:   "write files to `DIR`",
				Value:   "site",
			},
		},
		Commands: []*cli.Command{
			{
				Name:    "build",
				Aliases: []string{"b"},
				Usage:   "Perform a one-off site build",
				Action:  buildCmd,
			},
			{
				Name:    "clean",
				Aliases: []string{"c"},
				Usage:   "Remove all generated files",
				Action:  cleanCmd,
			},
			{
				Name:    "serve",
				Aliases: []string{"s"},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "addr",
						Aliases: []string{"a"},
						Usage:   "listen at `host:port`",
						Value:   "localhost:3000",
					},
				},
				Usage:  "Build and serve the site locally",
				Action: serveCmd,
			},
			{
				Name:    "new",
				Aliases: []string{"n"},
				Usage:   "Generate a new site",
				Action:  newCmd,
			},
		},
	}
}

func buildCmd(c *cli.Context) (err error) {
	return site.New(c.String("source"), c.String("destination")).Build()
}

func cleanCmd(c *cli.Context) (err error) {
	return site.New(c.String("source"), c.String("destination")).Clean()
}

func serveCmd(c *cli.Context) (err error) {
	return site.New(c.String("source"), c.String("destination")).Serve(c.String("addr"))
}

func newCmd(c *cli.Context) (err error) {
	dst := c.Args().Get(0)

	if dst == "" {
		return errors.New("directory is required, but not provided")
	}

	return scaffold.Create(dst)
}
