// Package ui provides embedded frontend assets.
package ui

import "embed"

// Assets embed build folder
//
//go:embed build/*
var Assets embed.FS
