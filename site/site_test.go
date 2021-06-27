package site_test

import (
	"os"
	"testing"

	"go.astrophena.name/gen/scaffold"
	"go.astrophena.name/gen/site"
)

func TestBuild(t *testing.T) {
	src, err := os.MkdirTemp("", "gen-site-test-src")
	if err != nil {
		t.Fatalf("Failed to create a temporary directory: %v", err)
	}
	defer os.RemoveAll(src)
	dst, err := os.MkdirTemp("", "gen-site-test-dst")
	if err != nil {
		t.Fatalf("Failed to create a temporary directory: %v", err)
	}
	defer os.RemoveAll(dst)

	if err := scaffold.Create(src); err != nil {
		t.Fatalf("Failed to generate a new site: %v", err)
	}

	s, err := site.New(src, dst, true, false)
	if err != nil {
		t.Fatalf("Failed to initialize a new site: %v", err)
	}

	if err := s.Build(); err != nil {
		t.Fatalf("Failed to build a site: %v", err)
	}
}
