// © 2020 Ilya Mateyko. All rights reserved.
// © 2019 Frédéric Guillot. All rights reserved.
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

	"go.astrophena.me/gen/internal/buildinfo"
	"go.astrophena.me/gen/pkg/fileutil"

	"github.com/russross/blackfriday/v2"
	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/css"
	"github.com/tdewolff/minify/v2/html"
	"gopkg.in/yaml.v2"
)

var (
	m   *minify.M
	tpl *template.Template

	tplFuncs = template.FuncMap{
		"css": func(s string) template.CSS {
			return template.CSS(s)
		},
		"content": func(p *Page) template.HTML {
			var c string
			switch p.Ext {
			case ".md":
				c = string(blackfriday.Run([]byte(p.Body)))
			default:
				c = p.Body
			}
			return template.HTML(c)
		},
		"year": func() int {
			return time.Now().Year()
		},
		"version": buildinfo.TplFunc(),
	}
)

func init() {
	m = minify.New()
	m.AddFunc("text/html", html.Minify)
	m.AddFunc("text/css", css.Minify)

	tpl = template.New("").Funcs(tplFuncs)
}

// Page represents a parsed page.
type Page struct {
	src         string
	Ext         string
	Body        string
	CSS         string
	Frontmatter *Frontmatter
}

// Frontmatter represents a page frontmatter.
type Frontmatter struct {
	URI         string            `yaml:"uri"`
	Title       string            `yaml:"title"`
	Description string            `yaml:"description"`
	MetaTags    map[string]string `yaml:"meta_tags"`
	Template    string            `yaml:"template"`
}

// Generate generates HTML from a parsed page and writes it to dst, returning
// an error otherwise.
func (p *Page) Generate(dst string) (err error) {
	dir := filepath.Join(dst, filepath.Dir(p.Frontmatter.URI))
	if err := fileutil.Mkdir(dir); err != nil {
		return err
	}

	var tbuf, mbuf bytes.Buffer

	if err := tpl.ExecuteTemplate(&tbuf, p.Frontmatter.Template, p); err != nil {
		return err
	}

	if err := m.Minify("text/html", &mbuf, &tbuf); err != nil {
		return err
	}

	f, err := os.Create(filepath.Join(dst, p.Frontmatter.URI))
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := mbuf.WriteTo(f); err != nil {
		return err
	}

	return nil
}

// ParseFile parses page source file and returns Page or an error.
func ParseFile(src string, css string) (*Page, error) {
	b, err := ioutil.ReadFile(src)
	if err != nil {
		return nil, err
	}

	content := string(b)

	separator := "\n---\n"
	position := strings.Index(content, separator)
	if position <= 0 {
		return nil, fmt.Errorf("%s: no header section detected", src)
	}

	frontmatter := content[:position]
	p := &Page{
		Body:        content[position+len(separator):],
		CSS:         css,
		src:         src,
		Ext:         filepath.Ext(src),
		Frontmatter: &Frontmatter{MetaTags: make(map[string]string)},
	}

	if err := yaml.Unmarshal([]byte(frontmatter), p.Frontmatter); err != nil {
		return nil, err
	}

	if p.Frontmatter.Title == "" || p.Frontmatter.Template == "" || p.Frontmatter.URI == "" {
		return nil, fmt.Errorf("%s: missing required frontmatter parameter (title, template, uri)", src)
	}

	if tpl.Lookup(p.Frontmatter.Template) == nil {
		return nil, fmt.Errorf("%s: the template %s specified is not defined", src, p.Frontmatter.Template)
	}

	return p, nil
}

// ParseTemplates parses templates tpls into template tpl, returning an error
// otherwise.
func ParseTemplates(tpls []string) error {
	for _, t := range tpls {
		f, err := ioutil.ReadFile(t)
		if err != nil {
			return err
		}

		tpl, err = tpl.Parse(string(f))
		if err != nil {
			return err
		}
	}

	return nil
}

// Minify minifies file src with the type mimetype, returning an error
// otherwise.
func Minify(mimetype, src string) (string, error) {
	if !fileutil.Exists(src) {
		return "", fmt.Errorf("%s: no such file", src)
	}

	b, err := ioutil.ReadFile(src)
	if err != nil {
		return "", err
	}

	s, err := m.Bytes(mimetype, b)
	if err != nil {
		return "", fmt.Errorf("%s: failed to minify: %w", src, err)
	}

	return string(s), nil
}
