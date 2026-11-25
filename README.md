# Portfolio

A statically generated portfolio for web development, data science, and AI
engineering projects. The site is written in Go with [templ](https://templ.guide)
components and Markdown content.

## Stack

- **Go 1.26+** — static site generator in `cmd/build`
- **[templ](https://templ.guide)** — type-safe HTML components
- **Markdown + YAML front matter** — project content
- **Plain CSS and JavaScript** — no Node, no bundler

## Structure

```text
.
├── cmd/
│   ├── build/                 # renders the site into dist/
│   └── preview/               # local static file server
├── components/
│   ├── layout.templ           # page shell, nav, footer, theme bootstrap
│   ├── index_page.templ       # hero, about, projects
│   ├── project_card.templ     # project card
│   └── helpers.go             # template helper functions
├── content/projects/          # one Markdown file per project
├── internal/content/          # content loading, validation, sorting + tests
├── public/                    # favicons and other static files
├── web/
│   ├── global.css             # all site styles
│   └── main.js                # theme toggle, filtering, scroll animations
└── dist/                      # generated output (gitignored)
```

## Development

The build is driven by the Makefile:

```sh
make build   # generate templ components, then render dist/
make serve   # serve dist/ at http://localhost:4321
make test    # run the Go tests
make clean   # remove dist/
```

Without Make:

```sh
go tool templ generate
go run ./cmd/build
```

`templ generate` must be re-run after editing any `.templ` file. `go run
./cmd/build` reads `content/projects`, copies `public/` and `web/`, and writes
`dist/index.html`.

## Adding a project

Create a Markdown file in `content/projects/`:

```yaml
---
title: "Neural Network Visualizer"
description: "A real-time tool to visualize weight adjustments in deep learning models."
category: "AI Engineering"
pubDate: 2026-01-26
tags: ["Python", "PyTorch", "React"]
link: "https://github.com"
---
```

`title`, `description`, `category`, `pubDate`, and `link` are required.
`category` must be one of `Web Dev`, `Data Science`, or `AI Engineering`.
`link` must be an `http` or `https` URL. The build fails with the offending file
name if validation does not pass.

## Deployment

`.github/workflows/deploy.yml` builds and deploys `dist/` to GitHub Pages on
every push to `main`. Enable **Settings → Pages → Build and deployment → Source:
GitHub Actions** in the repository.
