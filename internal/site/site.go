// © 2020 Ilya Mateyko. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE.md file.

// Package site implements site building.
package site // import "go.astrophena.name/gen/internal/site"

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"time"

	"go.astrophena.name/gen/fileutil"
	"go.astrophena.name/gen/internal/page"

	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/css"
	"github.com/tdewolff/minify/js"
	"github.com/tdewolff/minify/json"
	"github.com/tdewolff/minify/svg"
)

// Site represents a site.
type Site struct {
	minify bool

	src string
	dst string
}

// New returns a new Site.
func New(src, dst string, minify bool) *Site {
	return &Site{src: src, dst: dst, minify: minify}
}

// Build builds the Site.
func (s *Site) Build() (err error) {
	var (
		pagesDir     = filepath.Join(s.src, "pages")
		templatesDir = filepath.Join(s.src, "templates")
		staticDir    = filepath.Join(s.src, "static")

		start = time.Now()
	)

	for _, dir := range []string{pagesDir, templatesDir, staticDir} {
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
		log.Print("Copying static files...")
		if s.minify {
			if err := minifyStaticFiles(staticDir, s.dst); err != nil {
				return err
			}
		} else {
			if err := fileutil.CopyDirContents(staticDir, s.dst); err != nil {
				return err
			}
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
			log.Printf("Parsing and generating %d page...\n", len(pages))
		} else {
			log.Printf("Parsing and generating %d pages...\n", len(pages))
		}
	}

	for _, p := range pages {
		p, err := page.Parse(tpl, p)
		if err != nil {
			return err
		}

		if p != nil {
			if err := p.Generate(tpl, s.dst, s.minify); err != nil {
				return err
			}
		}
	}

	log.Printf("Built in %v.", time.Since(start))

	return nil
}

func minifyStaticFiles(src, dst string) error {
	files, err := fileutil.Files(src)
	if err != nil {
		return err
	}

	m := minify.New()
	m.AddFunc("text/css", css.Minify)
	m.AddFuncRegexp(regexp.MustCompile("^(application|text)/(x-)?(java|ecma)script$"), js.Minify)
	m.AddFuncRegexp(regexp.MustCompile("[/+]json$"), json.Minify)
	m.AddFunc("image/svg+xml", svg.Minify)

	for _, file := range files {
		mimetype := mime.TypeByExtension(filepath.Ext(file))

		from, err := os.Open(file)
		if err != nil {
			return err
		}
		defer from.Close()

		dir, err := filepath.Rel(src, file)
		if err != nil {
			return err
		}
		dir = filepath.Dir(dir)

		if err := fileutil.Mkdir(filepath.Join(dst, dir)); err != nil {
			return err
		}

		to, err := os.Create(filepath.Join(dst, dir, filepath.Base(file)))
		if err != nil {
			return err
		}
		defer to.Close()

		err = m.Minify(mimetype, to, from)
		if err != nil {
			if err != minify.ErrNotExist {
				return err
			}

			if _, err := io.Copy(to, from); err != nil {
				return err
			}
		}
	}

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
		Handler:      s.fs(),
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

func (s *Site) fs() http.Handler {
	fs := http.Dir(s.dst)
	fsrv := http.FileServer(fs)

	// See https://stackoverflow.com/a/62747667.
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := fs.Open(path.Clean(r.URL.Path)) // Do not allow path traversals.
		if os.IsNotExist(err) {
			s.notFound(w, r)
			return
		}
		fsrv.ServeHTTP(w, r)
	})
}

func (s *Site) notFound(w http.ResponseWriter, r *http.Request) {
	nf, err := ioutil.ReadFile(filepath.Join(s.dst, "404.html"))
	if err != nil {
		http.NotFound(w, r)
		return
	}

	w.WriteHeader(http.StatusNotFound)
	fmt.Fprint(w, string(nf))
}
