// © 2020 Ilya Mateyko. All rights reserved.
// © 2019 Frédéric Guillot. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE.md file.

// An another static site generator.
package main // import "astrophena.me/gen"

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"strings"
	"time"

	"astrophena.me/gen/fileutil"

	"github.com/oxtoacart/bpool"
	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/css"
	"github.com/tdewolff/minify/v2/html"
	"github.com/urfave/cli/v2"
)

const (
	develVersion = "(devel)"
	bufpoolSize  = 48
)

var (
	bufpool *bpool.BufferPool
	m       *minify.M
	minCSS  template.CSS
	tpl     *template.Template
	version string
)

type page struct {
	URI         string
	Title       string
	Description string
	Body        string
	GoImport    string
	GoSource    string

	template string
	filename string
}

func (p *page) Generate(dst string) (err error) {
	dir := filepath.Join(dst, filepath.Dir(p.URI))
	if err := fileutil.Mkdir(dir); err != nil {
		return err
	}

	var (
		tbuf = bufpool.Get()
		mbuf = bufpool.Get()
	)
	defer bufpool.Put(tbuf)
	defer bufpool.Put(mbuf)

	if err := tpl.ExecuteTemplate(tbuf, p.template, p); err != nil {
		return err
	}

	if err := m.Minify("text/html", mbuf, tbuf); err != nil {
		return err
	}

	f, err := os.Create(filepath.Join(dst, p.URI))
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := mbuf.WriteTo(f); err != nil {
		return err
	}

	return nil
}

func init() {
	if version == "" {
		bi, ok := debug.ReadBuildInfo()
		if ok {
			version = strings.TrimPrefix(bi.Main.Version, "v")
		} else {
			version = develVersion
		}
	}

	m = minify.New()
	m.AddFunc("text/html", html.Minify)
	m.AddFunc("text/css", css.Minify)

	bufpool = bpool.NewBufferPool(bufpoolSize)
}

func main() {
	app := &cli.App{
		Name:    "gen",
		Usage:   "An another static site generator.",
		Version: version,
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
				Name:    "clean",
				Aliases: []string{"c"},
				Usage:   "Removes all generated files",
				Action:  clean,
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

		tplFuncs = template.FuncMap{
			"minCSS": func() template.CSS {
				return minCSS
			},
			"noescape": func(s string) template.HTML {
				return template.HTML(s)
			},
			"strdate": func(ts time.Time) string {
				return ts.Format("January 2, 2006")
			},
			"year": func() int {
				return time.Now().Year()
			},
			"version": func() string {
				return fmt.Sprintf("%s, %s (%s/%s)", version, runtime.Version(), runtime.GOOS, runtime.GOARCH)
			},
		}
	)

	tpl = template.New("main").Funcs(tplFuncs)
	if err != nil {
		return err
	}

	tpls, err := fileutil.Files(templatesDir, "html")
	if err != nil {
		return err
	}

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

	fmt.Printf("Building into %s.\n", dst)

	if err := clean(c); err != nil {
		return err
	}

	if err := fileutil.Mkdir(dst); err != nil {
		return err
	}

	if err := fileutil.CopyDirContents(staticDir, dst); err != nil {
		return err
	}

	content, err := fileutil.Files(contentDir, "html")
	if err != nil {
		return err
	}

	if fileutil.Exists(css) {
		minCSS, err = minifyCSS(css)
		if err != nil {
			return err
		}
	}

	for _, f := range content {
		p, err := parseFile(f)
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

	fmt.Printf("Serving site on port %v.\n", port)
	if err := srv.ListenAndServe(); err != nil {
		if err != http.ErrServerClosed {
			return err
		}
	}

	return nil
}

func clean(c *cli.Context) (err error) {
	dst := c.String("destination")

	if fileutil.Exists(dst) {
		if err := os.RemoveAll(dst); err != nil {
			return err
		}
	}

	return nil
}

func parseFile(filename string) (*page, error) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	content := string(b)
	separator := "\n---\n"
	position := strings.Index(content, separator)
	if position <= 0 {
		return nil, fmt.Errorf("%s: no header section detected", filename)
	}

	header := content[:position]
	p := &page{
		Body:     content[position+len(separator):],
		filename: filename,
	}

	for _, line := range strings.Split(header, "\n") {
		switch {
		case strings.HasPrefix(line, "title: "):
			p.Title = line[7:]
		case strings.HasPrefix(line, "description: "):
			p.Description = line[13:]
		case strings.HasPrefix(line, "template: "):
			p.template = line[10:]
		case strings.HasPrefix(line, "uri: "):
			p.URI = line[5:]
		case strings.HasPrefix(line, "go-import: "):
			p.GoImport = line[11:]
		case strings.HasPrefix(line, "go-source: "):
			p.GoSource = line[11:]
		}
	}

	if p.Title == "" || p.template == "" || p.URI == "" {
		return nil, fmt.Errorf("%s: missing required header parameter (title, template, uri)", filename)
	}

	if tpl.Lookup(p.template) == nil {
		return nil, fmt.Errorf("%s: the template %s specified is not defined", filename, p.template)
	}

	return p, nil
}

func minifyCSS(path string) (template.CSS, error) {
	if !fileutil.Exists(path) {
		return "", fmt.Errorf("%s: no such file", path)
	}

	b, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}

	s, err := m.Bytes("text/css", b)
	if err != nil {
		return "", fmt.Errorf("%s: failed to minify: %w", path, err)
	}

	return template.CSS(s), nil
}
