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

	"astrophena.me/gen/fileutil"
)

const name = "files.go"

var (
	tpl = template.Must(template.New("").Parse(`// © 2020 Ilya Mateyko. All rights reserved.
// © 2019 Frédéric Guillot. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE.md file.

// Code generated by go generate; DO NOT EDIT.

package scaffold // import "astrophena.me/gen/scaffold"

var files = map[string][]byte{
	{{ range $filename, $content := . -}}
	{{ printf "%#v" $filename }}: {{ printf "%#v" $content }},
	{{ end -}}
}`))

	filesMap = make(map[string][]byte)
)

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

	for _, filename := range files {
		if filepath.Ext(filename) == ".go" {
			continue
		}

		b, err := ioutil.ReadFile(filename)
		if err != nil {
			fatal(err)
		}

		path := strings.TrimPrefix(filename, "site"+string(os.PathSeparator))
		filesMap[path] = b
	}

	var buf bytes.Buffer

	if err := tpl.ExecuteTemplate(&buf, "", filesMap); err != nil {
		fatal(err)
	}

	src, err := format.Source(buf.Bytes())
	if err != nil {
		fatal(err)
	}

	if err := ioutil.WriteFile(name, src, 0644); err != nil {
		fatal(err)
	}
}
