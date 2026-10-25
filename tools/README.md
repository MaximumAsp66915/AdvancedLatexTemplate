# tools/ — for maintainers

You do **not** need anything in this folder to use the templates. These scripts keep the repository healthy and publish it.

| Script | What it does | Needs |
|---|---|---|
| `check.sh` | The quality gate. Checks that every template exists in English and Persian with identical files, tests the generator, audits every color pair, builds all templates, and tests code colors across a page break. Ends with `PASS` or `FAILED`. `--no-build` runs only the fast checks, `--palettes` also builds with every palette, `check.sh homework paper` limits it to some templates. | bash, Go, XeLaTeX, latexmk; optional: python3, mutool |
| `audit-colors.py` | Contrast check for every palette and code theme (part of `check.sh`, also runnable alone). | python3 |
| `gallery.sh` | Rebuilds the preview images in `docs/gallery/` after a visual change. | bash, Go, XeLaTeX, pdftoppm |
| `build-release.sh VERSION` | Builds the `mktex` executables for Windows, macOS and Linux into `dist/` (git-ignored) and writes `SHA256SUMS`. Prints the `gh release create` command to publish them. | bash, Go |

The generator itself is the Go program in `cmd/mktex/` (tests: `go test ./...`).
