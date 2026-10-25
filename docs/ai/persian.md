# Persian (Farsi) typesetting guide

Everything an agent needs to write or change Persian templates correctly. Layers: the Persian language layer is
`core/atl/atl-lang-fa.tex`; the templates are `templates/<category>/fa/`. Read [engine.md](engine.md) for the engine and
[templates.md](templates.md) for parity. **"Verified"** below means compiled or inspected with XeLaTeX (TeX Live 2023,
XePersian v25.0, siunitx 3.3, Linux). **"Reported"** means taken from a source document and not independently checked.

## Sources behind this guide

- A university department's internship-report **Word template** (not redistributed): useful as a *specification* of fonts,
  sizes, margins, numbering, Persian typography and document skeleton. Its rules are in the sections below.
- A third-party modular XeLaTeX project by Hossein Dadashzadeh (MIT license, 2026), used for **ideas only** (color roles
  with semantic tokens, labelled boxes, left-to-right code boxes, siunitx conventions). No code was copied; the engine here
  is an independent implementation. If you ever reuse its code you must keep its copyright notice.

## How the Persian side is built

- XeLaTeX only; `main.tex` starts with `% !TeX program = xelatex`; `.latexmkrc` sets `$pdf_mode = 5`.
- Self-contained: `mktex` copies `core/atl/` and `core/fonts/` into the project, so nothing depends on system fonts or the repository.
- `report` class for the formal report (chapters, «فصل ۱»), `article` for the rest. XePersian ships adapters for `article report
  book extreport extbook scrartcl scrreprt scrbook memoir` and for `natbib tocloft listings enumitem footnote hyperref
  glossaries beamer fancyref varioref` (verified in the TeX Live tree).
- Numbering uses a hyphen for sections, floats and equations (`۲-۱`); the first component (chapter) is at the **right**,
  as in the Word template (verified).
- `\today` prints the **Jalali** date; `\latintoday` the Gregorian one. Front matter uses letters: `\harfi{page}`
  (آ ب پ ت ...); `\abjad{...}` and `\Abjad{...}` also exist (verified).
- Page numbers, digits in text and in math, and the decimal separator come out Persian (`۵٫۴`) (verified).

## Fonts

- Bundled in `core/fonts/`: **Vazirmatn** (Saber Rastikerdar) and **HM FElmi** (Hossein Movahhedian), both **SIL OFL 1.1**
  (read from the font name tables). The OFL requires each redistributed copy to keep the copyright notice and license text,
  so the `OFL-*.txt` files travel with the fonts. Do not add proprietary fonts (B Nazanin, B Lotus, B Zar, Times New
  Roman). A Persian template must not fail because a system font is missing: the engine falls back to Amiri or Noto where
  the bundled file is absent.
- Latin text in Persian documents uses TeX Gyre (Termes with FElmi, Heros with Vazirmatn); both ship with every TeX distribution.
- XePersian commands: `\settextfont[Path=...,BoldFont=...]{file.ttf}`, `\setdigitfont`, `\setlatintextfont`,
  `\setpersiansansfont`. Verified: a source project that names *Times New Roman* fails on Linux
  (`The font "Times New Roman" cannot be found`).
- Check glyph coverage for `ک ی` ZWNJ `٫ ٪ ؛ ؟ « »` and Persian digits before adding any font.
- The engine's keys: `fafont = vazirmatn | felmi`, `fascale`.

## Writing and typesetting rules

1. **Half-space (ZWNJ, U+200C)** inside compound words and between a stem and its suffix or prefix: `می‌شود`, `اندازه‌گیری`,
   `کتاب‌ها`. Never leave a real space inside one word. Type the character itself.
2. Use the Persian letters **ک (U+06A9)** and **ی (U+06CC)**, never the Arabic `ك` / `ي`.
3. **Ezafe after a final heh is written heh + ZWNJ + yeh** (`شماره‌ی دانشجویی`, `پایگاه‌داده‌ی فیلم`), not with the hamza
   forms U+06C0 or heh + U+0654. Audit (must print nothing apart from `docs/USER-MANUAL.md`, which names the forms as a warning):
   `grep -rnP '\x{06C0}|\x{0647}\x{0654}' --include=*.tex --include=*.md .`
4. Punctuation: `. ، : ؛ ؟ ! ) ] }` attach to the previous word and are followed by one space; `( [ {` are preceded by a space and
   attach to the next word. Use `،` `؛` `؟` and «...» for quotes. A sentence that ends in a displayed formula ends with `:`.
   A full stop after a reference marker goes after the bracket: `... نتیجه [۳].`
5. **Latin inside Persian text goes in `\lr{...}`**: words, abbreviations, codes, ratios, anything mixing letters and digits
   (`\lr{RF}`, `\lr{R1}`, `\lr{1:N}`, `\lr{IS-A}`). A bare `R1` renders as `۱R`; a bare `1:N` as `N:۱`.
6. **Decimals** typed as text (`1.38`) render reversed (`۳۸.۱`). Write them in math (`$1.38$`), `\num{1.38}` / `\qty{...}{...}`,
   or siunitx `S` columns. Units always through siunitx.
7. The Persian comma `،` inside displayed math is missing (Computer Modern): use `\text{،}`.
8. Equations numbered per chapter `(۲-۱)`; figures and tables per chapter in order of appearance (`شکل ۲-۱`);
   **table captions above, figure captions below**; every float referenced in the text with `\cref`.
9. References: IEEE style (BibTeX + IEEEtran), numeric `[1]` in citation order, typeset left-to-right in Latin. Persian-language
   entries were not tried (they would need `persian-bib.bst` or manual RTL handling).
10. Every chapter ends with a «جمع‌بندی» section.
11. Footnote numbers restart on every page in the Word convention; this is **not implemented** (numbering is continuous).

