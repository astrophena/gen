// © 2020 Ilya Mateyko. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE.md file.

// Package page implements page parsing and generation.
package page // import "go.astrophena.name/gen/internal/page"

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"go.astrophena.name/gen/fileutil"
	"go.astrophena.name/gen/frontmatter"
	"go.astrophena.name/gen/internal/version"

	"github.com/russross/blackfriday/v2"
	minifypkg "github.com/tdewolff/minify"
	"github.com/tdewolff/minify/html"
)

// Page represents a page.
type Page struct {
	URI         string            `yaml:"uri"`
	Content     string            `yaml:"-"`
	Title       string            `yaml:"title"`
	Description string            `yaml:"description"`
	MetaTags    map[string]string `yaml:"meta_tags"`
	Template    string            `yaml:"template"`
}

// SupportedFormats contains page formats supported by gen.
var SupportedFormats = []string{".html", ".md"}

// Generate generates HTML from a Page and writes it to the file dst.
func (p *Page) Generate(tpl *template.Template, dst string, minify bool) (err error) {
	dir := filepath.Join(dst, filepath.Dir(p.URI))
	if err := fileutil.Mkdir(dir); err != nil {
		return err
	}

	var pbuf, minbuf bytes.Buffer

	if err := tpl.ExecuteTemplate(&pbuf, p.Template, p); err != nil {
		return err
	}

	if minify {
		m := minifypkg.New()
		m.Add("text/html", &html.Minifier{
			KeepDocumentTags: true,
			KeepEndTags:      true,
		})

		if err := m.Minify("text/html", &minbuf, &pbuf); err != nil {
			return err
		}
	}

	f, err := os.Create(filepath.Join(dst, p.URI))
	if err != nil {
		return err
	}
	defer f.Close()

	if minify {
		if _, err := minbuf.WriteTo(f); err != nil {
			return err
		}
	} else {
		if _, err := pbuf.WriteTo(f); err != nil {
			return err
		}
	}

	return nil
}

// Parse parses a file from src and returns a Page.
func Parse(tpl *template.Template, src string) (*Page, error) {
	b, err := ioutil.ReadFile(src)
	if err != nil {
		return nil, err
	}

	p := &Page{MetaTags: make(map[string]string)}

	c, err := frontmatter.Parse(string(b), p)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to parse frontmatter: %w", src, err)
	}

	if tpl.Lookup(p.Template) == nil {
		return nil, fmt.Errorf("%s: the template %s specified is not defined", src, p.Template)
	}

	if p.Title == "" || p.Template == "" || p.URI == "" {
		return nil, fmt.Errorf("%s: missing required frontmatter parameter (title, template, uri)", src)
	}

	switch ext := filepath.Ext(src); ext {
	case ".html":
		p.Content = c
	case ".md":
		p.Content = string(blackfriday.Run([]byte(c)))
	default:
		return nil, fmt.Errorf("%s: format %s doesn't supported", src, strings.Trim(ext, "."))
	}

	return p, nil
}

// Template returns a template that is used for page generation.
func Template() *template.Template {
	return template.New("").Funcs(template.FuncMap{
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
}

// ParseTemplates parses tpls into a tpl.
func ParseTemplates(tpl *template.Template, tpls []string) (*template.Template, error) {
	for _, t := range tpls {
		f, err := ioutil.ReadFile(t)
		if err != nil {
			return nil, err
		}

		tpl, err = tpl.Parse(string(f))
		if err != nil {
			return nil, err
		}
	}

	return tpl, nil
}
