// © 2020 Ilya Mateyko. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE.md file.

package scaffold // import "go.astrophena.me/gen/internal/scaffold"

import (
	"io/ioutil"
	"path/filepath"

	"go.astrophena.me/gen/pkg/fileutil"
)

//go:generate go run generate.go

// Create creates a new site in the directory dst or returns an error.
func Create(dst string) (err error) {
	for name, content := range files {
		path := filepath.Join(dst, name)

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