## Formal academic report conventions (reported by the Word template)

| Item | Value | In `templates/formal-report/fa` |
|---|---|---|
| Paper, margins | A4; top/bottom 2.5 cm; binding side (right) 3 cm; other side 2.5 cm | `margins = formal` |
| Fonts | Persian 14 pt, Latin 12 pt | `12pt` with FElmi at scale 1.1 (about 13 pt); raise `fascale` to approach 14 |
| Line spacing | 1.5 | `linespread = 1.4` (LaTeX single spacing is already about 1.2) |
| Paragraph | justified, first-line indent 0.71 cm | `parstyle = indent` |
| Chapter heading | centered bold, «فصل اول: ...» | chapter label above a large title with a rule |
| Heading numbers | `1`, `1-2`, `1-2-3` | yes |
| Captions | centered; table above, figure below; numbered per chapter | yes |
| Page number | bottom centre, none on the title page; letters in front matter, digits from the first chapter | yes (`\harfi`) |
| Skeleton | title page, abstract (چکیده), contents (فهرست مطالب), list of figures and tables, chapters (first and last mandatory), references (مراجع), appendices (پیوست) | yes |

## Traps (all reproduced and fixed in `core/atl`; keep this list current)

1. `\sffamily` is the **Latin** sans. In Persian text it selects Heros and every Persian glyph is silently missing.
   Use `\persiansffamily` (the engine hook is `\atl@sf`).
2. **Listings in Persian documents are corrupted** by XePersian's digit conversion even under `\setLTR`
   (`h:.0e` becomes `h0.:e`, `3.14` becomes `14.3`). Fix: also `\XeTeXinterchartokenstate=0` and `\resetlatinfont`.
   The engine bundles all of it in `\atl@ltr`; apply it **before** the font switch. The bibliography uses it too.
3. Persian fonts have **no italic**: map `ItalicFont` to the upright file, otherwise XePersian warns and substitutes
   (also needed for the digit font used in math scripts).
4. fancyhdr: a header rule built from an hbox with leaders demands 35 pt of `\headheight`; the `\hrule` form does not.
5. Floats cannot sit inside tcolorbox; use `\captionof` in a `center` environment.
6. Diagrams (`tikzpicture`, `circuitikz`, pgfplots `axis`) go inside `\begin{LTR}...\end{LTR}`. Persian words in `axis` legends work.
7. `pgfplots`, `circuitikz` and other `\usepackage` lines must come **before** XePersian: use the `preamble/packages.tex` hook.
8. **Text inside TikZ nodes**: a node's text is an hbox, so `\setRTL` does nothing for word order. A multi-word Persian label
   comes out reversed ("به مربوط" for "مربوط به"). What works: `execute at begin node=\beginR, execute at end node=\endR` in the
   node style, or `\rl{...}` by hand. Every Latin token, digit group or ratio in a node still needs `\lr{...}`; write `\lr{1:N}` as one unit.
9. **`\underline`, `\mbox`, `\fbox` and other hbox commands** typeset a left-to-right box and reverse multi-word Persian.
   For underlined key attributes use a TikZ node helper that applies `\beginR...\endR` itself.
10. **Do not judge word order by eye** (a reader of Latin script gets it wrong easily). Check coordinates:
    `pdftotext -bbox main.pdf out.html`; in a correct RTL phrase the first logical word has the larger `xMin`.
11. **`\raggedright` is physically left** in XePersian. Use `\atl@flushstart` in the engine and
    `>{\raggedleft\arraybackslash}p{...}` in tables.
12. **`[?]`, `??` or `[key]` in the output** means a broken reference or a missing label: `[?]` / `??` = undefined `\cite` / `\ref`;
    `[key]` = an `\atlstr{key}` that is not defined in both language layers. `check.sh` greps the PDF text for them.
13. Overriding a label (for example «سؤال» instead of «مسئله») is done in `preamble/local.tex` with `\atldefstr{problem}{...}`.
14. Numbers like `۶-۱` are read right to left: chapter first, at the right.
15. **Persian punctuation and digits inside `\lr{...}` are missing glyphs** (verified: `،`, `۱`, `۲` give `Missing character` in the Latin font and print nothing). Keep `\lr{}` for Latin text only.
16. **Inline `\lstinline`** used to scramble code in Persian text; the engine now fixes it (see [engine.md](engine.md) rule 14). For inline code in
    running text the manual-friendly form is still `\mbox{\lr{\texttt{...}}}`: the `\mbox` stops a line break in the middle of code, which would
    otherwise scramble the order of brackets in right-to-left text. Do not put it inside bold or headings (the mono font has no bold).
17. Reported from a real homework build, **not reproduced in a minimal test**: TikZ labels placed on a path (`\draw ... -- node[...]{...}`) with
    `execute at begin node=\beginR` lost their text in a larger figure; a stand-alone `\node` at the middle of the edge worked reliably.
    Wrap inline formulas that must not break in `\mbox{$...$}`.
18. **Dark code themes and page breaks**: see [engine.md](engine.md) rules 12 and 13; `check.sh` tests it in English and Persian.

## Quick self-check before you finish a Persian change

- [ ] Latin tokens in `\lr{}`; decimals in math / `\num`; no bare `R1`, `1:N`, `3.14` in running text.
- [ ] ZWNJ in compound words; `ه‌ی` ezafe; Persian `ک ی`; the audit grep above prints nothing.
- [ ] Diagrams in `LTR`; node text with `\beginR ... \endR`; no `\underline` on multi-word Persian.
- [ ] No `\sffamily` in Persian paths; no bare `\raggedright` in Persian paths.
- [ ] `tools/check.sh` passes; rendered pages looked at; word order verified by coordinates where it matters.
