// © 2020 Ilya Mateyko. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE.md file.

// Package page implements page parsing and generation.
package page // import "go.astrophena.me/gen/internal/page"

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"go.astrophena.me/gen/internal/version"
	"go.astrophena.me/gen/pkg/fileutil"

	"github.com/russross/blackfriday/v2"
	"gopkg.in/yaml.v2"
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

// Generate generates HTML from a Page and writes it to the file by the
// path dst, returning an error otherwise.
func (p *Page) Generate(tpl *template.Template, dst string) (err error) {
	dir := filepath.Join(dst, filepath.Dir(p.URI))
	if err := fileutil.Mkdir(dir); err != nil {
		return err
	}

	var buf bytes.Buffer

	if err := tpl.ExecuteTemplate(&buf, p.Template, p); err != nil {
		return err
	}

	f, err := os.Create(filepath.Join(dst, p.URI))
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := buf.WriteTo(f); err != nil {
		return err
	}

	return nil
}

// Parse parses a file and returns Page or an error.
func Parse(tpl *template.Template, src string) (*Page, error) {
	b, err := ioutil.ReadFile(src)
	if err != nil {
		return nil, err
	}

	// TODO: Improve frontmatter detection code.
	all := string(b)

	separator := "\n---\n"
	position := strings.Index(all, separator)
	if position <= 0 {
		return nil, fmt.Errorf("%s: no frontmatter detected", src)
	}

	frontmatter := all[:position]
	content := all[position+len(separator):]

	p := &Page{
		MetaTags: make(map[string]string),
	}

	switch filepath.Ext(src) {
	case ".html":
		p.Content = content
	case ".md":
		p.Content = string(blackfriday.Run([]byte(content)))
	default:
		return nil, fmt.Errorf("%s: format doesn't supported", src)
	}

	if err := yaml.Unmarshal([]byte(frontmatter), p); err != nil {
		return nil, err
	}

	if p.Title == "" || p.Template == "" || p.URI == "" {
		return nil, fmt.Errorf("%s: missing required frontmatter parameter (title, template, uri)", src)
	}

	if tpl.Lookup(p.Template) == nil {
		return nil, fmt.Errorf("%s: the template %s specified is not defined", src, p.Template)
	}

	return p, nil
}

// Template returns a *template.Template that is used for generating pages.
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

// ParseTemplates parses tpls into a tpl, returning it back or an error.
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
