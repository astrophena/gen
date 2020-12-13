// © 2020 Ilya Mateyko. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE.md file.

package page_test

import (
	"fmt"
	"html/template"
	"path/filepath"
	"reflect"
	"testing"

	"go.astrophena.name/gen/fileutil"
	"go.astrophena.name/gen/internal/page"
)

func testTpl(t *testing.T) *template.Template {
	tpls, err := fileutil.Files("testdata", ".html")
	if err != nil {
		t.Error(err)
	}

	tpl, err := page.ParseTemplates(page.Template(), tpls)
	if err != nil {
		t.Error(err)
	}

	return tpl
}

func TestValidParse(t *testing.T) {
	tpl := testTpl(t)
	f := filepath.Join("testdata", "valid.md")

	p, err := page.Parse(tpl, f)
	if err != nil {
		t.Error(err)
	}

	// Keep this updated with the testdata/valid.md.
	expected := &page.Page{
		URI:      "index.html",
		Content:  fmt.Sprintf("<p>Hello, world!</p>\n"),
		Title:    "Hello, world!",
		MetaTags: make(map[string]string),
		Template: "layout",
	}
	parsed := p

	if !reflect.DeepEqual(expected, parsed) {
		t.Errorf("expected %v, but parsed %v", expected, parsed)
	}
}
