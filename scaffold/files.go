// © 2020 Ilya Mateyko. All rights reserved.
// © 2019 Frédéric Guillot. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE.md file.

// Code generated by go generate; DO NOT EDIT.

package scaffold // import "astrophena.me/gen/scaffold"

var files = map[string][]byte{
	".gitignore":              []byte{0x23, 0x20, 0x49, 0x67, 0x6e, 0x6f, 0x72, 0x65, 0x20, 0x67, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x74, 0x65, 0x64, 0x20, 0x66, 0x69, 0x6c, 0x65, 0x73, 0x2e, 0xa, 0x73, 0x69, 0x74, 0x65, 0xa},
	"README.md":               []byte{0x54, 0x68, 0x69, 0x73, 0x20, 0x69, 0x73, 0x20, 0x61, 0x20, 0x6e, 0x65, 0x77, 0x20, 0x73, 0x69, 0x74, 0x65, 0x20, 0x62, 0x75, 0x69, 0x6c, 0x74, 0x20, 0x77, 0x69, 0x74, 0x68, 0x20, 0x60, 0x67, 0x65, 0x6e, 0x60, 0x2e, 0xa, 0xa, 0x54, 0x6f, 0x20, 0x72, 0x75, 0x6e, 0x20, 0x69, 0x74, 0x20, 0x6c, 0x6f, 0x63, 0x61, 0x6c, 0x6c, 0x79, 0x20, 0x6a, 0x75, 0x73, 0x74, 0x20, 0x72, 0x75, 0x6e, 0x20, 0x60, 0x67, 0x65, 0x6e, 0x20, 0x73, 0x65, 0x72, 0x76, 0x65, 0x60, 0x20, 0x61, 0x6e, 0x64, 0x20, 0x6f, 0x70, 0x65, 0x6e, 0xa, 0x79, 0x6f, 0x75, 0x72, 0x20, 0x62, 0x72, 0x6f, 0x77, 0x73, 0x65, 0x72, 0x20, 0x61, 0x74, 0x20, 0x68, 0x74, 0x74, 0x70, 0x3a, 0x2f, 0x2f, 0x6c, 0x6f, 0x63, 0x61, 0x6c, 0x68, 0x6f, 0x73, 0x74, 0x3a, 0x33, 0x30, 0x30, 0x30, 0x2e, 0xa, 0xa, 0x5b, 0x67, 0x65, 0x6e, 0x5d, 0x3a, 0x20, 0x68, 0x74, 0x74, 0x70, 0x73, 0x3a, 0x2f, 0x2f, 0x61, 0x73, 0x74, 0x72, 0x6f, 0x70, 0x68, 0x65, 0x6e, 0x61, 0x2e, 0x6d, 0x65, 0x2f, 0x67, 0x65, 0x6e, 0xa},
	"assets/sitewide.css":     []byte{0x2f, 0x2a, 0x20, 0x53, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x3a, 0x20, 0x68, 0x74, 0x74, 0x70, 0x3a, 0x2f, 0x2f, 0x62, 0x65, 0x74, 0x74, 0x65, 0x72, 0x6d, 0x6f, 0x74, 0x68, 0x65, 0x72, 0x66, 0x75, 0x63, 0x6b, 0x69, 0x6e, 0x67, 0x77, 0x65, 0x62, 0x73, 0x69, 0x74, 0x65, 0x2e, 0x63, 0x6f, 0x6d, 0x20, 0x2a, 0x2f, 0xa, 0x62, 0x6f, 0x64, 0x79, 0x20, 0x7b, 0xa, 0x20, 0x20, 0x6d, 0x61, 0x72, 0x67, 0x69, 0x6e, 0x3a, 0x20, 0x34, 0x30, 0x70, 0x78, 0x20, 0x61, 0x75, 0x74, 0x6f, 0x3b, 0xa, 0x20, 0x20, 0x6d, 0x61, 0x78, 0x2d, 0x77, 0x69, 0x64, 0x74, 0x68, 0x3a, 0x20, 0x36, 0x35, 0x30, 0x70, 0x78, 0x3b, 0xa, 0x20, 0x20, 0x6c, 0x69, 0x6e, 0x65, 0x2d, 0x68, 0x65, 0x69, 0x67, 0x68, 0x74, 0x3a, 0x20, 0x31, 0x2e, 0x36, 0x3b, 0xa, 0x20, 0x20, 0x66, 0x6f, 0x6e, 0x74, 0x2d, 0x73, 0x69, 0x7a, 0x65, 0x3a, 0x20, 0x31, 0x38, 0x70, 0x78, 0x3b, 0xa, 0x20, 0x20, 0x63, 0x6f, 0x6c, 0x6f, 0x72, 0x3a, 0x20, 0x23, 0x34, 0x34, 0x34, 0x3b, 0xa, 0x20, 0x20, 0x70, 0x61, 0x64, 0x64, 0x69, 0x6e, 0x67, 0x3a, 0x20, 0x30, 0x20, 0x31, 0x30, 0x70, 0x78, 0x3b, 0xa, 0x7d, 0xa, 0xa, 0x68, 0x31, 0x2c, 0x20, 0x68, 0x32, 0x2c, 0x20, 0x68, 0x33, 0x20, 0x7b, 0xa, 0x20, 0x20, 0x6c, 0x69, 0x6e, 0x65, 0x2d, 0x68, 0x65, 0x69, 0x67, 0x68, 0x74, 0x3a, 0x20, 0x31, 0x2e, 0x32, 0x3b, 0xa, 0x7d, 0xa},
	"content/home.html":       []byte{0x74, 0x69, 0x74, 0x6c, 0x65, 0x3a, 0x20, 0x48, 0x65, 0x6c, 0x6c, 0x6f, 0x2c, 0x20, 0x77, 0x6f, 0x72, 0x6c, 0x64, 0x21, 0xa, 0x74, 0x65, 0x6d, 0x70, 0x6c, 0x61, 0x74, 0x65, 0x3a, 0x20, 0x73, 0x69, 0x74, 0x65, 0x77, 0x69, 0x64, 0x65, 0xa, 0x75, 0x72, 0x69, 0x3a, 0x20, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x2e, 0x68, 0x74, 0x6d, 0x6c, 0xa, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x3a, 0x20, 0x4e, 0x65, 0x77, 0x20, 0x73, 0x69, 0x74, 0x65, 0x20, 0x62, 0x75, 0x69, 0x6c, 0x74, 0x20, 0x77, 0x69, 0x74, 0x68, 0x20, 0x67, 0x65, 0x6e, 0x2e, 0xa, 0x2d, 0x2d, 0x2d, 0xa, 0x3c, 0x70, 0x3e, 0xa, 0x20, 0x20, 0x54, 0x68, 0x69, 0x73, 0x20, 0x69, 0x73, 0x20, 0x61, 0x20, 0x6e, 0x65, 0x77, 0x20, 0x73, 0x69, 0x74, 0x65, 0x20, 0x62, 0x75, 0x69, 0x6c, 0x74, 0x20, 0x77, 0x69, 0x74, 0x68, 0xa, 0x20, 0x20, 0x3c, 0x61, 0x20, 0x68, 0x72, 0x65, 0x66, 0x3d, 0x22, 0x68, 0x74, 0x74, 0x70, 0x73, 0x3a, 0x2f, 0x2f, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x61, 0x73, 0x74, 0x72, 0x6f, 0x70, 0x68, 0x65, 0x6e, 0x61, 0x2f, 0x67, 0x65, 0x6e, 0x22, 0x3e, 0x67, 0x65, 0x6e, 0x3c, 0x2f, 0x61, 0x3e, 0x2e, 0xa, 0x3c, 0x2f, 0x70, 0x3e, 0xa},
	"static/robots.txt":       []byte{0x55, 0x73, 0x65, 0x72, 0x2d, 0x61, 0x67, 0x65, 0x6e, 0x74, 0x3a, 0x20, 0x2a, 0xa, 0x41, 0x6c, 0x6c, 0x6f, 0x77, 0x3a, 0x20, 0x2f, 0xa},
	"templates/sitewide.html": []byte{0x7b, 0x7b, 0x20, 0x64, 0x65, 0x66, 0x69, 0x6e, 0x65, 0x20, 0x22, 0x73, 0x69, 0x74, 0x65, 0x77, 0x69, 0x64, 0x65, 0x22, 0x20, 0x7d, 0x7d, 0xa, 0x3c, 0x21, 0x64, 0x6f, 0x63, 0x74, 0x79, 0x70, 0x65, 0x20, 0x68, 0x74, 0x6d, 0x6c, 0x3e, 0xa, 0x3c, 0x68, 0x74, 0x6d, 0x6c, 0x20, 0x6c, 0x61, 0x6e, 0x67, 0x3d, 0x22, 0x65, 0x6e, 0x22, 0x3e, 0xa, 0x20, 0x20, 0x3c, 0x68, 0x65, 0x61, 0x64, 0x3e, 0xa, 0x20, 0x20, 0x20, 0x20, 0x3c, 0x6d, 0x65, 0x74, 0x61, 0x20, 0x63, 0x68, 0x61, 0x72, 0x73, 0x65, 0x74, 0x3d, 0x22, 0x75, 0x74, 0x66, 0x2d, 0x38, 0x22, 0x3e, 0xa, 0x20, 0x20, 0x20, 0x20, 0x3c, 0x6d, 0x65, 0x74, 0x61, 0x20, 0x6e, 0x61, 0x6d, 0x65, 0x3d, 0x22, 0x76, 0x69, 0x65, 0x77, 0x70, 0x6f, 0x72, 0x74, 0x22, 0x20, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x3d, 0x22, 0x77, 0x69, 0x64, 0x74, 0x68, 0x3d, 0x64, 0x65, 0x76, 0x69, 0x63, 0x65, 0x2d, 0x77, 0x69, 0x64, 0x74, 0x68, 0x2c, 0x69, 0x6e, 0x69, 0x74, 0x69, 0x61, 0x6c, 0x2d, 0x73, 0x63, 0x61, 0x6c, 0x65, 0x3d, 0x31, 0x22, 0x3e, 0xa, 0x20, 0x20, 0x20, 0x20, 0x7b, 0x7b, 0x20, 0x72, 0x61, 0x6e, 0x67, 0x65, 0x20, 0x24, 0x6e, 0x61, 0x6d, 0x65, 0x2c, 0x20, 0x24, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x20, 0x3a, 0x3d, 0x20, 0x2e, 0x4d, 0x65, 0x74, 0x61, 0x54, 0x61, 0x67, 0x73, 0x20, 0x7d, 0x7d, 0xa, 0x20, 0x20, 0x20, 0x20, 0x3c, 0x6d, 0x65, 0x74, 0x61, 0x20, 0x6e, 0x61, 0x6d, 0x65, 0x3d, 0x22, 0x7b, 0x7b, 0x20, 0x24, 0x6e, 0x61, 0x6d, 0x65, 0x20, 0x7d, 0x7d, 0x22, 0x20, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x3d, 0x22, 0x7b, 0x7b, 0x20, 0x24, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x20, 0x7d, 0x7d, 0x22, 0x3e, 0xa, 0x20, 0x20, 0x20, 0x20, 0x7b, 0x7b, 0x20, 0x65, 0x6e, 0x64, 0x20, 0x7d, 0x7d, 0xa, 0x20, 0x20, 0x20, 0x20, 0x3c, 0x73, 0x74, 0x79, 0x6c, 0x65, 0x3e, 0x7b, 0x7b, 0x20, 0x63, 0x73, 0x73, 0x20, 0x2e, 0x43, 0x53, 0x53, 0x20, 0x7d, 0x7d, 0x3c, 0x2f, 0x73, 0x74, 0x79, 0x6c, 0x65, 0x3e, 0xa, 0x20, 0x20, 0x20, 0x20, 0x3c, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x3e, 0x7b, 0x7b, 0x20, 0x2e, 0x54, 0x69, 0x74, 0x6c, 0x65, 0x20, 0x7d, 0x7d, 0x3c, 0x2f, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x3e, 0xa, 0x20, 0x20, 0x3c, 0x2f, 0x68, 0x65, 0x61, 0x64, 0x3e, 0xa, 0x20, 0x20, 0x3c, 0x62, 0x6f, 0x64, 0x79, 0x3e, 0xa, 0x20, 0x20, 0x20, 0x20, 0x3c, 0x68, 0x31, 0x3e, 0x7b, 0x7b, 0x20, 0x2e, 0x54, 0x69, 0x74, 0x6c, 0x65, 0x20, 0x7d, 0x7d, 0x3c, 0x2f, 0x68, 0x31, 0x3e, 0xa, 0x20, 0x20, 0x20, 0x20, 0x7b, 0x7b, 0x20, 0x68, 0x74, 0x6d, 0x6c, 0x20, 0x2e, 0x42, 0x6f, 0x64, 0x79, 0x20, 0x7d, 0x7d, 0xa, 0x20, 0x20, 0x3c, 0x2f, 0x62, 0x6f, 0x64, 0x79, 0x3e, 0xa, 0x3c, 0x2f, 0x68, 0x74, 0x6d, 0x6c, 0x3e, 0xa, 0x7b, 0x7b, 0x20, 0x65, 0x6e, 0x64, 0x20, 0x7d, 0x7d, 0xa},
}
