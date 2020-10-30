// © 2020 Ilya Mateyko. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE.md file.

// Package site implements site building.
package site // import "go.astrophena.name/gen/internal/site"

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"go.astrophena.name/gen/internal/page"
	"go.astrophena.name/gen/pkg/fileutil"
)

// Site represents a site.
type Site struct {
	src string
	dst string
}

// New returns a new Site.
func New(src, dst string) *Site {
	return &Site{src: src, dst: dst}
}

// Build builds the Site.
func (s *Site) Build() (err error) {
	var (
		pagesDir     = filepath.Join(s.src, "pages")
		templatesDir = filepath.Join(s.src, "templates")
		staticDir    = filepath.Join(s.src, "static")

		start = time.Now()
	)

	dirs := []string{pagesDir, templatesDir, staticDir}
	for _, dir := range dirs {
		if !fileutil.Exists(dir) && dir != "static" {
			return fmt.Errorf("%s: doesn't exist, this directory is required for building a site", dir)
		}
	}

	if err := s.Clean(); err != nil {
		return err
	}

	if err := fileutil.Mkdir(s.dst); err != nil {
		return err
	}

	if fileutil.Exists(staticDir) {
		fmt.Println("Copying static files...")
		if err := fileutil.CopyDirContents(staticDir, s.dst); err != nil {
			return err
		}
	}

	tpls, err := fileutil.Files(templatesDir, ".html")
	if err != nil {
		return err
	}

	if len(tpls) < 1 {
		return fmt.Errorf("no templates found in the directory %s", templatesDir)
	}

	tpl, err := page.ParseTemplates(page.Template(), tpls)
	if err != nil {
		return err
	}

	pages, err := fileutil.Files(pagesDir, page.SupportedFormats...)
	if err != nil {
		return err
	}

	if len(pages) > 0 {
		if len(pages) == 1 {
			fmt.Printf("Parsing and generating %d page...\n", len(pages))
		} else {
			fmt.Printf("Parsing and generating %d pages...\n", len(pages))
		}
	}

	for _, p := range pages {
		p, err := page.Parse(tpl, p)
		if err != nil {
			return err
		}

		if p != nil {
			if err := p.Generate(tpl, s.dst); err != nil {
				return err
			}
		}
	}

	log.Printf("Built in %v.", time.Since(start))

	return nil
}

// Clean removes all generated files.
func (s *Site) Clean() (err error) {
	if fileutil.Exists(s.dst) {
		return os.RemoveAll(s.dst)
	}
	return nil
}

// Serve starts local HTTP server, serving the Site.
func (s *Site) Serve(addr string) (err error) {
	srv := &http.Server{
		Addr:         addr,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 15,
		Handler:      http.FileServer(http.Dir(s.dst)),
	}

	if err := s.Build(); err != nil {
		return err
	}

	log.Printf("Listening on %s.", addr)
	log.Println("Use Ctrl+C to stop.")
	if err := srv.ListenAndServe(); err != nil {
		if err != http.ErrServerClosed {
			return err
		}
	}

	return nil
}
