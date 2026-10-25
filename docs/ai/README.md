# Agent manual

This manual is for AI agents (and the people who steer them) working on **AdvancedLatexTemplate**. It is
model-neutral: everything is plain Markdown with relative links, nothing depends on one vendor's features.
The short, always-read file is [`AGENTS.md`](../../AGENTS.md) at the repository root; this folder holds the details.
People who want to contribute should also read [`CONTRIBUTING.md`](../../CONTRIBUTING.md).

## How to wire it to your tool

| Tool | What to do |
|---|---|
| **OpenAI Codex / GPT agents** | They read `AGENTS.md` from the repository root automatically. |
| **Claude Code** | It reads `CLAUDE.md`. Create one containing a single line, `@AGENTS.md`, so it loads the same file. This repository git-ignores `CLAUDE.md` (it is treated as a personal, local file), so keep it local or track it in your fork. |
| **Gemini CLI** | It reads `GEMINI.md` by default. Either create a `GEMINI.md` that says "Follow AGENTS.md and docs/ai/README.md", or configure Gemini CLI's context file name to `AGENTS.md` (see its documentation). |
| **Groq** | Groq is an inference provider, not an agent. Use it through an agent tool (for example Aider or OpenCode) and point that tool at `AGENTS.md` (Aider: `--read AGENTS.md`). |
| **Cursor, Copilot, Zed and other editors** | Many support `AGENTS.md`; otherwise add a rules file that points to it. |
| **Plain chat (no file access)** | Paste `AGENTS.md` plus the one topic file for your task (below). |

## Context budget: load as little as you need

Small-context models should read `AGENTS.md` and **one** topic file. Each topic file is self-contained.

| If your task is... | Read |
|---|---|
| Understanding the repo, how a project is generated, why decisions were taken | [architecture.md](architecture.md) |
| Changing the TeX engine (`core/atl/`): a key, palette, code theme, label, box, title page, language layer | [engine.md](engine.md) |
| Adding or editing a template, or keeping English and Persian in step | [templates.md](templates.md) |
| Anything involving Persian text, fonts, right-to-left layout, or `fa/` templates | [persian.md](persian.md) |
| Running checks, verifying output, git and publishing rules, recipes, troubleshooting | [workflow.md](workflow.md) |
| Using the system as an end user (not changing it) | [`docs/USER-MANUAL.md`](../USER-MANUAL.md) |

Editing Persian templates always requires [persian.md](persian.md) as well: most failures come from its rules.

## The five facts everything else depends on

1. **One engine, two language layers.** `core/atl/` is language-neutral; only `atl-lang-en.tex` and
   `atl-lang-fa.tex` know about fonts, direction, numbering and labels.
2. **Self-contained projects.** `mktex` copies the engine (and Persian fonts) into each project, so projects need no
   installation of this repository and build on Windows, macOS, Linux and Overleaf.
3. **Configuration is data.** Every template reads `\atlset{key=value}`; precedence is template defaults, then
   `config.tex`, then `preamble/local.tex`.
4. **Parity is enforced by a script.** `tools/check.sh` fails when `en` and `fa` file lists differ and builds
   all variants. It cannot judge whether the Persian text means the same as the English: that is a manual review.
5. **The repository is public and lives inside a TeX search path.** Never compile in it, and never commit personal
   or third-party material (see [workflow.md](workflow.md#public-repository-hygiene)).

## Keeping this manual correct

The manual is the single maintained description of how the project works. When you change behavior, update the
file that describes it in the same change:

| You changed | Update |
|---|---|
| a configuration key | `core/project/config-options.{en,fa}.tex`, the same block at the end of every `templates/*/*/latex/config.tex`, `docs/USER-MANUAL.md` section 5, and [engine.md](engine.md) |
| a palette, code theme or label | [engine.md](engine.md) "Adding things", `docs/USER-MANUAL.md` section 8 |
| the catalog (a template added or removed) | [templates.md](templates.md), `README.md`, `docs/USER-MANUAL.md` section 6, regenerate `docs/gallery/` |
| a Persian finding | [persian.md](persian.md), with how you verified it |
| a rule learnt the hard way | [engine.md](engine.md) rules list or [persian.md](persian.md) traps list |

Do not record anything you did not run. Mark what is verified and what is only reported by a source.
