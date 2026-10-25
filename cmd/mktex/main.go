// SPDX-License-Identifier: MPL-2.0

// mktex creates a new LaTeX project from the AdvancedLatexTemplate catalog.
//
//	mktex                                   interactive
//	mktex -l fa -c lab-report -n lab3 -y    non-interactive
//	mktex --list                            show templates, palettes and code themes
//
// A project is self-contained: the engine (core/atl) and, for Persian, the fonts are copied
// into <name>/latex/, so it builds on Windows, macOS, Linux and Overleaf. The templates,
// engine and fonts are embedded in this executable; "--root DIR" reads them from a checkout
// of the repository instead (used by the maintainers' checks).
package main

import (
	"bufio"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"

	assets "github.com/MaximumAsp66915/AdvancedLatexTemplate"
)

var version = "dev" // set at build time: -ldflags "-X main.version=v1.0.0"

const usageText = `Usage: mktex [options]
  -l, --lang LANG        en | fa
  -c, --category NAME    template category (see --list)
  -n, --name FOLDER      name of the folder to create
  -p, --palette NAME     color palette (default: the template's own)
      --listing NAME     code theme: darcula | dark | bright | github
      --logo FILE        logo image copied into the project (cover pages only)
  -s, --set KEY=VALUE    any config key (title, author, course, ...); repeatable
  -o, --out DIR          parent directory for the new folder (default: current folder)
  -y, --yes              do not ask; use defaults for everything not given
      --no-defaults      ignore your personal defaults file (reproducible output)
      --no-open          do not offer to open VS Code
      --list             list templates, palettes and code themes
      --defaults FILE    use this personal defaults file
      --init-defaults    create a personal defaults file and print where it is
      --root DIR         read templates from a repository checkout instead of the built-in copy
      --version
  -h, --help
Environment: MKTEX_LANG sets the default language; MKTEX_DEFAULTS the defaults file.
Personal defaults (name, university, logo...): run "mktex --init-defaults" and edit the file it prints.`

var listings = []string{"darcula", "dark", "bright", "github"}

var promptLabels = map[string]string{
	"title": "Title", "subtitle": "Subtitle", "author": "Author name(s)", "studentid": "Student ID",
	"course": "Course name", "coursecode": "Course code", "instructor": "Instructor",
	"supervisor": "Supervisor", "advisor": "Advisor", "team": "Team members",
	"partners": "Lab partners", "experiment": "Experiment (number and name)",
	"labsection": "Lab section", "performed": "Date performed", "university": "University",
	"department": "Department", "term": "Term (e.g. Spring 2026)", "duedate": "Due date",
	"submissiondate": "Submission date", "affiliation": "Affiliation",
	"keywords": "Keywords (comma separated)", "date": "Date",
}

const vscodeSettings = `{
  "latex-workshop.latex.recipes": [
    { "name": "latexmk (XeLaTeX)", "tools": ["xelatexmk"] }
  ],
  "latex-workshop.latex.tools": [
    {
      "name": "xelatexmk",
      "command": "latexmk",
      "args": ["-xelatex", "-synctex=1", "-interaction=nonstopmode", "-file-line-error", "%DOC%"]
    }
  ],
  "latex-workshop.latex.recipe.default": "first"
}
`

// ---------------------------------------------------------------------------- options ---

type meta struct{ keys, vals []string }

func (m *meta) set(k, v string) {
	for i := range m.keys {
		if m.keys[i] == k {
			m.vals[i] = v
			return
		}
	}
	m.keys = append(m.keys, k)
	m.vals = append(m.vals, v)
}

func (m *meta) has(k string) bool {
	for _, x := range m.keys {
		if x == k {
			return true
		}
	}
	return false
}

type options struct {
	lang, category, name, palette, listing, logo, out string
	root, defaultsFile                                string
	yes, noOpen, noDefaults                           bool
	list, help, showVersion, initDefaults             bool
	meta                                              meta
}

type fatalError struct{ msg string }

