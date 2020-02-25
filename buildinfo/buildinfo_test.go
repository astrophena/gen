// © 2020 Ilya Mateyko. All rights reserved.
// © 2019 Frédéric Guillot. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE.md file.

package buildinfo // import "astrophena.me/gen/buildinfo"

import "testing"

func TestTplFunc(t *testing.T) {
	Func := TplFunc()

	expected := Func()
	returned := Version

	if expected != returned {
		t.Errorf("expected %s, but returned %s", expected, returned)
	}
}
