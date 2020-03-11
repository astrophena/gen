// © 2020 Ilya Mateyko. All rights reserved.
// © 2019 Frédéric Guillot. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE.md file.

// Package pack implements asset bundling.
package pack // import "astrophena.me/gen/internal/pack"

import (
	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/css"
	"github.com/tdewolff/minify/v2/html"
)

var minifierFuncs = map[string]minify.MinifierFunc{
	"text/html": html.Minify,
	"text/css":  css.Minify,
}

// Minifier returns a configured minifier.
func Minifier() *minify.M {
	m := minify.New()
	for t, f := range minifierFuncs {
		m.AddFunc(t, f)
	}
	return m
}
