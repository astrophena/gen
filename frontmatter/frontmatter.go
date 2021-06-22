// © 2020 Ilya Mateyko. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE.md file.

// Package frontmatter implements frontmatter-related functions.
package frontmatter

import (
	"bufio"
	"errors"
	"strings"

	"gopkg.in/yaml.v2"
)

const delim = "---\n"

// ErrNotDetected is returned when no frontmatter has been detected.
var ErrNotDetected = errors.New("no frontmatter detected")

// Contains returns true if the supplied text includes frontmatter.
func Contains(text string) (contains bool, err error) {
	r := bufio.NewReader(strings.NewReader(text))

	b, err := r.Peek(len([]byte(delim)))
	if err != nil {
		return false, err
	}

	return string(b) == delim, nil
}

// Extract extracts frontmatter from supplied text, returning
// frontmatter and content.
func Extract(text string) (frontmatter, content string, err error) {
	c, err := Contains(text)
	if err != nil {
		return "", "", err
	}

	if !c {
		return "", "", ErrNotDetected
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

	return frontmatter, content, scanner.Err()
}

// Parse extracts frontmatter from supplied text and unmarshals it
// into obj, returning content without frontmatter and an error.
func Parse(text string, obj interface{}) (content string, err error) {
	fm, c, err := Extract(text)
	if err != nil {
		return "", err
	}

	if err := yaml.Unmarshal([]byte(fm), obj); err != nil {
		return "", err
	}

	return c, nil
}
