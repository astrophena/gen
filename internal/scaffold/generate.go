// © 2020 Ilya Mateyko. All rights reserved.
// © 2019 Frédéric Guillot. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE.md file.

// +build ignore

package main

import (
	"bytes"
	"fmt"
	"go/format"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"go.astrophena.me/gen/pkg/fileutil"
)

const name = "files.go"

var tpl = template.Must(template.New("").Parse(`// © 2020 Ilya Mateyko. All rights reserved.
// © 2019 Frédéric Guillot. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE.md file.

// Code generated by go generate; DO NOT EDIT.

package scaffold // import "go.astrophena.me/gen/internal/scaffold"

var files = map[string][]byte{
	{{ range $file, $content := . -}}
	{{ printf "%#v" $file }}: {{ printf "%#v" $content }},
	{{ end -}}
}`))

func fatal(err error) {
	fmt.Println(err)
	os.Exit(1)
}

func main() {
	dir := filepath.Join(".", "site")

	files, err := fileutil.Files(dir)
	if err != nil {
		fatal(err)
	}

	fmap := make(map[string][]byte)

	for _, f := range files {
		if filepath.Ext(f) == ".go" {
			continue
		}

		b, err := ioutil.ReadFile(f)
		if err != nil {
			fatal(err)
		}

		p := strings.TrimPrefix(f, "site"+string(os.PathSeparator))
		fmap[p] = b
	}

	var b bytes.Buffer

	if err := tpl.ExecuteTemplate(&b, "", fmap); err != nil {
		fatal(err)
	}

	s, err := format.Source(b.Bytes())
	if err != nil {
		fatal(err)
	}

	if err := ioutil.WriteFile(name, s, 0644); err != nil {
		fatal(err)
	}
}
