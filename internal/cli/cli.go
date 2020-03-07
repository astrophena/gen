// © 2020 Ilya Mateyko. All rights reserved.
// © 2019 Frédéric Guillot. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE.md file.

// Package cli implements command line interface.
package cli // import "astrophena.me/gen/internal/cli"

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"astrophena.me/gen/internal/buildinfo"
	"astrophena.me/gen/internal/page"
	"astrophena.me/gen/internal/scaffold"
	"astrophena.me/gen/pkg/fileutil"

	"github.com/urfave/cli/v2"
)

var minCSS string

// Run invokes command line interface.
func Run() {
	app := &cli.App{
		Name:    "gen",
		Usage:   "An another static site generator.",
		Version: buildinfo.Version,
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
				Usage:   "Performs a one off site build",
				Action:  build,
			},
			{
				Name:    "create",
				Aliases: []string{"c"},
				Usage:   "Creates a new site in the provided directory",
				Action:  create,
			},
			{
				Name:    "serve",
				Aliases: []string{"s"},
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:    "port",
						Aliases: []string{"p"},
						Usage:   "listen on `PORT`",
						Value:   3000,
					},
				},
				Usage:  "Builds site and serves it locally",
				Action: serve,
			},
			{
				Name:    "remove",
				Aliases: []string{"r"},
				Usage:   "Removes all generated files",
				Action:  remove,
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func build(c *cli.Context) (err error) {
	var (
		src = c.String("source")
		dst = c.String("destination")

		assetsDir    = filepath.Join(src, "assets")
		contentDir   = filepath.Join(src, "content")
		templatesDir = filepath.Join(src, "templates")
		staticDir    = filepath.Join(src, "static")

		css = filepath.Join(assetsDir, "sitewide.css")

		start = time.Now()
	)

	fmt.Printf("Building into %s.\n", dst)

	if err := remove(c); err != nil {
		return err
	}

	tpls, err := fileutil.Files(templatesDir, ".html")
	if err != nil {
		return err
	}

	if err := page.ParseTemplates(tpls); err != nil {
		return err
	}

	if err := fileutil.Mkdir(dst); err != nil {
		return err
	}

	if err := fileutil.CopyDirContents(staticDir, dst); err != nil {
		return err
	}

	content, err := fileutil.Files(contentDir, ".html", ".md")
	if err != nil {
		return err
	}

	if fileutil.Exists(css) {
		minCSS, err = page.Minify("text/css", css)
		if err != nil {
			return err
		}
	}

	for _, f := range content {
		p, err := page.ParseFile(f, minCSS)
		if err != nil {
			return err
		}

		if p != nil {
			if err := p.Generate(dst); err != nil {
				return err
			}
		}
	}

	fmt.Printf("Built in %v.\n", time.Since(start))

	return nil
}

func create(c *cli.Context) (err error) {
	dst := c.Args().Get(0)

	if dst == "" {
		return fmt.Errorf("directory is required, but not provided")
	}

	if err := scaffold.Create(dst); err != nil {
		return err
	}

	return nil
}

func serve(c *cli.Context) (err error) {
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

	if err := build(c); err != nil {
		return err
	}

	fmt.Printf("Serving on port %v.\n", port)
	if err := srv.ListenAndServe(); err != nil {
		if err != http.ErrServerClosed {
			return err
		}
	}

	return nil
}

func remove(c *cli.Context) (err error) {
	dst := c.String("destination")

	if fileutil.Exists(dst) {
		if err := os.RemoveAll(dst); err != nil {
			return err
		}
	}

	return nil
}
