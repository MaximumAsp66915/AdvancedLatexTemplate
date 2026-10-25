# Contributing to AdvancedLatexTemplate

Thank you for using these templates. They exist so that students and researchers in English and Persian
can produce good-looking documents without fighting LaTeX. They are free to use, and the best way to make them
better for everyone is to **send your improvements back to this project** instead of keeping them in a private copy.

## Using the templates (no contribution needed)

- Use them for any document, for any purpose, free of charge. **The documents you write are yours**: your text, your
  PDFs and your own content files carry no obligation at all.
- You may copy the templates into your projects (the generator does exactly that) and adapt the *content files*
  freely. Changing `config.tex`, `preamble/local.tex` and your sections is normal use, not a fork.

## Why we ask you to contribute changes

The project is licensed under the [Mozilla Public License 2.0](LICENSE). In plain words:

- If you **distribute** a modified copy of the engine, a template, or any other file from this repository, you must
  publish your changes to those files under the same license.
- Private modification for your own use is always allowed; no open license can forbid that, and we do not ask you to hide it.
- What we **ask**, beyond the license, is simple: if you fixed something or added something useful, open a pull request.
  A change that lives in one person's copy helps one person; merged here, it helps every student who uses the
  templates, in both languages, and you stop having to re-apply it after every update.

## Ways to help

| You found / want | Do this |
|---|---|
| A bug (wrong layout, missing glyph, build error) | Open an issue with the template, language, `config.tex`, the relevant part of the log and, if you can, a screenshot |
| A mistake in the Persian or English text | Fix it in a pull request. Native review of the Persian text is especially valuable |
| A new palette or code theme | One line in `core/atl/atl-palettes.tex` or `atl-code.tex`; see `docs/ai/engine.md` |
| A new template category | Both languages, same file names; see `docs/ai/templates.md` |
| A new language | See "Language" in `docs/ai/engine.md`; open an issue first so we can agree the plan |
| Testing on Windows or macOS, CI, code signing of the executables | Very welcome; open an issue to coordinate |

## Before you open a pull request

1. Fork the repository and create a branch for one logical change.
2. Make the change in **both languages** when it touches a template (English and Persian must keep the same file
   names). Read `docs/ai/README.md` for the map of the project: it is written for AI agents but is the most complete
   description of the rules for people too.
3. Run `tools/check.sh`. It must end with `PASS`. It checks the structure of every template in both languages, the
   color contrast of every palette and code theme, builds every template, and checks that dark code themes never draw
   black glyphs. It needs XeLaTeX, `latexmk`, bash and Go (to build the generator); `python3` and `mutool` enable the color checks. The scripts are in `tools/`.
4. **Look at the output.** Render the pages (`pdftoppm -r 60 -png main.pdf p`) and check them in both languages. For a visual
   change attach before/after screenshots to the pull request.
5. Update the documentation that describes what you changed (the table in `docs/ai/README.md` lists which file covers which change).
6. Do not include personal data (names, student IDs, e-mails), university logos, proprietary fonts, build files
   (`*.pdf`, `*.aux`, ...), or code and text copied from somewhere without a compatible license.

Keep pull requests small and focused, and describe what changed and what you verified (and what you could not verify,
for example Windows or macOS).

## Persian contributions

Persian typesetting has several traps (digits, word order in diagrams, fonts). Read `docs/ai/persian.md` before editing
any `fa` file. The most common ones: Latin words in `\lr{...}`, decimals in math or `\num{}`, zero-width non-joiner in
compound words, diagrams inside `\begin{LTR}`.

## AI-assisted contributions

AI tools are welcome. Point your agent at [`AGENTS.md`](AGENTS.md); the manual in `docs/ai/` is written for it. You are
responsible for what you submit: run `tools/check.sh` yourself and look at the rendered pages.

## License of contributions

By submitting a contribution you agree that it is licensed under the Mozilla Public License 2.0, like the rest of the
project, and that you have the right to submit it. The bundled fonts in `core/fonts/` keep their own license (SIL Open
Font License 1.1); do not remove their license texts.
