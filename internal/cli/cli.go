// © 2020 Ilya Mateyko. All rights reserved.
// © 2019 Frédéric Guillot. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE.md file.

// Package cli implements the command line interface of gen.
package cli // import "go.astrophena.me/gen/internal/cli"

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"go.astrophena.me/gen/internal/page"
	"go.astrophena.me/gen/internal/scaffold"
	"go.astrophena.me/gen/internal/version"
	"go.astrophena.me/gen/pkg/fileutil"

	"github.com/urfave/cli/v2"
)

// Run invokes the command line interface of gen.
func Run(args []string) (err error) {
	return App().Run(args)
}

// App returns the structure of the command line interface of gen.
func App() *cli.App {
	return &cli.App{
		Name:    "gen",
		Usage:   "An another static site generator.",
		Version: version.Version,
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
				Name:    "server",
				Aliases: []string{"s"},
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:    "port",
						Aliases: []string{"p"},
						Usage:   "serve at `PORT`",
						Value:   3000,
					},
				},
				Usage:  "Start local HTTP server",
				Action: serverCmd,
			},
			{
				Name:    "new",
				Aliases: []string{"n"},
				Usage:   "Create a new site",
				Action:  newCmd,
			},
		},
	}
}

// buildCmd implements the "build" command.
func buildCmd(c *cli.Context) (err error) {
	var (
		src = c.String("source")
		dst = c.String("destination")

		dirs = map[string]string{
			"pages":     filepath.Join(src, "pages"),
			"templates": filepath.Join(src, "templates"),
			"static":    filepath.Join(src, "static"),
		}
	)

	if err := cleanCmd(c); err != nil {
		return err
	}

	for _, dir := range dirs {
		if !fileutil.Exists(dir) && dir != "static" {
			return fmt.Errorf("%s: doesn't exist, this directory is required for building a site", dir)
		}
	}

	if err := fileutil.Mkdir(dst); err != nil {
		return err
	}

	if fileutil.Exists(dirs["static"]) {
		if err := fileutil.CopyDirContents(dirs["static"], dst); err != nil {
			return err
		}
	}

	tpls, err := fileutil.Files(dirs["templates"], ".html")
	if err != nil {
		return err
	}

	tpl, err := page.ParseTemplates(page.Template(), tpls)
	if err != nil {
		return err
	}

	pages, err := fileutil.Files(dirs["pages"], ".html", ".md")
	if err != nil {
		return err
	}

	for _, pg := range pages {
		pg, err := page.ParseFile(tpl, pg)
		if err != nil {
			return err
		}

		if pg != nil {
			return pg.Generate(tpl, dst)
		}
	}

	return nil
}

// newCmd implements the "new" command.
func newCmd(c *cli.Context) (err error) {
	dst := c.Args().Get(0)

	if dst == "" {
		return fmt.Errorf("directory is required, but not provided")
	}

	return scaffold.Create(dst)
}

// serverCmd implements the "server" command.
func serverCmd(c *cli.Context) (err error) {
	var (
		dst  = c.String("destination")
		port = c.Int("port")

		srv = &http.Server{
			Addr:         fmt.Sprintf("localhost:%v", port),
			WriteTimeout: time.Second * 15,
			ReadTimeout:  time.Second * 15,
			IdleTimeout:  time.Second * 15,
			Handler:      http.FileServer(http.Dir(dst)),
		}
	)

	fmt.Printf("Listening on a port %v...\nUse Ctrl+C to stop.\n", port)
	if err := srv.ListenAndServe(); err != nil {
		if err != http.ErrServerClosed {
			return err
		}
	}

	return nil
}

// cleanCmd implements the "clean" command.
func cleanCmd(c *cli.Context) (err error) {
	dst := c.String("destination")

	if fileutil.Exists(dst) {
		return os.RemoveAll(dst)
	}

	return nil
}
