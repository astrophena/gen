// © 2020 Ilya Mateyko. All rights reserved.
// © 2019 Frédéric Guillot. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE.md file.

// Package buildinfo contains build information.
package buildinfo // import "astrophena.me/gen/buildinfo"

import (
	"runtime/debug"
	"strings"
)

const defaultVersion = "devel"

// Version of gen (generated by "make" or "goreleaser").
var Version = defaultVersion

// TplFunc returns a template function that returns gen's version.
func TplFunc() func() string {
	return func() string {
		return Version
	}
}

func init() {
	if Version == defaultVersion {
		bi, ok := debug.ReadBuildInfo()
		if ok {
			Version = strings.TrimPrefix(bi.Main.Version, "v")
		}
	}
}
