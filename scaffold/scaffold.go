// © 2020 Ilya Mateyko. All rights reserved.
// © 2019 Frédéric Guillot. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE.md file.

// Package scaffold implements a function for generating new sites.
package scaffold // import "astrophena.me/gen/scaffold"

import (
	"io/ioutil"
	"path/filepath"

	"astrophena.me/gen/fileutil"
)

//go:generate go run generate.go

// Generate generates a new site in directory dst or returns an error.
func Generate(dst string) (err error) {
	for filename, content := range files {
		path := filepath.Join(dst, filename)

		dir := filepath.Dir(path)
		if err := fileutil.Mkdir(dir); err != nil {
			return err
		}

		if err := ioutil.WriteFile(path, content, 0644); err != nil {
			return err
		}
	}

	return nil
}