func (e fatalError) Error() string { return e.msg }

func die(format string, a ...any) { panic(fatalError{fmt.Sprintf(format, a...)}) }

// parseArgs accepts "-x value", "--long value" and "--long=value".
func parseArgs(args []string) options {
	o := options{out: ".", lang: os.Getenv("MKTEX_LANG"), defaultsFile: os.Getenv("MKTEX_DEFAULTS")}
	valueOf := func(i *int, name string) string {
		if *i+1 >= len(args) {
			die("option %s needs a value", name)
		}
		*i++
		return args[*i]
	}
	for i := 0; i < len(args); i++ {
		a := args[i]
		inline := ""
		hasInline := false
		if strings.HasPrefix(a, "--") {
			if k := strings.Index(a, "="); k > 0 {
				a, inline, hasInline = a[:k], a[k+1:], true
			}
		}
		val := func() string {
			if hasInline {
				return inline
			}
			return valueOf(&i, a)
		}
		switch a {
		case "-l", "--lang":
			o.lang = val()
		case "-c", "--category":
			o.category = val()
		case "-n", "--name":
			o.name = val()
		case "-p", "--palette":
			o.palette = val()
		case "--listing":
			o.listing = val()
		case "--logo":
			o.logo = val()
		case "-s", "--set":
			kv := val()
			k := strings.Index(kv, "=")
			if k < 0 {
				die("--set needs KEY=VALUE")
			}
			o.meta.set(kv[:k], kv[k+1:])
		case "-o", "--out":
			o.out = val()
		case "--root":
			o.root = val()
		case "--defaults":
			o.defaultsFile = val()
		case "-y", "--yes":
			o.yes = true
		case "--no-open":
			o.noOpen = true
		case "--no-defaults":
			o.noDefaults = true
		case "--list":
			o.list = true
		case "--init-defaults":
			o.initDefaults = true
		case "--version":
			o.showVersion = true
		case "-h", "--help":
			o.help = true
		default:
			fmt.Fprintln(os.Stderr, usageText)
			die("unknown option: %s", args[i])
		}
	}
	return o
}

// ---------------------------------------------------------------------------- catalog ---

var assignRe = regexp.MustCompile(`^([A-Za-z_][A-Za-z0-9_]*)=(.*)$`)

// parseAssignments reads KEY="value" / KEY=value lines (template.conf and the defaults file).
func parseAssignments(text string) map[string]string {
	out := map[string]string{}
	for _, line := range strings.Split(text, "\n") {
		line = strings.TrimSpace(strings.TrimSuffix(line, "\r"))
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		m := assignRe.FindStringSubmatch(line)
		if m == nil {
			continue
		}
		v := strings.TrimSpace(m[2])
		if len(v) >= 2 && (v[0] == '"' || v[0] == '\'') && v[len(v)-1] == v[0] {
			v = v[1 : len(v)-1]
		}
		out[m[1]] = v
	}
	return out
}

type catalog struct{ src fs.FS }

func (c catalog) categories() []string {
	entries, err := fs.ReadDir(c.src, "templates")
	if err != nil {
		return nil
	}
	var out []string
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		if _, err := fs.Stat(c.src, path.Join("templates", e.Name(), "template.conf")); err == nil {
			out = append(out, e.Name())
		}
	}
	sort.Strings(out)
	return out
}

func (c catalog) conf(category string) map[string]string {
	b, err := fs.ReadFile(c.src, path.Join("templates", category, "template.conf"))
	if err != nil {
		return map[string]string{}
	}
	return parseAssignments(string(b))
}

var paletteRe = regexp.MustCompile(`(?m)^\\atldefpalette\{([a-z]+)\}`)

func (c catalog) palettes() []string {
	b, err := fs.ReadFile(c.src, "core/atl/atl-palettes.tex")
	if err != nil {
		return nil
	}
	var out []string
	for _, m := range paletteRe.FindAllStringSubmatch(string(b), -1) {
		out = append(out, m[1])
	}
	return out
}

