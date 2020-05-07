// © 2020 Ilya Mateyko. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE.md file.

// Package site implements site building.
package site // import "go.astrophena.me/gen/internal/site"

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"go.astrophena.me/gen/internal/page"
	"go.astrophena.me/gen/pkg/fileutil"
)

// Build builds the site from the directory src
// to the directory dst, creating it if needed.
func Build(src, dst string) (err error) {
	var (
		dirs = map[string]string{
			"pages":     filepath.Join(src, "pages"),
			"templates": filepath.Join(src, "templates"),
			"static":    filepath.Join(src, "static"),
		}

		start = time.Now()
	)

	// Remove the previously generated site.
	if fileutil.Exists(dst) {
		if err := os.RemoveAll(dst); err != nil {
			return err
		}
	}

	// Check if the required directories exist.
	for _, dir := range dirs {
		if !fileutil.Exists(dir) && dir != "static" {
			return fmt.Errorf("%s: doesn't exist, this directory is required for building a site", dir)
		}
	}

	// Create the site directory.
	if err := fileutil.Mkdir(dst); err != nil {
		return err
	}

	// Copy static files if we have them.
	if fileutil.Exists(dirs["static"]) {
		fmt.Printf("Copying static files...\n")
		if err := fileutil.CopyDirContents(dirs["static"], dst); err != nil {
			return err
		}
	}

	// Parse templates if we have them, otherwise return an error.
	tpls, err := fileutil.Files(dirs["templates"], ".html")
	if err != nil {
		return err
	}

	if len(tpls) < 1 {
		return fmt.Errorf("no templates found in the directory %s", dirs["templates"])
	}

	tpl, err := page.ParseTemplates(page.Template(), tpls)
	if err != nil {
		return err
	}

	// Parse and generate pages if we have them.
	pages, err := fileutil.Files(dirs["pages"], ".html", ".md")
	if err != nil {
		return err
	}

	if len(pages) == 1 {
		fmt.Printf("Parsing and generating %d page...\n", len(pages))
	} else {
		fmt.Printf("Parsing and generating %d pages...\n", len(pages))
	}
	for _, pg := range pages {
		pg, err := page.Parse(tpl, pg)
		if err != nil {
			return err
		}

		if pg != nil {
			if err := pg.Generate(tpl, dst); err != nil {
				return err
			}
		}
	}

	fmt.Printf("Successfully built in %v.\n", time.Since(start))

	return nil
}
