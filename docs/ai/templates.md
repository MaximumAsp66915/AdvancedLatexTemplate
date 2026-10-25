# Templates

How to add or change a template, and how to keep English and Persian in step.
Engine details: [engine.md](engine.md). Persian rules (required reading for `fa` files): [persian.md](persian.md).

## Layout and the parity rule

```
templates/<category>/
  template.conf        catalog entry: LABEL, DESCRIPTION, DEFAULT_PALETTE, DEFAULT_LISTING, TITLEPAGE, ASK
  en/  fa/             exactly what mktex copies into <project>/   (the same file list in both!)
    README.md          optional project notes, in that language
    latex/             compile root: main.tex, config.tex, preamble/{packages,local}.tex,
                       sections/ or chapters/ (+ appendix/, frontmatter/), fig/, refs.bib ...
    preprovided/ simulation/ data/ ...   sibling folders the category needs (each with a .gitkeep)
```

**Every category must exist in both languages with identical file names.** Content differs by language rules
(direction, labels, fonts, numbering, typography), never in structure. `tools/check.sh` fails if the file lists differ
and builds every variant.

`mktex` adds on top of a template: `latex/atl/`, `latex/fonts/` (Persian), `latex/.latexmkrc`, a rewritten
`latex/config.tex`, `.vscode/settings.json` and a logo in `latex/fig/`. **Do not commit those inside a template.**

**`config.tex` ends with the options reference.** The block comes from `core/project/config-options.<lang>.tex`;
`mktex` appends it to the generated `config.tex`, and every `templates/<cat>/<lang>/latex/config.tex` must end with
the identical block (`check.sh` verifies this). When you add or change a configuration key: edit both `config-options`
files, replace the block in every template's `config.tex`, and update `docs/USER-MANUAL.md` section 5.

## Catalog

| Category | Class | Title page | Default palette | Use |
|---|---|---|---|---|
| `homework` | article | block | blue | problem sets: `problem` / `solution` boxes, code, results |
| `lab-report` | article | block | cyan | objective, theory (circuit figure), apparatus, procedure, data table and plot, analysis, references |
| `project-report` | article | cover | sunset | cover, abstract, contents, requirements, design diagram, implementation, tests |
| `formal-report` | report | cover | gray | thesis or internship report: chapters, roman / abjad front matter, lists of figures and tables, appendices |
| `essay` | article | classic | gray | indented prose, footnotes, quotation, references |
| `paper` | article, two-column | classic | gray | IEEE-like paper with abstract, keywords, theorem, figure, table |

## Writing a template

1. **`main.tex` skeleton** (copy a sibling): `\documentclass` -> `\input{atl/atl-keys}` ->
   `\atlset{category=..., titlepage=..., margins=..., infofields={...}, ...defaults}` -> `\input{config}` ->
   `\input{atl/atl}` -> `\input{preamble/local}` -> `\begin{document}\atlmaketitle ...`. Persian adds `language=fa`
   to the first `\atlset`, and starts with `% !TeX program = xelatex` like the English one.
2. **Defaults live in `main.tex`.** Anything needed to compile without a `config.tex` (title, author, affiliation,
   keywords ...) goes in its `\atlset`; `config.tex` is overwritten by `mktex` and holds only what the user typed.
3. **`infofields`** may only use keys that have a label in both language layers; add a missing label with
   `\atldefstr` in both.
4. **Use palette roles** (`atl-primary` ...), never raw colors; use the boxes, `\atlrows`, `\atlhead`, `codebox`.
5. **Tables**: inside a box use `center` + `\captionof{table}`; otherwise a normal `table`. Table captions go above,
   figure captions below.
6. **Persian files** follow every rule in [persian.md](persian.md) (Latin in `\lr{}`, decimals in math or `\num`,
   diagrams in `LTR`, the ezafe spelling, TikZ node rules ...).
7. **`template.conf`**: `ASK` lists the config keys `mktex` prompts for (labels come from `promptLabels` in
   `cmd/mktex/main.go`), `TITLEPAGE` must match `main.tex` (the logo is only copied for `cover`).
8. **Verify**: `tools/check.sh <category>`, then render both languages to images and look at them
   ([workflow.md](workflow.md#verification-protocol)).

## What the two language versions differ in (besides the text)

| | en | fa |
|---|---|---|
| Direction and sides | LTR, binding on the left | RTL, binding on the right; bars and header sides mirror automatically |
| Body font | TeX Gyre Termes or Pagella, Heros headings | Vazirmatn or HM FElmi (`fafont`), Vazirmatn headings, Latin text in TeX Gyre |
| Numbering | `2.1` | `۲-۱` (Persian digits; components read right to left) |
| Front-matter pages | roman numerals | abjad letters (`\harfi`) |
| `\today` | Gregorian | Jalali |
| Line spacing | slightly lower | slightly higher (Persian fonts need more) |

## Adding a new category

1. Create `templates/<name>/{template.conf, en/, fa/}` with the same file lists.
2. Add a `README.md` per language only if the category has one (parity applies to it as well).
3. Run `tools/check.sh <name>`; update the catalog table above, the README, `docs/USER-MANUAL.md` section 6,
   and regenerate the gallery with `tools/gallery.sh`.
