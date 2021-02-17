// © 2020 Ilya Mateyko. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE.md file.

// Package scaffold implements the creation of new sites.
package scaffold // import "go.astrophena.name/gen/scaffold"

import (
	"embed"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"go.astrophena.name/gen/fileutil"
)

//go:embed site/*
var site embed.FS

// Create creates a new site in the path dst.
func Create(dst string) (err error) {
	return fs.WalkDir(site, "site", func(path string, d fs.DirEntry, err error) error {
		sp := strings.TrimPrefix(path, "site"+string(os.PathSeparator))
		dp := filepath.Join(dst, sp)

		if d.IsDir() {
			return nil
		}

		dir := filepath.Dir(dp)
		if err := fileutil.Mkdir(dir); err != nil {
			return err
		}

		in, err := site.Open(path)
		if err != nil {
			return err
		}
		defer in.Close()

		out, err := os.Create(dp)
		if err != nil {
			return err
		}

		_, err = io.Copy(out, in)
		if err != nil {
			return err
		}
		return out.Close()
	})
}
