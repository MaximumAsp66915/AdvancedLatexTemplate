# Architecture

## Goal

A catalog of professional LaTeX templates in **English and Persian (Farsi)** for Windows, macOS and Linux:
homework, lab report, project report, formal report (thesis), essay and paper. Every template exists in both
languages with the same structure, and every template can use any of 11 color palettes.

## Repository map

```
.
├── README.md  LICENSE  CONTRIBUTING.md  AGENTS.md       the public face (AGENTS.md starts this manual)
├── templates/<category>/                                 the catalog
│   ├── template.conf         label, description, default palette/code theme, title page, prompts
│   └── en/   fa/             exactly what the generator copies into a project (same file list in both)
├── core/
│   ├── atl/                  the engine (atl.tex loader + atl-*.tex modules), copied into every project
│   ├── fonts/                bundled Persian fonts + their SIL OFL license texts
│   └── project/              files every project gets: latexmkrc, config-options.{en,fa}.tex,
│                             mktex.defaults.example
├── cmd/mktex/                the generator (Go): main.go, main_test.go
├── go.mod  assets.go         Go module; assets.go embeds core/ and templates/ into the executable
├── tools/                    maintainer scripts: check.sh, gallery.sh, audit-colors.py, build-release.sh
├── docs/                     USER-MANUAL.md, ai/ (this manual), gallery/ (preview images)
└── .github/                  issue and pull-request templates
```

Local-only folders that a clone does not have (all git-ignored): `legacy/`, `fa/`, `dist/`, `CLAUDE.md` files.

## The generator (`cmd/mktex`)

One Go program, no dependencies beyond the standard library. `assets.go` embeds `core/` and `templates/`, so a built
executable is self-contained (about 4.7 MB). Releases carry six builds: Windows (amd64, arm64), macOS (arm64, amd64),
Linux (amd64, arm64). The same source produces all of them, which is why a Linux test run is meaningful for Windows too;
the Windows build was also run under Wine and produced byte-identical projects.

- Flags: `-l -c -n -p --listing --logo -s KEY=VAL -o -y --no-defaults --no-open --list --defaults --init-defaults --root --version`.
- `--root DIR` reads the catalog from a checkout instead of the embedded copy: this is how `tools/check.sh` tests
  edited templates without rebuilding.
- Personal defaults (name, ID, university, logo) are read from the user's config folder
  (`~/.config/mktex/defaults`, `%AppData%\mktex\defaults`) or `MKTEX_DEFAULTS` / `--defaults`; `mktex --init-defaults`
  creates the file from `core/project/mktex.defaults.example`. They are never inside the repository.
- Double-clicked on Windows (no arguments) it keeps the window open at the end.
- Template metadata comes from `templates/<category>/template.conf` (`KEY="value"` lines); palettes are read from
  `core/atl/atl-palettes.tex`. Adding a template or palette therefore needs no Go change, except new prompt labels
  (`promptLabels` in `main.go`) and new code themes (`listings` in `main.go`).
- `go test ./...` generates every template in both languages and checks the result.

## How a project is made

The generator copies `templates/<category>/<lang>/` to `<name>/`, then adds:

- `latex/atl/` : the engine (`core/atl/`)
- `latex/fonts/` : the Persian fonts and licenses (Persian projects only)
- `latex/.latexmkrc` : from `core/project/latexmkrc` (XeLaTeX via `$pdf_mode = 5`)
- `latex/config.tex` : overwritten with `\atlset{...}` from the prompts/flags, plus the options reference from
  `core/project/config-options.<lang>.tex` appended as comments
- `.vscode/settings.json` : an XeLaTeX recipe for LaTeX Workshop
- `latex/fig/logo.<ext>` : only for cover-page templates, if a logo was given

Build with `latexmk -xelatex main.tex` inside `latex/`.

Every template `main.tex` loads, in this order: `atl/atl-keys` (keys) -> template defaults (`\atlset{...}`) ->
`config` (the user's values) -> `atl/atl` (the engine) -> `preamble/local` (the user's macros).
So precedence is **template defaults < config.tex < preamble/local.tex**.

## Decisions (change only deliberately, and update this manual)

- **Self-contained projects**, not a shared preamble: the engine and fonts are copied in. Updating the repository
  never changes existing projects; to upgrade one, copy `core/atl/*` over its `latex/atl/`.
- **XeLaTeX for both languages.** XePersian requires it; one toolchain, fonts through fontspec.
- **One generator, in Go, distributed as executables** (GitHub Releases), instead of shell scripts: it runs on Windows
  without bash, needs nothing installed, and is tested once on Linux for all platforms. Executables are never committed.
- **Persian fonts**: Vazirmatn and HM FElmi, both SIL OFL 1.1, license texts kept next to them. No proprietary font is
  required or bundled (B Nazanin, B Lotus, Times New Roman, ...).
- **Palettes are data**: one `\atldefpalette` line each (no conditional chains). **Code themes are data** too.
- **Math stays Computer Modern in both languages**; only text and digits use the Persian/Latin text fonts.
- **License: MPL-2.0** (file-level copyleft). Free to use, including commercially; a distributed modified copy of any file
  here must be published under the same license, while the user's own documents stay theirs. Chosen so that improvements
  flow back to the project. See `LICENSE` and `CONTRIBUTING.md`. New engine, tool and Go files carry an SPDX notice.
- **No personal data in the repository.** Personal defaults live in the user's config folder; templates and examples are neutral.
- **Naming**: lowercase letters, digits and hyphens for files and folders (`lab-report`, `01-introduction.tex`); the
  sibling folder for course material is `preprovided`; no underscores or spaces; American spelling in text.
- **History**: the public history was reset to a single first commit. Do not reintroduce removed material (a university
  logo, personal defaults, the old English system).

## Rules that come from the TeX search path

If this repository is cloned into `~/texmf/tex/latex` (TeX's `TEXMFHOME`), TeX searches it recursively:

- **Never compile inside the repository.** A stray `main.toc` or `main.aux` anywhere in the tree is found by every
  project's first compile and breaks it (`main.toc:5: Undefined control sequence` happened for real). Use
  `tools/check.sh` or the generator in a scratch directory.
- **Shared filenames must be unique**; engine files are all prefixed `atl-`.
- A `.nosearch` marker file does **not** hide a directory from the search in TeX Live 2023 (tested).

## Open items

1. **Code signing.** The executables are not signed or notarized: Windows SmartScreen and macOS Gatekeeper show a warning
   on first run (the README explains how to continue). Signing needs a certificate the owner would have to obtain.
2. **Tested environments.** Everything is verified on Linux with TeX Live 2023; the Windows build was run under Wine
   (project generation only). The generated projects were never compiled on Windows or macOS TeX distributions or on Overleaf.
3. **Persian behavior on long documents** (real theses) was tested only on sample content.
4. **Persian-language bibliography entries** were not tried (references are typeset left-to-right in Latin).
5. **Footnote numbering restarting on every page** (a requirement of one Persian formal-report convention) is not implemented.
6. **Meaning parity** between English and Persian text is a manual review; `check.sh` checks structure and builds only.
