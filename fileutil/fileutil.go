// © 2020 Ilya Mateyko. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE.md file.

// Package fileutil contains utility functions for working with files
// and directories.
package fileutil

import (
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// CopyDirContents recursively copies contents of src to dst.
func CopyDirContents(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		sp := strings.TrimPrefix(path, src+string(os.PathSeparator))
		dp := filepath.Join(dst, sp)

		if path == src || info.IsDir() {
			return nil
		}

		dir := filepath.Dir(dp)
		if err := Mkdir(dir); err != nil {
			return err
		}

		if err := CopyFile(path, dp); err != nil {
			return err
		}

		return nil
	})
}

// CopyFile copies the file src to dst.
func CopyFile(src, dst string) error {
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

	return nil
}

// Exists returns true if a file or directory does exist and false
// otherwise.
func Exists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

// Files returns a slice of files in the directory dir recursively
// with extensions exts, or an error. If no file extensions are
// supplied, all files are returned.
func Files(dir string, exts ...string) (files []string, err error) {
	return files, filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		if len(exts) > 0 {
			var matches bool
			for _, ext := range exts {
				if filepath.Ext(path) == ext {
					matches = true
				}
			}
			if !matches {
				return nil
			}
		}

		files = append(files, path)

		return nil
	})
}

// Mkdir creates the directory, if it does not already exist. It also
// creates parent directories as needed.
func Mkdir(dir string) error {
	if !Exists(dir) {
		return os.MkdirAll(dir, 0755)
	}
	return nil
}
