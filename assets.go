// SPDX-License-Identifier: MPL-2.0

// Package assets embeds the template catalog, the engine and the bundled fonts, so the
// generator (cmd/mktex) is one self-contained executable that needs nothing installed.
package assets

import "embed"

// FS holds core/ (engine, fonts, project files) and templates/ (the catalog).
// The "all:" prefix keeps dot-files such as .gitkeep, which preserve empty folders.
//
//go:embed all:core all:templates
var FS embed.FS
