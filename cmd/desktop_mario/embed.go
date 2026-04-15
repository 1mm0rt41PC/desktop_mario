package main

import "embed"

// assetsFS embeds the entire assets directory into the binary so that
// sprite PNG files are available without external files on disk.
//
//go:embed all:assets
var assetsFS embed.FS
