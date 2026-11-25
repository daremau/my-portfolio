package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunBuildsSite(t *testing.T) {
	root := t.TempDir()
	contentDir := filepath.Join(root, "content")
	publicDir := filepath.Join(root, "public")
	webDir := filepath.Join(root, "web")
	outDir := filepath.Join(root, "dist")

	mustWrite(t, filepath.Join(contentDir, "example.md"), `---
title: "Test Project"
description: "A project used by the generator test."
stack: ["Go", "HTML"]
cover: "projects/example.png"
repo: "https://example.com"
year: 2026
order: 1
---
`)
	mustWrite(t, filepath.Join(publicDir, "favicon.svg"), "<svg></svg>")
	mustWrite(t, filepath.Join(publicDir, "projects", "example.png"), "not really a png")
	mustWrite(t, filepath.Join(webDir, "global.css"), "body { color: black; }")
	mustWrite(t, filepath.Join(webDir, "main.js"), "console.log('ok');")

	// A stale file must be removed when the output directory is rebuilt.
	mustWrite(t, filepath.Join(outDir, "stale.txt"), "stale")

	if err := run(contentDir, publicDir, webDir, outDir); err != nil {
		t.Fatalf("run: %v", err)
	}

	html := mustRead(t, filepath.Join(outDir, "index.html"))
	for _, want := range []string{
		"<title>Mauricio Da Re — Full-Stack Engineer</title>",
		"Test Project",
		"./styles.css",
		"./main.js",
		"./favicon.svg",
		"./projects/example.png",
		"https://example.com",
	} {
		if !strings.Contains(html, want) {
			t.Errorf("index.html does not contain %q", want)
		}
	}

	for _, name := range []string{
		"styles.css",
		"main.js",
		"favicon.svg",
		"projects/example.png",
	} {
		if _, err := os.Stat(filepath.Join(outDir, name)); err != nil {
			t.Errorf("expected %s in output: %v", name, err)
		}
	}

	if _, err := os.Stat(filepath.Join(outDir, "stale.txt")); !os.IsNotExist(err) {
		t.Errorf("stale.txt should have been removed, stat err = %v", err)
	}
}

func TestRunRejectsMissingCover(t *testing.T) {
	root := t.TempDir()
	contentDir := filepath.Join(root, "content")
	publicDir := filepath.Join(root, "public")
	webDir := filepath.Join(root, "web")
	outDir := filepath.Join(root, "dist")

	mustWrite(t, filepath.Join(contentDir, "example.md"), `---
title: "Missing cover"
description: "A project without its cover image."
stack: ["Go"]
cover: "projects/missing.png"
year: 2026
order: 1
---
`)
	mustWrite(t, filepath.Join(webDir, "global.css"), "body {}")
	mustWrite(t, filepath.Join(webDir, "main.js"), "")

	err := run(contentDir, publicDir, webDir, outDir)
	if err == nil {
		t.Fatal("expected an error")
	}
	if !strings.Contains(err.Error(), "cover") {
		t.Fatalf("error = %v, want cover error", err)
	}
}

func TestRunRendersEmptyState(t *testing.T) {
	root := t.TempDir()
	contentDir := filepath.Join(root, "content")
	publicDir := filepath.Join(root, "public")
	webDir := filepath.Join(root, "web")
	outDir := filepath.Join(root, "dist")

	if err := os.MkdirAll(contentDir, 0o755); err != nil {
		t.Fatal(err)
	}
	mustWrite(t, filepath.Join(webDir, "global.css"), "body {}")
	mustWrite(t, filepath.Join(webDir, "main.js"), "")

	if err := run(contentDir, publicDir, webDir, outDir); err != nil {
		t.Fatalf("run: %v", err)
	}

	html := mustRead(t, filepath.Join(outDir, "index.html"))
	if !strings.Contains(html, "No projects yet.") {
		t.Error("empty state was not rendered")
	}
}

func mustWrite(t *testing.T, path, contents string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(contents), 0o644); err != nil {
		t.Fatal(err)
	}
}

func mustRead(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}