func contains(list []string, s string) bool {
	for _, x := range list {
		if x == s {
			return true
		}
	}
	return false
}

var texEscaper = strings.NewReplacer("&", `\&`, "%", `\%`, "#", `\#`, "_", `\_`)

// ---------------------------------------------------------------------------- prompts ---

type app struct {
	o    options
	cat  catalog
	defs map[string]string
	in   *bufio.Reader
}

func (a *app) readLine(prompt string) string {
	fmt.Print(prompt)
	line, _ := a.in.ReadString('\n')
	return strings.TrimRight(line, "\r\n")
}

func (a *app) ask(prompt, def string) string {
	if a.o.yes {
		return def
	}
	var ans string
	if def != "" {
		ans = a.readLine(fmt.Sprintf("%s [%s]: ", prompt, def))
	} else {
		ans = a.readLine(fmt.Sprintf("%s (optional): ", prompt))
	}
	if ans == "" {
		return def
	}
	return ans
}

func (a *app) choose(prompt, def string, opts []string) string {
	if a.o.yes {
		return def
	}
	fmt.Println(prompt)
	for i, o := range opts {
		mark := ""
		if o == def {
			mark = "  (default)"
		}
		fmt.Printf("  %2d) %s%s\n", i+1, o, mark)
	}
	ans := a.readLine(fmt.Sprintf("Choice [%s]: ", def))
	if ans == "" {
		ans = def
	}
	if n, err := strconv.Atoi(ans); err == nil && n >= 1 && n <= len(opts) {
		ans = opts[n-1]
	}
	return ans
}

// defaultFor returns DEFAULT_<KEY>_<LANG> or DEFAULT_<KEY> from the personal defaults file.
func (a *app) defaultFor(key, lang string) string {
	k := "DEFAULT_" + strings.ToUpper(key)
	if v, ok := a.defs[k+"_"+strings.ToUpper(lang)]; ok && v != "" {
		return v
	}
	return a.defs[k]
}

// ---------------------------------------------------------------------------- files ---

func copyTree(src fs.FS, from, to string, keep func(name string) bool) {
	err := fs.WalkDir(src, from, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel := strings.TrimPrefix(strings.TrimPrefix(p, from), "/")
		target := filepath.Join(to, filepath.FromSlash(rel))
		if d.IsDir() {
			return os.MkdirAll(target, 0o755)
		}
		if keep != nil && !keep(d.Name()) {
			return nil
		}
		b, err := fs.ReadFile(src, p)
		if err != nil {
			return err
		}
		return os.WriteFile(target, b, 0o644)
	})
	if err != nil {
		die("cannot copy %s: %v", from, err)
	}
}

func copyFile(src, dst string) error {
	b, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dst, b, 0o644)
}

func defaultsPath(o options) string {
	if o.defaultsFile != "" {
		return o.defaultsFile
	}
	dir, err := os.UserConfigDir()
	if err != nil {
		return ""
	}
	return filepath.Join(dir, "mktex", "defaults")
}

// ---------------------------------------------------------------------------- main flow ---

