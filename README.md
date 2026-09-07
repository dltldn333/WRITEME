# WRITEME

> The blazing fast README assembler for monorepos.

Stop copy-pasting documentation across packages. Write a fragment once, include
it everywhere, and let WRITEME compile plain `README.md` files that GitHub, npm,
and pkg.go.dev render exactly as they do today.

> [!WARNING]
> **Not released yet.** `0.0.x` reserves the name while the compiler is being
> built. Nothing is functional. Follow the repository for progress.

## The idea

Markdown has no component system, so shared documentation gets duplicated into
every package and then drifts. Site frameworks solve this by turning docs into a
website — but a `README.md` has no runtime. It has to stay a static file.

So WRITEME goes the other direction: it compiles *down* to ordinary Markdown.

```
WRITEME.md  ──[ writeme build ]──>  README.md
```

Sources are plain `.md` files, so every editor, linter, and formatter you
already use keeps working.

## Syntax

```markdown
---
props:
  - title
  - version
---

# {{ title }} (v{{ version }})

::badge[Build Status]{color="green"}
::include{src="./install-guide.md"}
```

- **`props`** — a file declares what it accepts, like a function signature
- **`{{ ... }}`** — interpolation
- **`::name[label]{attrs}`** — components, using MDC/remark-directive syntax
- **`::include`** — pull in a fragment; edit it once, every README updates

## Why Go

A monorepo tool should not force `node_modules` onto a Go, Rust, or Python
repository. A single static binary runs in any CI, for any language.

## License

MIT
