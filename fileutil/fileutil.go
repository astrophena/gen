// © 2020 Ilya Mateyko. All rights reserved.
// © 2019 Frédéric Guillot. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE.md file.

// Package fileutil implements some functions for working with files.
package fileutil // import "astrophena.me/gen/fileutil"

import (
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Browse returns a slice of file paths in the directory dir recursively, with
// file extensions exts or an error in case of failure. If no file extensions
// are supplied, all files are included.
func Browse(dir string, exts ...string) (files []string, err error) {
	return files, filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		for _, ext := range exts {
			if !strings.HasSuffix(path, "."+ext) {
				return nil
			}
		}

		files = append(files, path)

		return nil
	})
}

// CopyDirContents recursively copies contents of the src directory to dst.
func CopyDirContents(src, dst string) (err error) {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		sp := strings.TrimPrefix(path, src+string(os.PathSeparator))
		dp := filepath.Join(dst, sp)

		if path == src || info.IsDir() {
			return nil
		}

		dir := filepath.Dir(dp)
		if err := MkDir(dir); err != nil {
			return err
		}

		if err := CopyFile(path, dp); err != nil {
			return err
		}

		return nil
	})
}

// Exists returns true if a file or directory exists and false
// otherwise.
func Exists(path string) (exists bool) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

// MkDir creates a directory if it does not exist.
func MkDir(dir string) (err error) {
	if !Exists(dir) {
		return os.MkdirAll(dir, 0755)
	}
	return nil
}

// CopyFile copies the src file to dst. Any existing file will be overwritten
// and it will not copy file attributes.
func CopyFile(src, dst string) (err error) {
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
