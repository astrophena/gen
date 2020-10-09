// © 2020 Ilya Mateyko. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE.md file.

package fileutil_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"go.astrophena.name/gen/pkg/fileutil"
)

func remove(t *testing.T, f string) func() {
	return func() {
		if fileutil.Exists(f) {
			if err := os.RemoveAll(f); err != nil {
				t.Error(err)
			}
		}
	}
}

func TestCopyDirContents(t *testing.T) {
	f1 := filepath.Join("testdata", "files")
	f2 := filepath.Join("testdata", "files2")

	if err := fileutil.CopyDirContents(f1, f2); err != nil {
		t.Error(err)
	}
	t.Cleanup(remove(t, f2))
}

func ExampleCopyDirContents() {
	if err := fileutil.CopyDirContents("phryne", "fisher"); err != nil {
		log.Fatal(err)
	}
}

func TestCopyFile(t *testing.T) {
	f1 := filepath.Join("testdata", "copyfile.txt")
	f2 := filepath.Join("testdata", "copyfile2.txt")

	if err := fileutil.CopyFile(f1, f2); err != nil {
		t.Error(err)
	}
	t.Cleanup(remove(t, f2))

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
}

func ExampleCopyFile() {
	if err := fileutil.CopyFile("/etc/hostname", "/tmp/hostname"); err != nil {
		log.Fatal(err)
	}
}

func TestExists(t *testing.T) {
	dir := filepath.Join("testdata", "exists")

	if fileutil.Exists(dir) {
		t.Errorf("%s shouldn't exist", dir)
	}
}

func ExampleExists() {
	if fileutil.Exists("example") {
		fmt.Println("example exists")
	} else {
		fmt.Println("no example!")
	}
}

func TestFiles(t *testing.T) {
	dir := filepath.Join("testdata", "files")

	// Keep this synced with testdata/files directory.
	exp := []string{
		// dot.md is excluded.
		filepath.Join(dir, "jack.txt"),
		filepath.Join(dir, "phryne", "fisher.txt"),
	}

	ret, err := fileutil.Files(dir, ".txt")
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(exp, ret) {
		t.Errorf("expected %s, but returned %s", exp, ret)
	}
}

func ExampleFiles() {
	dir := filepath.Join("testdata/files")
	fmt.Println(fileutil.Files(dir))
}

func TestMkdir(t *testing.T) {
	dir := filepath.Join("testdata", "mkdir")

	if err := fileutil.Mkdir(dir); err != nil {
		t.Error(err)
	}
	t.Cleanup(remove(t, dir))

	if !fileutil.Exists(dir) {
		t.Errorf("%s should exist", dir)
	}
}

func ExampleMkdir() {
	if err := fileutil.Mkdir("example"); err != nil {
		log.Fatal(err)
	}
}
