// © 2020 Ilya Mateyko. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE.md file.

package frontmatter_test

import (
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/astrophena/gen/pkg/frontmatter"
)

var (
	validFile   = filepath.Join("testdata", "valid.md")
	invalidFile = filepath.Join("testdata", "invalid.md")
)

func TestValidExtract(t *testing.T) {
	text, err := ioutil.ReadFile(validFile)
	if err != nil {
		t.Error(err)
	}

	ret1, ret2, err := frontmatter.Extract(string(text))
	if err != nil {
		t.Error(err)
	}

	exp1 := "hello: world\n"

	if ret1 != exp1 {
		t.Errorf("returned %s, but expected %s", ret1, exp1)
	}

	exp2 := "# Hello, world!\n"

	if ret2 != exp2 {
		t.Errorf("returned %s, but expected %s", ret2, exp2)
	}
}

func TestInvalidExtract(t *testing.T) {
	text, err := ioutil.ReadFile(invalidFile)
	if err != nil {
		t.Error(err)
	}

	_, _, err = frontmatter.Extract(string(text))
	if err != frontmatter.ErrNoFrontmatter {
		t.Error("frontmatter shouldn't be detected")
	}
}
