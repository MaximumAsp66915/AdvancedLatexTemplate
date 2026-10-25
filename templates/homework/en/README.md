# Homework

Compile from `latex/`:

```bash
latexmk -xelatex main.tex
```

| Path | Purpose |
|------|---------|
| `latex/main.tex` | Entry point. Template defaults live in the first `\atlset{...}` block |
| `latex/config.tex` | Your title, name, course, palette and code theme (written by `mktex`) |
| `latex/sections/` | One file per problem. Use `problem` and `solution` environments |
| `latex/appendix/` | Code listings |
| `latex/fig/` | Figures (referenced without the folder name) |
| `latex/preamble/` | `packages.tex` (extra packages) and `local.tex` (your macros) |
| `latex/atl/` | The engine, copied by `mktex`. You normally do not edit it |
| `preprovided/`, `simulation/` | Course material and simulation files |

Change the colors with `palette = ...` in `config.tex`
(blue, red, gray, cyan, magenta, yellow, green, sunset, ocean, berry, neon).
