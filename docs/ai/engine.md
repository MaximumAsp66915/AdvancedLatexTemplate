# The engine (`core/atl/`)

`core/atl/` ("atl" = AdvancedLatexTemplate) is the one shared, **language-neutral** TeX engine behind every
template. `mktex` copies it into each project as `latex/atl/`. **XeLaTeX only** (fontspec; XePersian for Persian).
Related: [templates.md](templates.md) (using the engine from a template), [persian.md](persian.md) (the Persian layer).

## Load order and module roles

A template's `main.tex` runs: `\input{atl/atl-keys}` -> `\atlset{template defaults}` -> `\input{config}` (user values)
-> `\input{atl/atl}` (loads everything below, in this order) -> `\input{preamble/local}`.

| File | Role |
|---|---|
| `atl-keys` | `\atlset`, `\atlmeta`, `\atlifset`, `\atlifeq`, `\atlifon`, `\atldefault`, `\atlstr`, `\atldefstr` |
| `atl-packages` | every package, once; then the user hook `preamble/packages.tex` (extra `\usepackage` lines go here, before XePersian) |
| `atl-lang-en` / `atl-lang-fa` | the **only** language-specific code: fonts, direction hooks, numbering style, front-matter numbering, all labels |
| `atl-palettes` | `\atldefpalette`, `\atlsetpalette`; roles `atl-primary secondary accent light tint success warning info result` plus derived colors |
| `atl-layout` | geometry (`margins`), line spacing, paragraph style |
| `atl-headings` | titlesec styles, numbering, `\atlunnumbered`, `\atlfrontmatter` / `\atlmainmatter` / `\atlappendix`, `abstract`, `\atlkeywords` |
| `atl-pagestyle` | fancyhdr header and footer (start and end sides mirror in Persian) |
| `atl-floats` | captions (tables above, figures below), `\atlhead`, `\atlrows`, siunitx setup |
| `atl-boxes` | `notebox warningbox tipbox resultbox definitionbox examplebox`, `problem`, `solution` |
| `atl-code` | code themes (`\atldeflisttheme`, `\atlsetlisting`), `codebox`, `\atlinputcode` |
| `atl-macros` | text helpers and a small general notation set (`\vect \mat \trans \qtag \valunit ...`) |
| `atl-theorems` | amsthm styles: theorem, lemma, proposition, corollary, definition, example, remark, proof |
| `atl-titles` | `\atlmaketitle` with `titlepage = block | cover | classic | none`, `\atlinfogrid` |
| `atl-bib` | `\atlbibliography` (BibTeX + IEEEtran) and its own `thebibliography` |

## Configuration keys

Set with `\atlset{key=value}`. Any name is accepted; unknown keys are simply stored and read with `\atlmeta{key}`.
A later setting wins.

- **Template and layout**: `language` (en|fa), `category`, `titlepage` (block|cover|classic|none),
  `margins` (normal|narrow|formal|wide), `linespread`, `parstyle` (block|indent), `headfont` (sans|serif),
  `headsize` (normal|compact), `numbering` (flat|section), `sectionbreak` (on|off), `header` (on|off),
  `headerstart`, `headerend`, `infofields` (comma list of title-page info keys), `bibstyle`.
- **Look**: `palette`, `listing` (darcula|dark|bright|github), `enfont` (termes|pagella|lm), `fafont` (vazirmatn|felmi), `fascale`.
- **Document data** (title page; labels come from the language layer): `title subtitle author studentid course
  coursecode instructor supervisor advisor team partners experiment labsection term date duedate submissiondate
  performed university department affiliation keywords logo doctype`.

Helper API: `\atlmeta{k}` value; `\atlifset{k}{yes}{no}`; `\atlifeq{k}{value}{yes}{no}`; `\atlifon{k}{yes}{no}`
(true/on/yes); `\atldefault{k}{v}` sets only if unset. `\atlstr{name}` gives a localised label (prints `[name]` if undefined).

## Rules learnt the hard way (each one was reproduced by compiling)

1. **Never put `&` or `\\` inside `\if...\fi`.** They end the table cell while the conditional is still open
   ("Incomplete \ifodd"). The `\atlif...` macros select the branch first through
   `\csname atl@if\atl@isset{...}\endcsname{yes}{no}`; keep new conditionals in that style. They are `\long`, so
   branches may contain `\par`.
2. **Every module ends with `\makeatother`.** Loader code after the first `\input` is read with `@` as an ordinary
   character, so wrap such code in `\makeatletter ... \makeatother` (as `atl.tex` does for its `\AtBeginDocument` block).
3. **A macro takes at most nine arguments.** That is why `\atldeflisttheme` and `\atldefpalette` take key=value lists.
4. **`literate=` and other listings keys do not expand macros.** Write the list out literally, or use tcolorbox's
   `listing options/.expanded=` when a value comes from a macro.
