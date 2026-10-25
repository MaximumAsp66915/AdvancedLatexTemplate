// SPDX-License-Identifier: MPL-2.0

package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// repoRoot is the checkout the tests read the catalog from (cmd/mktex/../..).
func repoRoot(t *testing.T) string {
	t.Helper()
	root, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	return root
}

func TestTexEscape(t *testing.T) {
	got := texEscaper.Replace(`Q&A_test 50% #1 {x}`)
	want := `Q\&A\_test 50\% \#1 {x}`
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestParseAssignments(t *testing.T) {
	m := parseAssignments("# comment\nA=\"x y\"\nB=plain\n\nC='q'\nDEFAULT_AUTHOR_FA=\"نام\"\r\n")
	for k, want := range map[string]string{"A": "x y", "B": "plain", "C": "q", "DEFAULT_AUTHOR_FA": "نام"} {
		if m[k] != want {
			t.Errorf("%s = %q, want %q", k, m[k], want)
		}
	}
}

func TestParseArgsForms(t *testing.T) {
	o := parseArgs([]string{"-l", "fa", "--category=lab-report", "-s", "title=A=B", "-y"})
	if o.lang != "fa" || o.category != "lab-report" || !o.yes || o.meta.vals[0] != "A=B" {
		t.Errorf("unexpected parse result: %+v", o)
	}
}

func TestCatalogHasBothLanguages(t *testing.T) {
	root := repoRoot(t)
	c := catalog{os.DirFS(root)}
	cats := c.categories()
	if len(cats) < 6 {
		t.Fatalf("expected at least 6 templates, got %v", cats)
	}
	for _, cat := range cats {
		for _, lang := range []string{"en", "fa"} {
			if _, err := os.Stat(filepath.Join(root, "templates", cat, lang, "latex", "main.tex")); err != nil {
				t.Errorf("%s/%s has no latex/main.tex", cat, lang)
			}
		}
	}
	if p := c.palettes(); len(p) != 11 || !contains(p, "sunset") {
		t.Errorf("palettes: %v", p)
	}
}

// Generating every template in both languages must produce a complete, self-contained project.
func TestGenerateEveryTemplate(t *testing.T) {
	root := repoRoot(t)
	out := t.TempDir()
	c := catalog{os.DirFS(root)}
	for _, cat := range c.categories() {
		for _, lang := range []string{"en", "fa"} {
			name := cat + "_" + lang
			run([]string{"--root", root, "--no-defaults", "--no-open", "-y", "-l", lang, "-c", cat,
				"-n", name, "-o", out, "-s", "title=T & U"}, strings.NewReader(""))
			proj := filepath.Join(out, name)
			must := []string{"latex/main.tex", "latex/atl/atl.tex", "latex/config.tex", "latex/.latexmkrc", ".vscode/settings.json", "latex/fig"}
			if lang == "fa" {
				must = append(must, "latex/fonts/Vazirmatn-Regular.ttf", "latex/fonts/OFL-Vazirmatn.txt")
			}
			for _, f := range must {
				if _, err := os.Stat(filepath.Join(proj, filepath.FromSlash(f))); err != nil {
					t.Errorf("%s: missing %s", name, f)
				}
			}
			cfg, _ := os.ReadFile(filepath.Join(proj, "latex", "config.tex"))
			if !strings.Contains(string(cfg), `title = {T \& U}`) || !strings.Contains(string(cfg), "language = "+lang) {
				t.Errorf("%s: config.tex is wrong:\n%s", name, cfg)
			}
		}
	}
}

func TestRefusesExistingFolder(t *testing.T) {
	out := t.TempDir()
	root := repoRoot(t)
	if err := os.Mkdir(filepath.Join(out, "taken"), 0o755); err != nil {
		t.Fatal(err)
	}
	defer func() {
		r := recover()
		if fe, ok := r.(fatalError); !ok || !strings.Contains(fe.msg, "already exists") {
			t.Errorf("expected an 'already exists' error, got %v", r)
		}
	}()
	run([]string{"--root", root, "--no-defaults", "-y", "-l", "en", "-c", "essay", "-n", "taken", "-o", out}, strings.NewReader(""))
}
