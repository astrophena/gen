// © 2020 Ilya Mateyko. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE.md file.

package page_test

import (
	"html/template"
	"path/filepath"
	"reflect"
	"testing"

	"go.astrophena.name/gen/page"
)

func newTpl(t *testing.T) *template.Template {
	tpl, err := page.ParseTemplates("testdata")
	if err != nil {
		t.Error(err)
	}

	return tpl
}

func TestRewrite(t *testing.T) {
	p, err := page.Parse(newTpl(t), filepath.Join("testdata", "rewrite.md"))
	if err != nil {
		t.Error(err)
	}

	if p.URI != "/hello/index.html" {
		t.Errorf("%s should be /hello/index.html", p.URI)
	}
}

func TestValidParse(t *testing.T) {
	parsed, err := page.Parse(newTpl(t), filepath.Join("testdata", "valid.md"))
	if err != nil {
		t.Error(err)
	}

	// See testdata/valid.md.
	expected := &page.Page{
		URI:      "index.html",
		Content:  "<p>Hello, world!</p>\n",
		Title:    "Hello, world!",
		MetaTags: make(map[string]string),
		Template: "layout",
	}

	if !reflect.DeepEqual(expected, parsed) {
		t.Errorf("expected %v, but parsed %v", expected, parsed)
	}
}
