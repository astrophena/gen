// © 2020 Ilya Mateyko. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE.md file.

// Package site implements site building.
package site

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"io"
	"log"
	"mime"
	"net/http"
	"os"
	"os/signal"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"syscall"
	"time"

	"go.astrophena.name/gen/fileutil"
	"go.astrophena.name/gen/frontmatter"
	"go.astrophena.name/gen/version"

	"github.com/russross/blackfriday/v2"
	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/css"
	"github.com/tdewolff/minify/html"
	"github.com/tdewolff/minify/js"
	"github.com/tdewolff/minify/json"
	"github.com/tdewolff/minify/svg"
)

// SupportedFormats contains supported page formats.
var SupportedFormats = []string{".html", ".md"}

// TemplateExt is a template file extension.
const TemplateExt = ".tmpl"

// Site represents a site.
type Site struct {
	pages    []*Page
	minify   bool
	src, dst string
	tpl      *template.Template
	quiet    bool
}

func (s *Site) logf(format string, args ...interface{}) {
	if !s.quiet {
		log.Printf(format, args...)
	}
}

// New returns a new site.
func New(src, dst string, quiet, minify bool) (*Site, error) {
	var (
		err error
		s   = &Site{src: src, dst: dst, quiet: quiet, minify: minify}
	)

	for _, dir := range []string{s.pagesDir(), s.templatesDir()} {
		if !fileutil.Exists(dir) {
			return nil, fmt.Errorf("%s: does not exist, this directory is required", dir)
		}
	}

	s.tpl, err = parseTemplates(s.templatesDir())
	if err != nil {
		return nil, err
	}

	return s, nil
}

// Build builds the site.
func (s *Site) Build() error {
	start := time.Now()

	if err := s.Clean(); err != nil {
		return err
	}

	if err := fileutil.Mkdir(s.dst); err != nil {
		return err
	}

	if fileutil.Exists(s.staticDir()) {
		s.logf("Copying static files...")

		if s.minify {
			if err := minifyStaticFiles(s.staticDir(), s.dst); err != nil {
				return err
			}
		} else {
			if err := fileutil.CopyDirContents(s.staticDir(), s.dst); err != nil {
				return err
			}
		}
	}

	pages, err := fileutil.Files(s.pagesDir(), SupportedFormats...)
	if err != nil {
		return err
	}

	if len(pages) > 0 {
		if len(pages) == 1 {
			s.logf("Parsing and generating %d page...", len(pages))
		} else {
			s.logf("Parsing and generating %d pages...", len(pages))
		}
	}

	for _, pp := range pages {
		p, err := s.parsePage(pp)
		if err != nil {
			return err
		}
		s.pages = append(s.pages, p)
	}

	for _, p := range s.pages {
		if err := p.Build(); err != nil {
			return err
		}
	}

	s.logf("Built in %v.", time.Since(start))

	return nil
}

// Clean removes all generated files.
func (s *Site) Clean() (err error) {
	if fileutil.Exists(s.dst) {
		return os.RemoveAll(s.dst)
	}
	return nil
}

// Serve starts local HTTP server, serving the site.
func (s *Site) Serve(addr string) error {
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

	s.logf("Listening on %s.", addr)
	s.logf("Use Ctrl+C to stop.")

	var (
		errc = make(chan error)
		stop = make(chan os.Signal)
	)

	signal.Notify(stop, os.Interrupt)
	signal.Notify(stop, syscall.SIGTERM)

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				errc <- err
			}
		}
	}()

	select {
	case err := <-errc:
		return err
	case <-stop:
		s.logf("Shutting down the server...")

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		return srv.Shutdown(ctx)
	}

	return nil
}

func (s *Site) pagesDir() string     { return filepath.Join(s.src, "pages") }
func (s *Site) staticDir() string    { return filepath.Join(s.src, "static") }
func (s *Site) templatesDir() string { return filepath.Join(s.src, "templates") }

func (s *Site) fs() http.Handler {
	dir := http.Dir(s.dst)
	fs := http.FileServer(dir)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := dir.Open(path.Clean(r.URL.Path))
		if os.IsNotExist(err) {
			s.notFound(w, r)
			return
		}
		fs.ServeHTTP(w, r)
	})
}

func (s *Site) notFound(w http.ResponseWriter, r *http.Request) {
	b, err := os.ReadFile(filepath.Join(s.dst, "404.html"))
	if err != nil {
		http.NotFound(w, r)
		return
	}
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprint(w, string(b))
}

// Page represents a site page.
type Page struct {
	URI         string            `yaml:"uri"`
	Content     string            `yaml:"-"`
	Title       string            `yaml:"title"`
	Description string            `yaml:"description"`
	MetaTags    map[string]string `yaml:"meta_tags"`
	Template    string            `yaml:"template"`

	s *Site // reference to the page owner
}

// Build builds a site page to dst.
func (p *Page) Build() error {
	dir := filepath.Join(p.s.dst, filepath.Dir(p.URI))
	if err := fileutil.Mkdir(dir); err != nil {
		return err
	}

	var buf bytes.Buffer

	if err := p.s.tpl.ExecuteTemplate(&buf, p.Template, p); err != nil {
		return err
	}

	f, err := os.Create(filepath.Join(p.s.dst, p.URI))
	if err != nil {
		return err
	}
	defer f.Close()

	if p.s.minify {
		m := minify.New()
		m.Add("text/html", &html.Minifier{
			KeepDocumentTags: true,
			KeepEndTags:      true,
		})
		return m.Minify("text/html", f, &buf)
	}

	if _, err := buf.WriteTo(f); err != nil {
		return err
	}

	return nil
}

// parsePage parses a file from src and returns a page.
func (s *Site) parsePage(src string) (*Page, error) {
	b, err := os.ReadFile(src)
	if err != nil {
		return nil, err
	}

	p := &Page{MetaTags: make(map[string]string), s: s}

	c, err := frontmatter.Parse(string(b), p)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to parse frontmatter: %w", src, err)
	}

	if s.tpl.Lookup(p.Template) == nil {
		return nil, fmt.Errorf("%s: the template %s specified is not defined", src, p.Template)
	}

	if p.Title == "" || p.Template == "" || p.URI == "" {
		return nil, fmt.Errorf("%s: missing required frontmatter parameter (title, template, uri)", src)
	}

	if !strings.HasSuffix(p.URI, ".html") {
		p.URI = p.URI + "/index.html"
	}

	switch filepath.Ext(src) {
	case ".html":
		p.Content = c
	case ".md":
		p.Content = string(blackfriday.Run([]byte(c)))
	default:
		return nil, fmt.Errorf("%s: format does not supported", src)
	}

	return p, nil
}

// parseTemplates parses templates from dir and returns a template
// that is used for generating pages.
func parseTemplates(dir string) (*template.Template, error) {
	tpls, err := fileutil.Files(dir, TemplateExt)
	if err != nil {
		return nil, err
	}

	if len(tpls) < 1 {
		return nil, fmt.Errorf("no templates found in %s", dir)
	}

	tpl := template.New("site").Funcs(template.FuncMap{
		"content": func(p *Page) template.HTML {
			return template.HTML(p.Content)
		},
		"year": func() int {
			return time.Now().Year()
		},
		"version": func() string {
			return version.Version
		},
	})

	for _, t := range tpls {
		b, err := os.ReadFile(t)
		if err != nil {
			return nil, err
		}

		tpl, err = tpl.Parse(string(b))
		if err != nil {
			return nil, err
		}
	}

	return tpl, nil
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
