package web

import "embed"

// Static holds all UI assets embedded into the binary.
// Templates and JS/CSS go into static/ — served at /ui/.
//
//go:embed static
var Static embed.FS