5. **Floats cannot live inside tcolorbox** (`solution`, all boxes): use `\captionof{table}{...}` inside a `center` environment.
6. **fancyhdr**: a header rule built as an hbox with leaders makes fancyhdr demand 35 pt of `\headheight` under XePersian.
   Use the `\hrule` form (`\atl@headrule`).
7. **fontspec breaks `\mathrm{\Omega}`** (the math "operators" font follows the text font). Both language layers reset it
   to Computer Modern in `\AtBeginDocument`.
8. **Italic**: Persian fonts have none, so `\atl@fa@try` maps italic to the upright file (no substitution warning, also for
   the digit font used in math scripts). Latin Modern Mono has no bold, so code keywords are colored, not bold.
9. **Persian sans is `\persiansffamily`**, not `\sffamily` (that is the Latin sans; Persian glyphs vanish silently).
   Use the hook `\atl@sf`. Code and Latin runs use `\atl@ltr` (see [persian.md](persian.md)).
10. **Missing label** shows as `[key]` in the output: add `\atldefstr{key}{...}` to **both** language layers.
11. **`\raggedright` is physically left in XePersian.** Use the hook `\atl@flushstart` (en: `\raggedright`, fa: `\raggedleft`).
12. **Every glyph needs its own color.** In listings, ordinary punctuation has no style of its own (`\lst@thestyle` is
    `\relax`), so it takes the surrounding text color; inside a breakable box that color is black again on every page after
    the first. Result on a dark theme: `& # % ^ ~ | `` and `)` drawn black (a real bug, found in a homework). The engine fixes it in
    `atl-code.tex` with an `OutputOther` hook that gives such tokens the code color, and sets `colupper` on `codebox`. Never rely
    on inherited color for text on a dark background.
13. **`)` in listings with `breaklines=true`**: listings re-defines `)` (to allow a break after it), which silently disables any
    `literate` rule for `)`. The same hook colors it with the bracket color. Do not try to color `)` with `literate`.
14. **Inline `\lstinline` inside Persian text** comes out scrambled (`pid(waitpid , status& ,0 );`). The engine wraps inline listings
    in `\beginL ... \endL` (listings `Init`/`DeInit` hooks, Persian only). Block listings are handled by `\atl@ltr`.
15. **Colors are audited**: `tools/audit-colors.py` checks every code-theme token against its background (>= 3:1, comments >= 2.5)
    and every palette role used for text (white on `primary` >= 4.5, `secondary-text` >= 3, box titles >= 4.5). It found that
    the original `neon` secondary (`08D9D6`) was too light for sub-headings; it is now `00B5B2`. `check.sh` runs it.
16. Front-matter numbering, appendix numbering, numbering separator and the start side are hooks:
    `\atl@frontnumbering` (`roman` / `harfi`), `\atl@appendixnumber`, `\atl@numsep` (`.` / `-`), `\atl@start` / `\atl@end`, `\atl@L` / `\atl@R`.

## Language-layer contract

A language layer (`atl-lang-<code>.tex`) must define: `\atl@start \atl@end \atl@L \atl@R \atl@sf \atl@ltr \atl@lr
\atl@numsep \atl@frontnumbering \atl@appendixnumber \atl@startskip \atl@flushstart`, every `\atldefstr` label
(copy the list from `atl-lang-en.tex`), cleveref names, and the font setup. Shared modules may only use these hooks
for direction- or language-dependent behavior.

## Adding things

- **New engine, tool or Go file**: start it with the one-line notice `% SPDX-License-Identifier: MPL-2.0` (`#` for scripts, `//` for Go). The engine files are
  copied into every project, so the notice travels with them.

- **Palette**: one line in `atl-palettes.tex` (then run `python3 tools/audit-colors.py`; it must stay green):
  `\atldefpalette{name}{primary=..., secondary=..., accent=..., light=..., tint=...}` (HTML hex, no `#`;
  `success warning info result` optional). It appears in `mktex --list` automatically. `primary` must be dark enough
  for text on white. Update `docs/USER-MANUAL.md` section 8 and the palette lists in the README and the options reference.
- **Code theme**: one `\atldeflisttheme` line in `atl-code.tex` and the name in `listings` in `cmd/mktex/main.go`; run the color audit.
- **Label**: `\atldefstr{key}{...}` in **both** language layers.
- **Configuration key**: read it with `\atlifeq`/`\atlmeta` in the module that needs it, give it `\atldefault` if it has a
  default, then update the options reference (see [README.md](README.md#keeping-this-manual-correct)).
- **Language**: new `atl-lang-<code>.tex`, a loader branch in `atl.tex`, a `<code>/` folder in every template.
- **Box or title-page style**: use only palette roles and the language hooks; test in both languages and both directions.

After any change run `tools/check.sh` (add `--palettes` for colors) and look at rendered pages.
