// © 2020 Ilya Mateyko. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE.md file.

package frontmatter_test

import (
	"os"
	"path/filepath"
	"testing"

	"go.astrophena.name/gen/frontmatter"
)

var (
	validFile   = filepath.Join("testdata", "valid.md")
	invalidFile = filepath.Join("testdata", "invalid.md")
)

func TestValidExtract(t *testing.T) {
	text, err := os.ReadFile(validFile)
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
	text, err := os.ReadFile(invalidFile)
	if err != nil {
		t.Error(err)
	}

	_, _, err = frontmatter.Extract(string(text))
	if err != frontmatter.ErrNotDetected {
		t.Error("frontmatter shouldn't be detected")
	}
}
