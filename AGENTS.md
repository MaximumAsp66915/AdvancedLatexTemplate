# AGENTS.md

Instructions for AI coding agents working in this repository (Claude Code, OpenAI Codex / GPT agents,
Gemini CLI, Cursor, Copilot, Aider or OpenCode on any model including Groq-hosted ones, or a plain chat
session). Humans: see the [README](README.md) and the [user manual](docs/USER-MANUAL.md).

## What this repository is

A catalog of professional LaTeX templates in **English and Persian (Farsi)**: six categories
(`homework`, `lab-report`, `project-report`, `formal-report`, `essay`, `paper`), each in both languages with
the same file structure, 11 selectable color palettes, built with **XeLaTeX**. A generator (`mktex`, a Go program in
`cmd/mktex/` that embeds `templates/` and `core/`) copies a template plus the shared engine (`core/atl/`) and, for Persian,
the bundled fonts into a self-contained project. It is published as executables for Windows, macOS and Linux.

```
templates/   the catalog (templates/<name>/{en,fa})      cmd/mktex/ + go.mod + assets.go   the generator
core/        engine (atl/), fonts, project files          tools/   maintainer scripts (check, gallery, release)
docs/        user manual, this agent manual, previews     .github/ issue and pull-request templates
```

## Read this next

**[docs/ai/README.md](docs/ai/README.md)** is the entry point of the agent manual. It routes you to the one
topic file you need ([architecture](docs/ai/architecture.md), [engine](docs/ai/engine.md),
[templates](docs/ai/templates.md), [Persian](docs/ai/persian.md), [workflow](docs/ai/workflow.md)), so you
do not have to load everything.

## Rules that prevent the most expensive mistakes

1. **Never compile inside the repository.** A stray `main.toc` or `main.aux` in the tree breaks every
   project's first build. Test with `tools/check.sh`, or generate a project with `mktex --root . -o <scratch folder>`.
2. **Parity is a hard rule.** Every template exists in `en` and `fa` with identical file names. After any change
   under `core/`, `templates/` or `cmd/` run `tools/check.sh`; it must print `PASS`.
3. **XeLaTeX only.** Do not add pdfLaTeX-only packages or `inputenc`/`fontenc` setup.
4. **Colors come from palette roles** (`atl-primary`, `atl-accent`, ...). Never write raw colors in templates.
5. **Every user-visible label is a string** defined with `\atldefstr` in **both** `atl-lang-en.tex` and
   `atl-lang-fa.tex`. Never hard-code English or Persian words in shared engine files.
6. **No `&` or `\\` inside `\if...\fi`.** Use the `\atlifset` / `\atlifeq` / `\atlifon` macros, which are safe in tables.
7. **Persian text rules** (details in [persian.md](docs/ai/persian.md)): Latin words and codes in `\lr{...}`,
   decimals in math or `\num{}`, ZWNJ in compound words, diagrams inside `\begin{LTR}`, headings with
   `\persiansffamily` (never `\sffamily`), ezafe after a final heh as `ه‌ی`.
8. **Names**: folders and files use lowercase letters, digits and hyphens (`lab-report`, `01-introduction.tex`);
   the one folder that merges words is `preprovided`. No underscores or spaces in names. American spelling in text
   (color, license, catalog).
9. **The repository is public and licensed MPL-2.0.** Never commit personal data (names, IDs, emails), personal defaults,
   university logos, proprietary fonts (B Nazanin, Times New Roman, ...), third-party code without a license check,
   executables, or build output. New engine, tool and Go files start with `SPDX-License-Identifier: MPL-2.0`.
   See [CONTRIBUTING.md](CONTRIBUTING.md).
10. **Look at the result.** A clean log does not prove a good layout: render pages with `pdftoppm` and check them,
    and say plainly what you did not verify.
11. **Commit or push only when the owner asks.** Do not rewrite history or force-push unless explicitly told to.

## Commands

```bash
tools/check.sh                 # generator tests, color audit, parity, build every template (en+fa), code-color test (~1.5 min)
tools/check.sh --no-build      # fast checks only
tools/check.sh --palettes      # also build with every palette
tools/check.sh paper           # one template
go run ./cmd/mktex --root . --list                       # templates, palettes, code themes
go run ./cmd/mktex --root . -l fa -c homework -n hw1 -p sunset -s title="..." -y --no-defaults --no-open -o /tmp/x
tools/gallery.sh               # regenerate docs/gallery previews
tools/build-release.sh v1.0.0  # build the Windows/macOS/Linux executables into dist/
```

Build a generated project with `latexmk -xelatex main.tex` inside its `latex/` folder.
