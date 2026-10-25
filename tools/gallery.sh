#!/bin/bash
# SPDX-License-Identifier: MPL-2.0
# gallery.sh -- render the first page of every template (both languages) to docs/gallery/*.png
# Needs XeLaTeX, latexmk and pdftoppm (poppler). Uses neutral sample data and the template's own palette.
set -u
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd -P)"
OUT="$ROOT/docs/gallery"; mkdir -p "$OUT"
export PATH="$PATH:/usr/local/go/bin"
MKTEX="${MKTEX:-}"
if [[ -z "$MKTEX" ]]; then MKTEX="$(mktemp -d "${TMPDIR:-/tmp}/atl-mktex.XXXXXX")/mktex"; (cd "$ROOT" && go build -o "$MKTEX" ./cmd/mktex) || exit 1; fi
TMP="$(mktemp -d "${TMPDIR:-/tmp}/atl-gallery.XXXXXX")"
for d in "$ROOT"/templates/*/template.conf; do
  c="$(basename "$(dirname "$d")")"
  for l in en fa; do
    if [[ $l == fa ]]; then A=(-s "title=عنوان نمونه" -s "author=نام نویسنده" -s "course=نام درس" -s "supervisor=استاد راهنما" -s "university=نام دانشگاه" -s "department=نام دانشکده")
    else A=(-s "title=Sample Title" -s "author=Author Name" -s "course=Course Name" -s "supervisor=Supervisor Name" -s "university=University Name" -s "department=Department Name"); fi
    "$MKTEX" --root "$ROOT" --no-defaults -y --no-open -l "$l" -c "$c" -n "${c}_$l" -o "$TMP" "${A[@]}" >/dev/null || continue
    ( cd "$TMP/${c}_$l/latex" && latexmk -xelatex -interaction=nonstopmode main.tex >/dev/null 2>&1 )
    pdftoppm -r 55 -png -f 1 -l 1 -singlefile "$TMP/${c}_$l/latex/main.pdf" "$OUT/${c}-$l" && echo "  $c-$l.png"
  done
done
rm -rf "$TMP"
