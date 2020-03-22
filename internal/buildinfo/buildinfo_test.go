// © 2020 Ilya Mateyko. All rights reserved.
// © 2019 Frédéric Guillot. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE.md file.

package buildinfo_test

import (
	"fmt"
	"testing"

	"go.astrophena.me/gen/internal/buildinfo"
)

func TestTplFunc(t *testing.T) {
	exp := buildinfo.TplFunc()()
	ret := buildinfo.Version

	if exp != ret {
		t.Errorf("expected %s, but returned %s", exp, ret)
	}
}

func ExampleTplFunc() {
	f := buildinfo.TplFunc()
	fmt.Println(f())
	// Output: devel
}
