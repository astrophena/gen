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