func run(args []string, in io.Reader) {
	o := parseArgs(args)
	if o.help {
		fmt.Println(usageText)
		return
	}
	if o.showVersion {
		fmt.Println("mktex", version)
		return
	}

	var src fs.FS = assets.FS
	if o.root != "" {
		src = os.DirFS(o.root)
	}
	a := &app{o: o, cat: catalog{src}, defs: map[string]string{}, in: bufio.NewReader(in)}

	if o.initDefaults {
		p := defaultsPath(o)
		if p == "" {
			die("cannot find a configuration folder; use --defaults FILE")
		}
		if _, err := os.Stat(p); err == nil {
			fmt.Println("Defaults file already exists:", p)
			return
		}
		ex, err := fs.ReadFile(src, "core/project/mktex.defaults.example")
		if err != nil {
			die("built-in example not found: %v", err)
		}
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			die("%v", err)
		}
		if err := os.WriteFile(p, ex, 0o644); err != nil {
			die("%v", err)
		}
		fmt.Println("Created", p, "\nEdit it to set your name, student ID, university, logo ...")
		return
	}

	cats := a.cat.categories()
	if o.list {
		fmt.Println("Templates (each in: en, fa):")
		for _, c := range cats {
			fmt.Printf("  %-16s %s\n", c, a.cat.conf(c)["DESCRIPTION"])
		}
		fmt.Println("Palettes:    ", strings.Join(a.cat.palettes(), " "))
		fmt.Println("Code themes: ", strings.Join(listings, " "))
		return
	}

	// personal defaults
	if !o.noDefaults {
		if p := defaultsPath(o); p != "" {
			if b, err := os.ReadFile(p); err == nil {
				a.defs = parseAssignments(string(b))
			} else if o.defaultsFile != "" {
				die("cannot read defaults file %s", o.defaultsFile)
			}
		}
	}

	// language
	lang := o.lang
	if lang == "" {
		def := a.defs["DEFAULT_LANG"]
		if def == "" {
			def = "en"
		}
		lang = a.choose("Language:", def, []string{"en", "fa"})
	}
	lang = strings.ToLower(lang)
	if lang != "en" && lang != "fa" {
		die("language must be en or fa")
	}

	// template
	if len(cats) == 0 {
		die("no templates found")
	}
	category := o.category
	if category == "" {
		if o.yes {
			category = cats[0]
		} else {
			fmt.Println("Templates:")
			for i, c := range cats {
				fmt.Printf("  %2d) %-16s %s\n", i+1, c, a.cat.conf(c)["DESCRIPTION"])
			}
			ans := a.readLine(fmt.Sprintf("Template [%s]: ", cats[0]))
			if ans == "" {
				ans = cats[0]
			}
			if n, err := strconv.Atoi(ans); err == nil && n >= 1 && n <= len(cats) {
				ans = cats[n-1]
			}
			category = ans
		}
	}
	if !contains(cats, category) {
		die("unknown template '%s' (try: mktex --list)", category)
	}
	conf := a.cat.conf(category)
	srcDir := path.Join("templates", category, lang)
	if _, err := fs.Stat(src, path.Join(srcDir, "latex")); err != nil {
		die("template '%s' has no '%s' version", category, lang)
	}

	// folder
	name := o.name
	if name == "" {
		if o.yes {
			die("--name is required with --yes")
		}
		name = a.readLine("Folder name: ")
	}
	if name == "" {
		die("a folder name is required")
	}
	dest := filepath.Join(o.out, name)
	if _, err := os.Stat(dest); err == nil {
		die("'%s' already exists", dest)
	}

	// palette and code theme
	palettes := a.cat.palettes()
	palette := o.palette
	if palette == "" {
		def := conf["DEFAULT_PALETTE"]
		if def == "" {
			def = "blue"
		}
		palette = a.choose("Color palette:", def, palettes)
	}
	if !contains(palettes, palette) {
		die("unknown palette '%s' (try: mktex --list)", palette)
	}
	listing := o.listing
	if listing == "" {
		def := conf["DEFAULT_LISTING"]
		if def == "" {
			def = "darcula"
		}
		listing = a.choose("Code theme:", def, listings)
	}
	if !contains(listings, listing) {
		die("unknown code theme '%s'", listing)
	}

	// metadata the template asks for
	m := o.meta
	for _, k := range strings.Split(conf["ASK"], ",") {
		k = strings.TrimSpace(k)
		if k == "" || m.has(k) {
			continue
		}
		label := promptLabels[k]
		if label == "" {
			label = k
		}
		if v := a.ask(label, a.defaultFor(k, lang)); v != "" {
			m.set(k, v)
		}
	}
	isCover := conf["TITLEPAGE"] == "cover"
	logo := o.logo
	if logo == "" {
		logo = a.defaultFor("logo", lang)
	}
	if logo == "" && !o.yes && isCover {
		logo = a.ask("Logo image file", "")
	}

	// create
	if err := os.MkdirAll(o.out, 0o755); err != nil {
		die("cannot create %s", o.out)
	}
	fmt.Printf("Creating %s (%s): %s\n", category, lang, dest)
	copyTree(src, srcDir, dest, nil)
	latexDir := filepath.Join(dest, "latex")
	if err := os.MkdirAll(latexDir, 0o755); err != nil {
		die("%v", err)
	}
	copyTree(src, "core/atl", filepath.Join(latexDir, "atl"), nil)
	if _, err := os.Stat(filepath.Join(latexDir, ".latexmkrc")); err != nil {
		if b, err := fs.ReadFile(src, "core/project/latexmkrc"); err == nil {
			_ = os.WriteFile(filepath.Join(latexDir, ".latexmkrc"), b, 0o644)
		}
	}
	if lang == "fa" {
		copyTree(src, "core/fonts", filepath.Join(latexDir, "fonts"), func(n string) bool {
			return strings.HasSuffix(n, ".ttf") || strings.HasPrefix(n, "OFL-")
		})
	}
	if err := os.MkdirAll(filepath.Join(latexDir, "fig"), 0o755); err != nil {
		die("%v", err)
	}
	if logo != "" && isCover {
		ext := filepath.Ext(logo)
		if err := copyFile(logo, filepath.Join(latexDir, "fig", "logo"+ext)); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: logo file '%s' not found; skipped.\n", logo)
		} else {
			m.set("logo", "fig/logo"+ext)
		}
	}

	var cfg strings.Builder
	fmt.Fprintf(&cfg, "%% Written by mktex on %s. Edit freely; anything set here overrides the template defaults.\n", time.Now().Format("2006-01-02"))
	cfg.WriteString("\\atlset{\n")
	fmt.Fprintf(&cfg, "  language = %s,\n  palette  = %s,\n  listing  = %s,\n", lang, palette, listing)
	for i, k := range m.keys {
		if m.vals[i] != "" {
			fmt.Fprintf(&cfg, "  %s = {%s},\n", k, texEscaper.Replace(m.vals[i]))
		}
	}
	cfg.WriteString("}\n")
	if b, err := fs.ReadFile(src, "core/project/config-options."+lang+".tex"); err == nil {
		cfg.Write(b)
	}
	if err := os.WriteFile(filepath.Join(latexDir, "config.tex"), []byte(cfg.String()), 0o644); err != nil {
		die("%v", err)
	}

	vs := filepath.Join(dest, ".vscode")
	if err := os.MkdirAll(vs, 0o755); err != nil {
		die("%v", err)
	}
	if err := os.WriteFile(filepath.Join(vs, "settings.json"), []byte(vscodeSettings), 0o644); err != nil {
		die("%v", err)
	}

	fmt.Printf("Done. Build it with:  cd \"%s\" && latexmk -xelatex main.tex\n", latexDir)

	if !o.yes && !o.noOpen {
		ans := strings.ToLower(a.readLine("Open in VS Code? [Y/n]: "))
		if ans == "" || ans == "y" {
			if p, err := exec.LookPath("code"); err == nil {
				_ = exec.Command(p, "--new-window", dest).Start()
			} else {
				fmt.Println("'code' not found in PATH.")
			}
		}
	}
}

func main() {
	// Double-clicked on Windows (no arguments): keep the window open at the end.
	pause := runtime.GOOS == "windows" && len(os.Args) == 1
	code := 0
	func() {
		defer func() {
			if r := recover(); r != nil {
				if fe, ok := r.(fatalError); ok {
					fmt.Fprintln(os.Stderr, "Error:", fe.msg)
					code = 1
					return
				}
				panic(r)
			}
		}()
		run(os.Args[1:], os.Stdin)
	}()
	if pause {
		fmt.Print("\nPress Enter to close this window...")
		_, _ = bufio.NewReader(os.Stdin).ReadString('\n')
	}
	os.Exit(code)
}
