// © 2020 Ilya Mateyko. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE.md file.

// Package cli implements the command line interface of gen.
package cli

import (
	"errors"

	"go.astrophena.name/gen/scaffold"
	"go.astrophena.name/gen/site"
	"go.astrophena.name/gen/version"

	"github.com/urfave/cli/v2"
)

// Run invokes the command line interface of gen.
func Run(args []string) error {
	app := &cli.App{
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
			&cli.BoolFlag{
				Name:    "minify",
				Aliases: []string{"m", "min"},
				Usage:   "minify files",
				Value:   false,
			},
			&cli.BoolFlag{
				Name:    "quiet",
				Aliases: []string{"q"},
				Usage:   "quiet mode",
				Value:   false,
			},
		},
		Commands: []*cli.Command{
			{
				Name:   "build",
				Usage:  "Perform a one-off site build",
				Action: build,
			},
			{
				Name:   "clean",
				Usage:  "Remove all generated files",
				Action: clean,
			},
			{
				Name: "serve",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "addr",
						Aliases: []string{"a"},
						Usage:   "listen at `host:port`",
						Value:   "localhost:3000",
					},
				},
				Usage:  "Build and serve the site locally",
				Action: serve,
			},
			{
				Name:  "new",
				Usage: "Generate a new site",
				Action: func(c *cli.Context) error {
					dst := c.Args().Get(0)

					if dst == "" {
						return errors.New("directory is required")
					}

					return scaffold.Create(dst)
				},
			},
		},
	}

	return app.Run(args)
}

func newSite(c *cli.Context) (*site.Site, error) {
	return site.New(c.String("source"), c.String("destination"), c.Bool("quiet"), c.Bool("minify"))
}

func build(c *cli.Context) error {
	s, err := newSite(c)
	if err != nil {
		return err
	}
	return s.Build()
}

func clean(c *cli.Context) error {
	s, err := newSite(c)
	if err != nil {
		return err
	}
	return s.Clean()
}

func serve(c *cli.Context) error {
	s, err := newSite(c)
	if err != nil {
		return err
	}
	return s.Serve(c.String("addr"))
}
