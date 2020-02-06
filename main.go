// © 2020 Ilya Mateyko. All rights reserved.
// © 2019 Frédéric Guillot. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE.md file.

// An another static site generator.
package main // import "astrophena.me/gen"

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/css"
	"github.com/tdewolff/minify/v2/html"
	"github.com/urfave/cli/v2"
)

type pageMetadata struct {
	URI         string
	Title       string
	Description string
	Body        string

	template string
	filename string
}

var (
	tpl *template.Template
	m   *minify.M

	minifiedCSS template.CSS

	version = "dev"
)

func init() {
	m = minify.New()
	m.AddFunc("text/html", html.Minify)
	m.AddFunc("text/css", css.Minify)
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
	srcDir := c.String("src")
	tplDir := c.String("tpl")
	outDir := c.String("out")
	pubDir := c.String("pub")
	cssPath := c.String("css")

	if err := mkDir(outDir); err != nil {
		return err
	}

	fmt.Printf("building into %s\n", outDir)
	start := time.Now()

	tplFuncs := template.FuncMap{
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

	tpl, err = template.New("main").Funcs(tplFuncs).ParseGlob(tplDir + "/*.html")
	if err != nil {
		return err
	}

	srcs, err := browse(srcDir, "html")
	if err != nil {
		return err
	}

	minifiedCSS, err = minifyCSS(cssPath)
	if err != nil {
		return err
	}

	for _, src := range srcs {
		page, err := parseFile(src)
		if err != nil {
			return err
		}

		if page != nil {
			if err := generatePage(page, outDir); err != nil {
				return err
			}
		}
	}

	if err := copyDirContents(pubDir, outDir); err != nil {
		return err
	}

	fmt.Printf("built in %v\n", time.Since(start))

	return nil
}

func serve(c *cli.Context) error {
	port := c.Int("port")
	dir := c.String("dir")

	if !exists(dir) {
		return fmt.Errorf("%s does not exist: run \"gen build\" to build the site", dir)
	}

	handler := http.FileServer(http.Dir(dir))

	srv := &http.Server{
		Addr:         fmt.Sprintf("localhost:%v", port),
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 15,
		Handler:      handler,
	}

	fmt.Printf("serving %s on port %v…\n", dir, port)
	if err := srv.ListenAndServe(); err != nil {
		if err != http.ErrServerClosed {
			return err
		}
	}

	return nil
}

func parseFile(filename string) (*pageMetadata, error) {
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
	page := &pageMetadata{
		Body:     content[position+len(separator):],
		filename: filename,
	}

	for _, line := range strings.Split(header, "\n") {
		switch {
		case strings.HasPrefix(line, "title: "):
			page.Title = line[7:]
		case strings.HasPrefix(line, "description: "):
			page.Description = line[13:]
		case strings.HasPrefix(line, "template: "):
			page.template = line[10:]
		case strings.HasPrefix(line, "uri: "):
			page.URI = line[5:]
		}
	}

	if page.Title == "" || page.template == "" || page.URI == "" {
		return nil, fmt.Errorf("%s: missing required header parameter (title, template, uri)", filename)
	}

	if tpl.Lookup(page.template) == nil {
		return nil, fmt.Errorf("%s: the template %s specified is not defined", filename, page.template)
	}

	return page, nil
}

func browse(dir string, ext string) (files []string, err error) {
	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(path, "."+ext) {
			files = append(files, path)
		}
		return nil
	})

	return files, err
}

func generatePage(page *pageMetadata, outDir string) error {
	if err := mkDir(filepath.Join(outDir, filepath.Dir(page.URI))); err != nil {
		return err
	}

	var tplBuf, minBuf bytes.Buffer
	if err := tpl.ExecuteTemplate(&tplBuf, page.template, page); err != nil {
		return fmt.Errorf("%s: failed to generate: %w", page.URI, err)
	}

	if err := m.Minify("text/html", &minBuf, &tplBuf); err != nil {
		return fmt.Errorf("%s: failed to minify: %w", page.URI, err)
	}

	file, err := os.Create(filepath.Join(outDir, page.URI))
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err := minBuf.WriteTo(file); err != nil {
		return fmt.Errorf("%s: failed to write: %w", page.URI, err)
	}

	return nil
}

func minifyCSS(path string) (template.CSS, error) {
	if !exists(path) {
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

// copyDirContents recursively copies contents of the src directory to dst.
func copyDirContents(src, dst string) (err error) {
	return filepath.Walk(src, func(file string, info os.FileInfo, err error) error {
		srcPath := strings.TrimPrefix(file, src+string(os.PathSeparator))
		dstPath := filepath.Join(dst, srcPath)

		if file == src || info.IsDir() {
			return nil
		}

		dir := filepath.Dir(dstPath)
		if err := mkDir(dir); err != nil {
			return fmt.Errorf("%s: failed to create: %w", dir, err)
		}

		if err := copyFile(file, dstPath); err != nil {
			return fmt.Errorf("%s: failed to copy to %s: %w", file, dst, err)
		}

		return nil
	})
}

// exists returns true if a file or directory exists and false
// otherwise.
func exists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

// mkDir creates a directory if it does not exist.
func mkDir(dir string) error {
	if !exists(dir) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}
	return nil
}

// copyFile copies the src file to dst. Any existing file will be
// overwritten and it will not copy file attributes.
func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return out.Close()
}
