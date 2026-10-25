# AdvancedLatexTemplate: user manual

Professional LaTeX templates in **English and Persian (Farsi)** for homework, lab reports, project
reports, formal reports (theses), essays and papers. This manual explains how to create a project, what
you can configure, how to write content with the template's building blocks, and how to fix common
problems. For a one-page overview see the [README](../README.md).

Contents: [1 Requirements](#1-requirements) · [2 Creating a project](#2-creating-a-project) ·
[3 Project layout](#3-project-layout) · [4 Building](#4-building) · [5 Configuration reference](#5-configuration-reference) ·
[6 The templates](#6-the-templates) · [7 Writing content](#7-writing-content) · [8 Colors and code themes](#8-colors-and-code-themes) ·
[9 Persian documents](#9-persian-documents) · [10 Windows, macOS, Overleaf](#10-windows-macos-overleaf) ·
[11 Troubleshooting](#11-troubleshooting) · [12 Upgrading an old project](#12-upgrading-an-old-project) · [13 License and contributing](#13-license-and-contributing)

---

## 1. Requirements

- A TeX distribution with **XeLaTeX** and `latexmk`: [TeX Live](https://tug.org/texlive/) (Linux, Windows), [MacTeX](https://tug.org/mactex/) (macOS)
  or [MiKTeX](https://miktex.org/) (Windows; `latexmk` needs Perl there). Both languages are built with XeLaTeX. Or use Overleaf and install nothing.
- The generator **`mktex`**: one executable file for your system, from the
  [latest release](https://github.com/MaximumAsp66915/AdvancedLatexTemplate/releases/latest) (`mktex-windows-amd64.exe`, `mktex-macos-arm64`,
  `mktex-linux-amd64`, ...). Nothing else is needed: no bash, Python or Git. It is not code-signed; see [section 10](#10-windows-macos-overleaf) for the first-run warning.
- `pdftotext` (poppler) only if you want the `[?]` checks described in [section 11](#11-troubleshooting).

## 2. Creating a project

Run the generator from the folder where the project should be created (on Windows you can simply double-click
`mktex-windows-amd64.exe`; the project is then created next to it):

```bash
mktex                                          # interactive (use the name of the file you downloaded)
```

It asks for the language, the template, a folder name, the palette, the code theme and the title-page
fields of that template (title, author, course, ...). Press Enter to accept the default shown in brackets.
Non-interactive form:

```bash
mktex -l fa -c homework -n hw2 -p ocean -s title="تمرین ۲" -s author="نام" -s course="طراحی پایگاه داده‌ها" -y
mktex --list                                   # templates, palettes, code themes
```

| Option | Meaning |
|--------|---------|
| `-l`, `--lang en\|fa` | language |
| `-c`, `--category NAME` | `homework`, `lab-report`, `project-report`, `formal-report`, `essay`, `paper` |
| `-n`, `--name FOLDER` | folder to create (it must not exist yet) |
| `-p`, `--palette NAME` | color palette (default: the template's own) |
| `--listing NAME` | code theme: `darcula`, `dark`, `bright`, `github` |
| `--logo FILE` | logo copied to `latex/fig/` (cover pages only) |
| `-s`, `--set KEY=VALUE` | any configuration key, repeatable (see [section 5](#5-configuration-reference)) |
| `-o`, `--out DIR` | parent folder of the new project (default: current folder) |
| `-y`, `--yes` | do not ask; use defaults for everything not given |
| `--no-open` | do not offer to open VS Code |
| `--no-defaults` | ignore your personal defaults file |
| `--list` | show templates, palettes and code themes |
| `--init-defaults` | create your personal defaults file and print where it is |
| `--version`, `-h` | version, help |

### Personal defaults

Run `mktex --init-defaults`. It creates a small file in your user configuration folder
(`~/.config/mktex/defaults` on Linux and macOS, `%AppData%\mktex\defaults` on Windows) and prints its path. Fill in your name,
student ID, university, department and the path of a logo once; the generator then proposes them every time. A value can be
language-specific: `DEFAULT_AUTHOR_EN="Jane Doe"`, `DEFAULT_AUTHOR_FA="جین دو"`. Any question the generator asks
can be pre-answered this way (`DEFAULT_COURSE`, `DEFAULT_INSTRUCTOR`, `DEFAULT_TERM`, ...). The file lives outside any
repository, so your personal data is never committed by accident.

## 3. Project layout

```
<project>/
├── latex/                  compile root
│   ├── main.tex            entry point: template defaults at the top, then the document body
│   ├── config.tex          YOUR values (title, name, palette, ...), with an options reference as comments
│   ├── sections/ or chapters/, appendix/, frontmatter/   content files, one per section or chapter
│   ├── fig/                figures and the logo (referenced by file name only)
│   ├── refs.bib            bibliography (templates that cite)
│   ├── preamble/packages.tex   extra \usepackage lines (loaded before the language layer)
│   ├── preamble/local.tex      your own macros and settings (loaded after the engine)
│   ├── atl/                the engine; copied in, normally never edited
│   └── fonts/              Persian fonts (Persian projects only)
├── preprovided/, simulation/, data/      sibling folders some templates give you
└── .vscode/settings.json   builds with XeLaTeX in VS Code (LaTeX Workshop)
```

Projects are **self-contained**: they do not need this repository or `TEXMFHOME`, so you can zip a project,
share it, or upload `latex/` to Overleaf. Edit `config.tex`, `preamble/local.tex` and your content files; leave `atl/` alone.

## 4. Building

```bash
cd <project>/latex
latexmk -xelatex main.tex          # builds main.pdf (runs BibTeX and repeats passes automatically)
latexmk -c                         # remove auxiliary files, keep the PDF
```

In VS Code with the *LaTeX Workshop* extension the generated `.vscode/settings.json` already selects the
XeLaTeX recipe. On Overleaf choose **XeLaTeX** as the compiler. Never build inside the template repository
itself (stray `main.toc`/`main.aux` files there break every project's first compile); always build in a project folder.

## 5. Configuration reference

Everything is set with `\atlset{key = value, ...}`. The same list is written as comments at the bottom of every
`config.tex`. Rules:

- Values with spaces or commas need braces: `title = {My Report, Part 1}`; plain words and numbers do not: `palette = ocean`, `linespread = 1.3`.
- A later setting overrides an earlier one: **template defaults** (top of `main.tex`) < **`config.tex`** < **`preamble/local.tex`**.
- Paths are relative to `latex/`: `logo = fig/logo.png`.
- Empty or unset data keys are simply skipped on the title page.

### Look

| Key | Values | Default | Notes |
|-----|--------|---------|-------|
| `palette` | `blue red gray cyan magenta yellow green sunset ocean berry neon` | template's own | see [section 8](#8-colors-and-code-themes) |
| `listing` | `darcula dark bright github` | `darcula` | code block theme |
| `enfont` | `termes pagella lm` | `termes` | English body font (TeX Gyre Termes / Pagella, Latin Modern) |
| `fafont` | `vazirmatn felmi` | `vazirmatn` | Persian body font |
| `fascale` | number, e.g. `0.95` | 0.95 (Vazirmatn), 1.1 (FElmi) | Persian font scale |
| `headfont` | `sans serif` | `sans` | headings and labels |
| `headsize` | `normal compact` | `normal` | `compact` is meant for two-column papers |

### Layout

| Key | Values | Default | Notes |
|-----|--------|---------|-------|
| `margins` | `normal narrow formal wide` | `normal` | `formal`: 3 cm on the binding side (right in Persian), 2.5 cm elsewhere |
| `linespread` | decimal number | `1.15` | no unit; `1.5` is one-and-a-half spacing |
| `parstyle` | `block indent` | `block` | spaced paragraphs / first-line indent (0.71 cm) |
| `numbering` | `flat section` | `flat` | `section`: equations, figures and tables numbered per section (2.1, 2.2). Report-class templates always number per chapter |
| `sectionbreak` | `on off` | `off` | `on`: every `\section` starts a new page |

### Title page

| Key | Values | Notes |
|-----|--------|-------|
| `titlepage` | `block cover classic none` | banner on page 1 / full cover page / centered journal style / nothing |
| `infofields` | `{author,studentid,course,...}` | keys shown in the information grid, comma-separated without spaces. Allowed: `author studentid course coursecode instructor supervisor advisor team partners experiment labsection term date duedate submissiondate performed university department` |
| `doctype` | text | replaces the automatic label («Homework», «Report», ...) |
| `logo` | file path | cover pages only; `png`, `jpg` or `pdf` |

### Document data

`title subtitle author studentid course coursecode instructor supervisor advisor team partners experiment labsection
term date duedate submissiondate performed university department affiliation keywords`: plain text.
`date = {\today}` (Gregorian in English, Jalali in Persian) or a literal date. `keywords` appears under the abstract
(paper, essay); `affiliation` on the classic title page (`\\` breaks the line).

### Header and references

| Key | Values | Default | Notes |
|-----|--------|---------|-------|
| `header` | `on off` | `on` | running header |
| `headerstart` | text | course, else title | start side of the header (left in English, right in Persian) |
| `headerend` | text | current section | end side |
| `bibstyle` | any installed BibTeX style | `IEEEtran` | e.g. `plain`, `ieeetr` |

`language` (`en`/`fa`) and `category` are fixed by the template; do not change them in an existing project.

## 6. The templates

| Template | Class | Title page | Default palette | Typical use |
|----------|-------|-----------|-----------------|-------------|
| `homework` | article | block | blue | problem sets: `problem`/`solution` boxes, code, results |
| `lab-report` | article | block | cyan | objective, theory, apparatus, procedure, data and plots, analysis, conclusion |
| `project-report` | article | cover | sunset | cover, abstract, contents, requirements, design, implementation, results |
| `formal-report` | report | cover | gray | thesis, internship or final report: chapters, front matter, lists of figures and tables, appendices |
| `essay` | article | classic | gray | analytical prose with footnotes, quotations and references |
| `paper` | article, two columns | classic | gray | research paper: abstract, keywords, theorems, IEEE-style references |

Every template exists in both languages with the same file names. Each `main.tex` lists its files with
`\input{...}`: to add a section, create the file and add one `\input` line. The first block of `main.tex`
(`\atlset{category=..., titlepage=..., ...}`) holds the template's own defaults; your overrides belong in `config.tex`.

## 7. Writing content

### Structure

Use the standard commands: `\section`, `\subsection`, `\paragraph`; `\chapter` in `formal-report`.
Front and back matter (see `formal-report/main.tex`):

| Command | Effect |
|---------|--------|
| `\atlmaketitle` | title page chosen by `titlepage` (call right after `\begin{document}`) |
| `\atlfrontmatter` / `\atlmainmatter` | Roman (English) or letter (Persian: آ، ب، پ) page numbers / back to Arabic from 1 |
| `\atlunnumbered{Abstract}` | heading that is in the table of contents but unnumbered |
| `\begin{abstract} ... \end{abstract}` | tinted abstract box; `\atlkeywords{a, b, c}` adds the keywords line |
| `\atlbibliography{refs}` | bibliography from `refs.bib` (BibTeX, `\cite{key}`) |
| `\atlappendix` | switches numbering to appendix letters (A, B, ... in English; letters in Persian) |

### Boxes

```latex
\begin{notebox}[Optional title] ... \end{notebox}     % also warningbox, tipbox, resultbox,
                                                       % definitionbox, examplebox
\begin{problem}[Title][10] ... \end{problem}           % [Title] and [points] are optional
\begin{solution} ... \end{solution}                    % optional [title]
```

Boxes break across pages. Floats (`table`, `figure`) cannot live inside a box: use
`\begin{center}\captionof{table}{...}...\end{center}` there.

### Tables, figures, captions

Table captions go above, figure captions below. Zebra rows and a colored header:

```latex
\begin{table}[ht]\centering\atlrows
  \caption{Results.}
  \begin{tabular}{lcc}
    \rowcolor{atl-primary}\atlhead{Name} & \atlhead{Value} & \atlhead{Unit}\\
    Resistor & 4.7 & k$\Omega$\\
  \end{tabular}
\end{table}
\includegraphics[width=0.7\linewidth]{diagram}   % file in fig/
```

References: `\cref{fig:x}` / `\Cref{tab:y}` (cleveref) produce "Figure 1" / «شکل ۱».

### Code

```latex
\begin{codebox}[title=Fast Fourier transform]{Python} ... \end{codebox}
\atlinputcode[title=main.c]{C}{../simulation/main.c}      % from a file
\begin{lstlisting}[language=Python] ... \end{lstlisting}  % plain listing
```

Any language known to `listings` works (Python, Matlab, C, C++, Verilog, VHDL, bash, assembler, ...). Code is always
typeset left-to-right, also in Persian documents. Switch theme mid-document with `\atlsetlisting{github}`.

### Maths, units, notation

Theorem-like environments: `theorem lemma proposition corollary definition example remark` and `proof`. Units go through
siunitx: `\qty{2.2}{\micro\farad}`, `\num{1.5e-3}`, `\valunit{10}{\kilo\ohm}`. A small notation set:
`\vect{x}`, `\mat{A}`, `\trans{A}`, `\herm{A}`, `\abs{x}`, `\norm{x}`, `\inner{a,b}`, `\deriv{y}{x}`, `\pderiv{f}{x}`,
`\integral{0}{1}{f(x)}{x}`, `\qtag{R}{in}` (R<sub>in</sub> with an upright subscript), `\N \Z \Q \R \C`, and operators
`\tr \rank \diag \sgn \argmin \argmax \Expect \Var \Cov \Prob`. Text helpers: `\keyterm{}`, `\alert{}`, `\ok{}`, `\highlight{}`,
`\muted{}`, `\todo{}`. Domain-specific macros belong in `preamble/local.tex`; extra packages in `preamble/packages.tex`.

### Diagrams

`tikz` is loaded with the libraries `calc positioning arrows.meta shapes.geometric backgrounds fit`. Load further
packages (`pgfplots`, `circuitikz`) in `preamble/packages.tex`. In Persian documents put pictures inside
`\begin{LTR}...\end{LTR}` (see [section 9](#9-persian-documents)).

## 8. Colors and code themes

Eleven palettes: `blue red gray cyan magenta yellow green sunset ocean berry neon`. Set one in `config.tex` (`palette = sunset`)
or switch mid-document with `\atlsetpalette{green}`. Four code themes: `darcula dark bright github`.

Define your own palette in `preamble/local.tex` (HTML hex values without `#`):

```latex
\atldefpalette{mine}{primary=003D82, secondary=0066CC, accent=D4AF37, light=E6F7FF, tint=EAF2FF}
\atlsetpalette{mine}
```

`success`, `warning`, `info` and `result` are optional. Inside your own content use the palette roles instead of raw
colors so the document follows the palette: `atl-primary`, `atl-secondary`, `atl-accent`, `atl-light`, `atl-tint`,
`atl-success`, `atl-warning`, `atl-info`, `atl-result`.

## 9. Persian documents

The Persian templates need nothing installed: Vazirmatn and HM FElmi (both SIL OFL) are copied into `latex/fonts/`.

- **Latin inside Persian text:** wrap every Latin word, abbreviation, code, ratio or number-with-letter in `\lr{...}`:
  `\lr{IS-A}`, `\lr{1:N}`, `\lr{R1}`. A bare `R1` shows as `۱R`. Decimals must be in math or `\num{}`: `$1.38$`, `\num{1.38}`.
- **Half-space (ZWNJ, U+200C)** inside compound words: `می‌شود`, `کتاب‌ها`. Use the Persian letters `ک` and `ی`, never the Arabic `ك`/`ي`.
- **Spelling of the ezafe after «ه»:** write `ه‌ی` (`شماره‌ی دانشجویی`, `پایگاه‌داده‌ی فیلم`), not the hamza forms `ۀ` or `هٔ`.
- **Dates and digits:** `\today` prints the Jalali date; digits, page numbers and numbering are Persian (`۲-۱`, chapter first).
- **Diagrams:** put `tikzpicture`, `circuitikz` and `axis` inside `\begin{LTR}...\end{LTR}`. Persian text in a node needs
  `execute at begin node=\beginR, execute at end node=\endR` in the node style (otherwise multi-word labels come out in reversed word
  order), and every Latin token or number in a node must be wrapped in `\lr{...}`. Do not use `\underline` on multi-word Persian text
  (it creates a left-to-right box); draw the underline with a small TikZ helper instead.
- **Overriding a label:** `\atldefstr{problem}{سؤال}` in `preamble/local.tex` renames «مسئله» to «سؤال».
- **Fonts:** `fafont = felmi` selects the traditional style; `fascale` fine-tunes the size. The Persian faces have no italic; avoid `\emph` for emphasis.
- **Checking:** after building, search the PDF text for markers of broken references: `pdftotext main.pdf - | grep -nE '\[\?\]|\?\?'` must print nothing.

## 10. Windows, macOS, Overleaf

The generator is a single executable per system; the projects it creates are plain LaTeX and build anywhere.

- **Windows**: download `mktex-windows-amd64.exe` (ARM PCs: `mktex-windows-arm64.exe`) and double-click it, or run it in PowerShell
  (`.\mktex-windows-amd64.exe`). Windows SmartScreen may say the program is unrecognized because it is not code-signed:
  choose *More info*, then *Run anyway*. The file's checksum is in `SHA256SUMS` on the release page. If a `latexmk` build fails because
  Perl is missing (MiKTeX), install Perl, or build with `xelatex main.tex` (twice, and `bibtex main` once if you cite sources), or use VS Code.
- **macOS**: `chmod +x mktex-macos-arm64`, then `./mktex-macos-arm64`. If Gatekeeper blocks the download, run
  `xattr -d com.apple.quarantine ./mktex-macos-arm64` once.
- **Linux**: `chmod +x mktex-linux-amd64`, then `./mktex-linux-amd64`. Move it to a folder on your `PATH` and rename it `mktex` for convenience.
- **Overleaf**: create the project on any computer, zip the `latex/` folder, upload it and set the compiler to **XeLaTeX**.
- **No executable at all**: copy `templates/<template>/<en|fa>/` from the repository to a new folder; copy `core/atl/` to `latex/atl/` (and
  `core/fonts/*` to `latex/fonts/` for Persian); copy `core/project/latexmkrc` to `latex/.latexmkrc`; edit `latex/config.tex` by hand.
- **Build the generator yourself** (needs [Go](https://go.dev/dl/)): `go install github.com/MaximumAsp66915/AdvancedLatexTemplate/cmd/mktex@latest`.

## 11. Troubleshooting

| Symptom | Cause and fix |
|---------|---------------|
| `[?]` or `??` in the PDF | undefined `\cite`/`\ref`: check the key spelling, rebuild with `latexmk` (it repeats passes and runs BibTeX) |
| `[something]` where a label should be | a label name is not defined for this language: add `\atldefstr{name}{text}` (for Persian and English) in `preamble/local.tex` |
| `main.toc:5: Undefined control sequence` on the first run | a stale `main.toc`/`main.aux` is being found in `TEXMFHOME`: delete it and never compile inside the repository |
| Persian glyphs missing or boxes instead of letters | the project lacks `latex/fonts/` or you used `\sffamily` in Persian text (use `\persiansffamily`) |
| Persian numbers or ratios appear reversed (`۱:N` as `N:۱`, `1.38` as `۳۸.۱`) | wrap in `\lr{...}`, or put decimals in math / `\num{}` |
| Multi-word Persian label reversed inside a TikZ node or after `\underline` | see the Diagrams rule in [section 9](#9-persian-documents) |
| `Font shape ... undefined` warning for Persian | harmless: Persian fonts have no italic |
| Table or figure inside a box fails | floats cannot be inside boxes; use `\captionof` in a `center` environment |
| Build fails on first run with missing packages | install the missing TeX Live packages (`tlmgr install <name>`) or use a full distribution |

## 12. Upgrading an old project

The engine is copied into each project, so updating this repository never changes existing projects. To upgrade one,
copy `core/atl/*` over its `latex/atl/` (and `core/fonts/*` over `latex/fonts/` for Persian) and rebuild. Keep your
`config.tex`, `preamble/` and content files. Compare `config.tex` with `core/project/config-options.<lang>.tex` to see
newer options.

## 13. License and contributing

The templates are free to use for any purpose, including commercially, under the [Mozilla Public License 2.0](../LICENSE).
**Your documents are yours**: the license says nothing about the text, figures and PDFs you create. It only asks that if you
*distribute a modified copy of the template files themselves* (the engine in `atl/`, a template, a script) you publish
those changes under the same license. We would much rather you send an improvement back as a pull request than keep
a private copy: see [CONTRIBUTING.md](../CONTRIBUTING.md). The bundled Persian fonts keep their own license (SIL OFL 1.1; the
license texts are in `latex/fonts/`).

---

*Maintainers:* `tools/check.sh` verifies that every template exists in both languages with identical file lists, that
each `config.tex` carries the current options reference (`core/project/config-options.<lang>.tex`) and that all
variants build without errors, missing glyphs, undefined references or `[?]` in the text. `tools/gallery.sh` regenerates `docs/gallery/`.
