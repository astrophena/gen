// © 2020 Ilya Mateyko. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE.md file.

// Package frontmatter implements frontmatter-related functions.
package frontmatter // import "go.astrophena.me/gen/pkg/frontmatter"

import (
	"bufio"
	"errors"
	"strings"
)

const delim = "---\n"

// ErrNoFrontmatter is returned when no frontmatter has been detected.
var ErrNoFrontmatter = errors.New("no frontmatter detected")

// Extract extracts frontmatter from a text, returning a frontmatter
// and a content without it.
func Extract(text string) (frontmatter, content string, err error) {
	contains, err := Contains(text)
	if err != nil {
		return "", "", err
	}

	if !contains {
		return "", "", ErrNoFrontmatter
	}

	scanner := bufio.NewScanner(strings.NewReader(text))

	// Skip the first line.
	scanner.Scan()

	var (
		line    string
		reached bool
	)
	for scanner.Scan() {
		if line = scanner.Text() + "\n"; line == delim {
			reached = true
			continue
		}
		if reached {
			content += line
		} else {
			frontmatter += line
		}
	}

	return frontmatter, content, nil
}

// Contains returns true if the text include frontmatter.
func Contains(text string) (contains bool, err error) {
	r := bufio.NewReader(strings.NewReader(text))

	// Why 4 bytes?
	// ---\n
	// That's why.
	b, err := r.Peek(4)
	if err != nil {
		return false, err
	}

	return string(b) == delim, nil
}
