// © 2020 Ilya Mateyko. All rights reserved.
// © 2019 Frédéric Guillot. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE.md file.

package fileutil // import "astrophena.me/gen/fileutil"

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestCopyDirContents(t *testing.T) {
	f1 := filepath.Join("testdata", "files")
	f2 := filepath.Join("testdata", "files2")

	if Exists(f2) {
		if err := os.RemoveAll(f2); err != nil {
			t.Error(err)
		}
	}

	if err := CopyDirContents(f1, f2); err != nil {
		t.Error(err)
	}

	// TODO(astrophena): Compare two directories.

	if Exists(f2) {
		if err := os.RemoveAll(f2); err != nil {
			t.Error(err)
		}
	}
}

func TestCopyFile(t *testing.T) {
	f1 := filepath.Join("testdata", "copyfile.txt")
	f2 := filepath.Join("testdata", "copyfile2.txt")

	if Exists(f2) {
		if err := os.RemoveAll(f2); err != nil {
			t.Error(err)
		}
	}

	if err := CopyFile(f1, f2); err != nil {
		t.Error(err)
	}

	b1, err := ioutil.ReadFile(f1)
	if err != nil {
		t.Error(err)
	}

	b2, err := ioutil.ReadFile(f2)
	if err != nil {
		t.Error(err)
	}

	if !bytes.Equal(b1, b2) {
		t.Errorf("%s and %s are not equal", b1, b2)
	}

	if Exists(f2) {
		if err := os.RemoveAll(f2); err != nil {
			t.Error(err)
		}
	}
}

func TestExists(t *testing.T) {
	dir := filepath.Join("testdata", "exists")

	if Exists(dir) {
		t.Errorf("%s shouldn't exist", dir)
	}
}

func TestFiles(t *testing.T) {
	dir := filepath.Join("testdata", "files")

	// Keep this synced with testdata/files directory.
	expected := []string{
		// dot.md is excluded.
		"testdata/files/jack.txt",
		"testdata/files/phryne/fisher.txt",
	}

	returned, err := Files(dir, "txt")
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(expected, returned) {
		t.Errorf("expected %s, but returned %s", expected, returned)
	}
}

func TestMkdir(t *testing.T) {
	dir := filepath.Join("testdata", "mkdir")

	if Exists(dir) {
		t.Errorf("%s shouldn't exist", dir)
	}

	if err := Mkdir(dir); err != nil {
		t.Error(err)
	}

	if !Exists(dir) {
		t.Errorf("%s should exist", dir)
	}

	if Exists(dir) {
		if err := os.RemoveAll(dir); err != nil {
			t.Error(err)
		}
	}
}
