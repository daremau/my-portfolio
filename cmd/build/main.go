// Command build renders the portfolio into a static site.
package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"github.com/a-h/templ"

	"github.com/daremau/my-portfolio/components"
	"github.com/daremau/my-portfolio/internal/content"
)

func main() {
	contentDir := flag.String("content", "content/projects", "directory containing project Markdown files")
	publicDir := flag.String("public", "public", "directory of static assets copied as-is")
	webDir := flag.String("web", "web", "directory containing the CSS and JavaScript sources")
	outDir := flag.String("out", "dist", "directory to write the generated site to")
	flag.Parse()

	if err := run(*contentDir, *publicDir, *webDir, *outDir); err != nil {
		log.Fatal(err)
	}
}

func run(contentDir, publicDir, webDir, outDir string) error {
	projects, err := content.LoadDir(contentDir)
	if err != nil {
		return err
	}

	if err := os.RemoveAll(outDir); err != nil {
		return fmt.Errorf("clean output directory: %w", err)
	}
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return fmt.Errorf("create output directory: %w", err)
	}

	if err := copyDir(publicDir, outDir); err != nil {
		return fmt.Errorf("copy public assets: %w", err)
	}
	for src, dst := range map[string]string{
		filepath.Join(webDir, "global.css"): "styles.css",
		filepath.Join(webDir, "main.js"):    "main.js",
	} {
		if err := copyFile(src, filepath.Join(outDir, dst)); err != nil {
			return fmt.Errorf("copy %s: %w", src, err)
		}
	}

	page := components.Layout(
		"Portfolio | Web Dev, Data Science & AI",
		components.IndexPage(projects),
	)
	if err := writeComponent(filepath.Join(outDir, "index.html"), page); err != nil {
		return fmt.Errorf("render index: %w", err)
	}

	fmt.Printf("built %d projects into %s\n", len(projects), outDir)
	return nil
}

func writeComponent(path string, component templ.Component) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}

	if err := component.Render(context.Background(), file); err != nil {
		_ = file.Close()
		return err
	}
	return file.Close()
}

func copyDir(src, dst string) error {
	info, err := os.Stat(src)
	if errors.Is(err, fs.ErrNotExist) {
		return nil
	}
	if err != nil {
		return err
	}
	if !info.IsDir() {
		return fmt.Errorf("%s is not a directory", src)
	}

	return filepath.WalkDir(src, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		if rel == "." {
			return nil
		}

		target := filepath.Join(dst, rel)
		if entry.IsDir() {
			return os.MkdirAll(target, 0o755)
		}
		return copyFile(path, target)
	})
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return err
	}

	out, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o644)
	if err != nil {
		return err
	}

	if _, err := io.Copy(out, in); err != nil {
		_ = out.Close()
		return err
	}
	return out.Close()
}
