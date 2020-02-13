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
	tpl         *template.Template
	m           *minify.M
	minifiedCSS template.CSS
	version     string
	bufpool     *bpool.BufferPool
)

type page struct {
	URI         string
	Title       string
	Description string
	Body        string

	template string
	filename string
}

func (p *page) Generate(dst string) (err error) {
	dir := filepath.Join(dst, filepath.Dir(p.URI))

	if err := fileutil.MkDir(dir); err != nil {
		return err
	}

	// TODO(astrophena): Maybe don't use defer there?

	tbuf := bufpool.Get()
	defer bufpool.Put(tbuf)

	mbuf := bufpool.Get()
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
	// If version was not already embedded, try to get it from module
	// version.
	if version == "" {
		bi, ok := debug.ReadBuildInfo()
		if ok {
			version = strings.TrimPrefix(bi.Main.Version, "v")
		} else {
			version = develVersion
		}
	}

	// Set up minifiers.
	m = minify.New()
	m.AddFunc("text/html", html.Minify)
	m.AddFunc("text/css", css.Minify)

	// Set up buffer pool.
	bufpool = bpool.NewBufferPool(bufpoolSize)
}

func main() {
	app := &cli.App{
		Name:    "gen",
		Usage:   "an another static site generator",
		Version: version,
		Authors: []*cli.Author{
			{
				Name:  "Ilya Mateyko",
				Email: "inbox@astrophena.me",
			},
		},
		Commands: []*cli.Command{
			{
				Name:    "build",
				Aliases: []string{"b"},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "src",
						Aliases: []string{"s"},
						Usage:   "use sources from `DIR`",
						Value:   "src",
					},
					&cli.StringFlag{
						Name:    "tpl",
						Aliases: []string{"t"},
						Usage:   "use templates from `DIR`",
						Value:   "tpl",
					},
					&cli.StringFlag{
						Name:    "css",
						Aliases: []string{"c"},
						Usage:   "use CSS from `FILE`",
						Value:   "sitewide.css",
					},
					&cli.StringFlag{
						Name:    "pub",
						Aliases: []string{"p"},
						Usage:   "copy files from `DIR`",
						Value:   "pub",
					},
					&cli.StringFlag{
						Name:    "out",
						Aliases: []string{"o"},
						Usage:   "place built files into `DIR`",
						Value:   "out",
					},
				},
				Usage:  "Build the site",
				Action: build,
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
					&cli.StringFlag{
						Name:    "dir",
						Aliases: []string{"d"},
						Usage:   "serve from `DIR`",
						Value:   "out",
					},
				},
				Usage:  "Serve site locally",
				Action: serve,
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
		srcDir = c.String("src")
		tplDir = c.String("tpl")
		outDir = c.String("out")
		pubDir = c.String("pub")

		cssPath = c.String("css")

		start = time.Now()

		tplFuncs = template.FuncMap{
			"noescape": func(s string) template.HTML {
				return template.HTML(s)
			},
			"strdate": func(ts time.Time) string {
				return ts.Format("January 2, 2006")
			},
			"version": func() string {
				return fmt.Sprintf("%s, %s (%s/%s)", version, runtime.Version(), runtime.GOOS, runtime.GOARCH)
			},
			"minifiedCSS": func() template.CSS {
				return minifiedCSS
			},
			"year": func() int {
				return time.Now().Year()
			},
		}
	)

	fmt.Printf("Building into %s.\n", outDir)

	if err := fileutil.MkDir(outDir); err != nil {
		return err
	}

	if err := fileutil.CopyDirContents(pubDir, outDir); err != nil {
		return err
	}

	tpl, err = template.New("main").Funcs(tplFuncs).ParseGlob(tplDir + "/*.html")
	if err != nil {
		return err
	}

	srcs, err := fileutil.Browse(srcDir, "html")
	if err != nil {
		return err
	}

	minifiedCSS, err = minifyCSS(cssPath)
	if err != nil {
		return err
	}

	for _, src := range srcs {
		p, err := parseFile(src)
		if err != nil {
			return err
		}

		if p != nil {
			if err := p.Generate(outDir); err != nil {
				return err
			}
		}
	}

	fmt.Printf("Built in %v.\n", time.Since(start))

	return nil
}

func serve(c *cli.Context) error {
	var (
		port = c.Int("port")
		dir  = c.String("dir")
	)

	if !fileutil.Exists(dir) {
		build(c)
	}

	handler := http.FileServer(http.Dir(dir))

	srv := &http.Server{
		Addr:         fmt.Sprintf("localhost:%v", port),
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 15,
		Handler:      handler,
	}

	fmt.Printf("Serving %s on port %v.\n", dir, port)
	if err := srv.ListenAndServe(); err != nil {
		if err != http.ErrServerClosed {
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
