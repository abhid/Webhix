package web

import "embed"

// Static holds all UI assets embedded into the binary.
// Built frontend assets go into static/ and are served from / and /ui/.
//
//go:embed static
var Static embed.FS
