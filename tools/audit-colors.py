#!/usr/bin/env python3
# SPDX-License-Identifier: MPL-2.0
"""audit-colors.py -- check the color definitions of the engine for readability.

Needs only Python 3 (no TeX). Exit status 1 when something is below target.
  code themes : every token color against its background   (WCAG contrast >= 3, comments >= 2.5)
  palettes    : white text on `primary` (title banner, problem bar, table header)      >= 4.5
                `primary` text on white (headings, links)                              >= 4.5
                `secondary-text` on white (sub-headings; = secondary!78!black)         >= 3
                `accent-text` on white (section numbers; = accent!72!black)            >= 3
                box titles (color!82!black on color!7!white)                         >= 4.5
"""
import re, sys, pathlib

ROOT = pathlib.Path(__file__).resolve().parent.parent / "core" / "atl"

def lum(rgb):
    def f(c):
        c /= 255
        return c / 12.92 if c <= 0.03928 else ((c + 0.055) / 1.055) ** 2.4
    r, g, b = rgb
    return 0.2126 * f(r) + 0.7152 * f(g) + 0.0722 * f(b)

def contrast(a, b):
    la, lb = lum(a), lum(b)
    hi, lo = max(la, lb), min(la, lb)
    return (hi + 0.05) / (lo + 0.05)

def hx(h):
    return tuple(int(h[i:i + 2], 16) for i in (0, 2, 4))

def mix(c, pct, other):          # xcolor  c!pct!other
    return tuple(round(c[i] * pct / 100 + other[i] * (100 - pct) / 100) for i in range(3))

WHITE, BLACK = (255, 255, 255), (0, 0, 0)
fails = []

def check(label, value, target):
    if value < target:
        fails.append(f"{label}: {value:.1f} < {target}")

code = (ROOT / "atl-code.tex").read_text(encoding="utf8")
themes = re.findall(r"^\\atldeflisttheme\{(\w+)\}\{([^}]*)\}", code, re.M)
for name, kv in themes:
    d = dict(x.split("=") for x in kv.split(","))
    bg = hx(d["bg"])
    for key in ("fg", "kw", "str", "cmt", "num", "bi", "deco", "op", "br"):
        check(f"code theme {name}/{key}", contrast(hx(d[key]), bg), 2.5 if key == "cmt" else 3.0)

pal = (ROOT / "atl-palettes.tex").read_text(encoding="utf8")
palettes = re.findall(r"^\\atldefpalette\{(\w+)\}\{([^}]*)\}", pal, re.M)
for name, kv in palettes:
    d = dict(x.strip().split("=") for x in kv.split(","))
    P, S, A = hx(d["primary"]), hx(d["secondary"]), hx(d["accent"])
    st, at = mix(S, 78, BLACK), mix(A, 72, BLACK)
    check(f"palette {name}: white on primary", contrast(WHITE, P), 4.5)
    check(f"palette {name}: primary on white", contrast(P, WHITE), 4.5)
    check(f"palette {name}: secondary-text on white", contrast(st, WHITE), 3.0)
    check(f"palette {name}: accent-text on white", contrast(at, WHITE), 3.0)
    for role, base in (("examplebox", st), ("definitionbox", P)):
        check(f"palette {name}: {role} title", contrast(mix(base, 82, BLACK), mix(base, 7, WHITE)), 4.5)
for role, h in {"info": "1565C0", "warning": "C62828", "success": "2E7D32", "result": "6A1B9A"}.items():
    c = hx(h)
    check(f"semantic {role} box title", contrast(mix(c, 82, BLACK), mix(c, 7, WHITE)), 4.5)

print(f"audit-colors: {len(themes)} code themes, {len(palettes)} palettes checked")
if fails:
    print("FAIL:\n  " + "\n  ".join(fails))
    sys.exit(1)
print("audit-colors: all color pairs meet their targets")
