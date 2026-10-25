# Workflow: checking, verifying, publishing

## Commands

| Command | What it does |
|---|---|
| `tools/check.sh` | in order: parity of every template (en vs fa file lists, `config.tex` options block); generator tests (`go vet`, `go test`); color audit; a build of every variant; a dark-theme code test across a page break (English and Persian). About 1.5 minutes. Must end with `PASS` |
| `python3 tools/audit-colors.py` | only the color audit (palettes and code themes; needs just Python 3) |
| `tools/check.sh --no-build` | parity only (instant) |
| `tools/check.sh homework paper` | only the named templates |
| `tools/check.sh --palettes` | additionally builds one template with every palette |
| `go run ./cmd/mktex --root . --list` | templates, palettes and code themes (`--root .` makes the generator read this checkout instead of its embedded copy) |
| `go run ./cmd/mktex --root . -l fa -c lab-report -n lab3 -p ocean -s title="..." -y --no-defaults --no-open -o /tmp/x` | generate a project non-interactively (always use `-o` outside the repository and `--no-defaults` for reproducible output) |
| `go test ./...` | generator tests: every template in both languages generates a complete project |
| `tools/build-release.sh v1.0.0` | build the six executables into `dist/` (git-ignored) and `SHA256SUMS` |
| `tools/gallery.sh` | regenerate `docs/gallery/*.png` after a visual change |
| `latexmk -xelatex main.tex` (inside a generated `latex/`) | build one project |

Requirements: [Go](https://go.dev/dl/), a TeX distribution with XeLaTeX and `latexmk`, `bash`, and for the optional checks `python3` (color audit), `mutool` (code-color test) and `pdftoppm` / `pdftotext` (poppler). A check whose tool is missing is reported as skipped, never as passed.

## Verification protocol

A change is not done until each applicable step below was run, in **both** languages when `core/` or a template changed.

1. `tools/check.sh` prints `PASS` (0 errors, 0 missing glyphs, 0 undefined references, no `[?]` / `??` in the PDF text).
2. **Look at the output.** Warnings alone do not show layout bugs. Generate a project in a scratch folder and render it:
   `pdftoppm -r 60 -png main.pdf p` then open the PNGs. Check title page, headers, boxes, tables, figures, the last page.
3. For Persian, check word order **by coordinates**, not by eye: `pdftotext -bbox main.pdf out.html` (first logical word of an
   RTL phrase has the larger `xMin`), and `pdftotext main.pdf - | grep -nE '\[\?\]|\?\?'` must print nothing.
4. After a palette change run `tools/check.sh --palettes` and look at a page with a colored box in the palette.
5. Report honestly: say which steps you ran, which you could not run (for example Windows, macOS, Overleaf), and anything
   you saw but did not fix. Never claim a build passes without having run it.

Typical false alarms: overfull boxes under 1 pt; the `hypcap` warning from `\captionof` is already silenced.

## Public repository hygiene

The repository is **public**. Before every commit read `git status` and `git diff --cached` and make sure none of this is staged:

- Personal data: real names, student IDs, e-mail addresses, absolute home paths. Personal defaults live in the user's config folder (`~/.config/mktex/defaults`), outside the repository.
- University or company logos, other people's documents, third-party code or text without a license check.
- Proprietary fonts (B Nazanin, B Lotus, B Zar, Times New Roman, ...). Only the OFL fonts in `core/fonts/` are bundled.
- Build output: `*.pdf *.aux *.log *.toc *.bbl ...` (git-ignored). The only tracked images are in `docs/gallery/`.
- **Executables**: `dist/` is git-ignored; binaries are published as GitHub release assets, never committed.
- Local-only folders and files (all git-ignored): `legacy/`, `fa/`, `dist/`, `CLAUDE.md`, `.claude/`, editor folders.

**License.** The repository is under the Mozilla Public License 2.0 (`LICENSE`). Contributions are accepted under the same license. Do not add code, text, images or fonts whose license is incompatible
with it or unknown; bundled fonts keep their own license file next to them. A university or company logo is never acceptable.

Quick audit of the working tree (should print nothing):

```bash
git grep -n -i -E '/home/[a-z]+|@gmail\.com|@[a-z]+\.edu' -- . ':!core/fonts' ':!docs/ai'
git ls-files | grep -i -E '\.pdf$|\.aux$|\.exe$|^dist/|^legacy/|^fa/|CLAUDE|\.claude'
```

Rules for history: the public history was reset to one first commit on purpose. Do **not** rewrite published history, do **not** force-push, and do not delete branches, unless the owner
explicitly asks. If you find something sensitive that is already public, tell the owner what and where; rewriting history is their decision.

Commits: commit or push only when asked. Look at `git log --oneline` and match its style (short imperative subject, then a body
that says what changed and what was verified). Keep one logical change per commit. Do not commit generated gallery images
unless the visuals changed.

## Recipes

| Task | Steps |
|---|---|
| Add a palette | one `\atldefpalette` line in `core/atl/atl-palettes.tex`; `tools/check.sh --palettes`; update the palette lists (README, `docs/USER-MANUAL.md`, the options reference) |
| Add a code theme | one `\atldeflisttheme` line in `atl-code.tex`; add the name to `listings` in `cmd/mktex/main.go`; update the options reference |
| Add a label | `\atldefstr{key}{...}` in `atl-lang-en.tex` **and** `atl-lang-fa.tex` |
| Add a configuration key | read it in the module; `\atldefault` if needed; update `config-options.{en,fa}.tex`, every template `config.tex` block, `docs/USER-MANUAL.md` section 5, [engine.md](engine.md) |
| Add a template | see [templates.md](templates.md#adding-a-new-category) |
| Add a language | new `atl-lang-<code>.tex`, loader branch in `atl.tex`, `<code>/` in every template, generator prompts, this manual |
| Change a Persian rule | [persian.md](persian.md), with the command you used to verify it |

## Troubleshooting

| Symptom | Likely cause |
|---|---|
| `main.toc:5: Undefined control sequence` on a first build | a stale `main.toc`/`main.aux` inside the repository tree is being found; delete it and build elsewhere |
| Persian text invisible in headings or labels | `\sffamily` used instead of `\persiansffamily` / `\atl@sf` |
| Code in a Persian document shows `h0.:e`, `14.3` | listing lost the `\atl@ltr` switch (or it was placed after the font) |
| Persian number shows reversed (`۳۸.۱`) or `۱R` | decimal typed as text, or Latin token without `\lr{}` |
| `[key]` printed in the document | label missing in one of the language layers |
| `[?]` or `??` | undefined `\cite` / `\ref`: wrong key, or the build did not repeat passes (use `latexmk`) |
| "Incomplete \ifodd" or alignment errors in titles/tables | `&` or `\\` inside an `\if` conditional (see [engine.md](engine.md) rule 1) |
| `Package fontspec Error: font ... cannot be found` | a template names a system font; use the bundled / TeX Gyre fonts |
| `check.sh` says file lists differ | a file exists in only one language; add its counterpart |
| `check.sh` says `config.tex lacks the current options reference` | re-append the block from `core/project/config-options.<lang>.tex` |
| Build fails with a missing package | install it (`tlmgr install <name>`) or use a full TeX distribution |

## Known gaps (do not claim these work)

Windows-native generator, runs on Windows / macOS / Overleaf, footnote restart per page in Persian, Persian-language bibliography entries,
Details in [architecture.md](architecture.md#open-items).
