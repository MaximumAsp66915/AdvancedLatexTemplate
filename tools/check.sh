#!/bin/bash
# SPDX-License-Identifier: MPL-2.0
# check.sh -- verify the template catalog.
#   1. parity: every template exists in en and fa with the same file list
#   2. colors: every palette / code-theme color pair is readable (needs python3, no TeX)
#   3. build : generate every variant with mktex and compile it with XeLaTeX;
#              fail on errors, missing glyphs, undefined references, '[?]' / '??' in the PDF text
#   4. code colors: a long dark-theme listing inside a breakable box, in en and fa, must not
#              contain a single black glyph (needs mutool; skipped when it is not installed)
#
#   tools/check.sh                    everything
#   tools/check.sh --no-build         parity only (fast)
#   tools/check.sh homework paper     only these categories
#   tools/check.sh --palettes         also build one template with every palette
set -u
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd -P)"
BUILD=1; PALETTES=0; CATS=()
for a in "$@"; do
  case "$a" in --no-build) BUILD=0;; --palettes) PALETTES=1;; -h|--help) sed -n '2,11p' "$0"; exit 0;; *) CATS+=("$a");; esac
done
if [[ ${#CATS[@]} -eq 0 ]]; then
  for d in "$ROOT"/templates/*/template.conf; do CATS+=("$(basename "$(dirname "$d")")"); done
fi
FAIL=0
bad() { echo "  FAIL: $*"; FAIL=1; }

# The generator is a Go program. MKTEX=/path/to/mktex uses a ready-made binary; otherwise it is built here.
# It always reads templates from this checkout (--root), never from its built-in copy.
MKTEX_BIN="${MKTEX:-}"
build_generator() {
  [[ -n "$MKTEX_BIN" ]] && return 0
  export PATH="$PATH:/usr/local/go/bin"
  command -v go >/dev/null 2>&1 || { bad "Go is needed to build the generator (https://go.dev/dl/) or set MKTEX=/path/to/mktex"; return 1; }
  MKTEX_BIN="$(mktemp -d "${TMPDIR:-/tmp}/atl-mktex.XXXXXX")/mktex"
  ( cd "$ROOT" && go build -o "$MKTEX_BIN" ./cmd/mktex ) || { bad "go build failed"; return 1; }
}
mktex() { "$MKTEX_BIN" --root "$ROOT" "$@"; }

echo "== parity"
for c in "${CATS[@]}"; do
  T="$ROOT/templates/$c"
  [[ -f "$T/template.conf" ]] || { bad "$c: no template.conf"; continue; }
  for l in en fa; do [[ -d "$T/$l/latex" ]] || bad "$c: missing $l/latex"; done
  [[ -d "$T/en/latex" && -d "$T/fa/latex" ]] || continue
  for l in en fa; do   # config.tex must end with the options reference (core/project/config-options.$l.tex)
    B="$ROOT/core/project/config-options.$l.tex"; n="$(wc -l < "$B")"
    tail -n "$n" "$T/$l/latex/config.tex" | diff -q - "$B" >/dev/null || bad "$c/$l: config.tex lacks the current options reference"
  done
  if ! diff <(cd "$T/en" && find . -type f | sort) <(cd "$T/fa" && find . -type f | sort) >/dev/null; then
    bad "$c: en and fa file lists differ"; diff <(cd "$T/en" && find . -type f | sort) <(cd "$T/fa" && find . -type f | sort) | sed 's/^/        /'
  else echo "  ok   $c ($(cd "$T/en" && find . -type f | wc -l | tr -d ' ') files per language)"; fi
done
echo "== generator (Go)"
if command -v go >/dev/null 2>&1 || [[ -x /usr/local/go/bin/go ]]; then
  export PATH="$PATH:/usr/local/go/bin"
  ( cd "$ROOT" && go vet ./... && go test ./... 2>&1 | sed 's/^/  /' ) || bad "go vet / go test failed"
else echo "  skipped (Go not found)"; fi
echo "== colors"
if command -v python3 >/dev/null 2>&1; then
  python3 "$ROOT/tools/audit-colors.py" | sed 's/^/  /'
  [[ ${PIPESTATUS[0]} -eq 0 ]] || bad "color contrast audit failed (see above)"
else echo "  skipped (python3 not found)"; fi
[[ $BUILD -eq 0 ]] && { [[ $FAIL -eq 0 ]] && echo "PASS" || echo "FAILED"; exit $FAIL; }

build_generator || { echo "FAILED"; exit 1; }
TMP="$(mktemp -d "${TMPDIR:-/tmp}/atl-check.XXXXXX")"
echo "== build (in $TMP)"
build_one() { # category lang palette [name]
  local c="$1" l="$2" p="${3:-}" n="${4:-${1}_${2}}" extra=() out
  [[ -n "$p" ]] && extra=(-p "$p")
  if [[ "$l" == fa ]]; then A=("-s" "title=عنوان نمونه" "-s" "author=نام نمونه" "-s" "course=درس نمونه"); else A=("-s" "title=Sample Title" "-s" "author=Sample Author" "-s" "course=Sample Course"); fi
  mktex --no-defaults -y --no-open -l "$l" -c "$c" -n "$n" -o "$TMP" ${extra[@]+"${extra[@]}"} "${A[@]}" >/dev/null 2>"$TMP/$n.mk.err" || { bad "$n: mktex failed ($(head -1 "$TMP/$n.mk.err"))"; return; }
  ( cd "$TMP/$n/latex" && latexmk -xelatex -interaction=nonstopmode -file-line-error main.tex >"$TMP/$n.out" 2>&1 )
  local rc=$? log="$TMP/$n/latex/main.log" pages errs=0
  pages="$(pdfinfo "$TMP/$n/latex/main.pdf" 2>/dev/null | awk '/^Pages/{print $2}')"
  local miss over undef qm
  qm="$(pdftotext "$TMP/$n/latex/main.pdf" - 2>/dev/null | grep -c -E '\[\?\]|\?\?')"   # unresolved cite/ref markers in the PDF text
  miss="$(grep -c 'Missing character' "$log" 2>/dev/null)"; over="$(grep -c '^Overfull' "$log" 2>/dev/null)"
  undef="$(grep -c -E "(Reference|Citation).*undefined" "$log" 2>/dev/null)"
  local msg="pages=${pages:--} missing_glyphs=$miss overfull=$over undefined=$undef unresolved_marks=$qm"
  if [[ $rc -ne 0 || -z "$pages" || "$miss" != 0 || "$undef" != 0 || "$qm" != 0 ]]; then bad "$n: $msg (rc=$rc, see $TMP/$n.out)"
  else printf '  ok   %-28s %s\n' "$n" "$msg"; fi
}

# A listing that runs over a page break inside a breakable box must keep its colors on every page.
# Regression test for: `)' `&' `#' `%' `^' `~' `|' drawn black on dark themes after the first page.
listing_colors() { # lang
  local l="$1" d="$TMP/lstcolors_$1" i
  mkdir -p "$d"; cp -R "$ROOT/core/atl" "$d/atl"
  if [[ "$l" == fa ]]; then mkdir -p "$d/fonts"; cp "$ROOT"/core/fonts/*.ttf "$d/fonts/"; fi
  {
    printf '%s\n' '% !TeX program = xelatex' '\documentclass[11pt,a4paper]{article}' '\input{atl/atl-keys}' \
      "\\atlset{language=$l,listing=darcula,header=off}" '\input{atl/atl}' '\pagestyle{empty}' \
      '\begin{document}' '\begin{solution}' '\begin{lstlisting}[language=C]'
    for i in $(seq 1 70); do printf '    r%d = wait(&status); if ((a & b) | ~c ^ d %% e) { g(x); } /* (c) # ` */\n' "$i"; done
    printf '%s\n' '\end{lstlisting}' '\end{solution}' '\end{document}'
  } > "$d/t.tex"
  ( cd "$d" && latexmk -xelatex -interaction=nonstopmode t.tex >"$d/out" 2>&1 )
  local res pages total black
  res="$(mutool draw -q -F trace -o - "$d/t.pdf" 2>/dev/null | awk '
    /<page /      { pg++ }
    /<fill_text/  { match($0, / color="[^"]*"/); col = substr($0, RSTART + 8, RLENGTH - 9); mono = 0 }
    /LMMono/      { mono = 1 }
    /<g unicode=/ { if (mono) { t++; if (col == "0") b++ } }
    END           { print pg + 0, t + 0, b + 0 }')"
  read -r pages total black <<<"$res"
  if [[ -z "$total" || "$total" -eq 0 || "$pages" -lt 2 ]]; then
    bad "listing colors ($l): test did not run as intended (pages=${pages:-?}, glyphs=${total:-?}; see $d/out)"
  elif [[ "$black" -ne 0 ]]; then
    bad "listing colors ($l): $black of $total code glyphs are black on a dark theme (pages=$pages)"
  else printf '  ok   %-28s pages=%s code_glyphs=%s black=0\n' "listing colors ($l)" "$pages" "$total"; fi
}
for c in "${CATS[@]}"; do for l in en fa; do build_one "$c" "$l"; done; done
echo "== code colors (dark theme, listing across a page break in a breakable box)"
if command -v mutool >/dev/null 2>&1; then for l in en fa; do listing_colors "$l"; done
else echo "  skipped (mutool not found)"; fi
if [[ $PALETTES -eq 1 ]]; then
  echo "== palettes (homework)"
  for p in $(grep -o '^\\atldefpalette{[a-z]*}' "$ROOT/core/atl/atl-palettes.tex" | sed 's/.*{\(.*\)}/\1/'); do build_one homework en "$p" "pal_$p"; done
fi
[[ $FAIL -eq 0 ]] && { echo "PASS"; rm -rf "$TMP"; } || echo "FAILED (outputs kept in $TMP)"
exit $FAIL
